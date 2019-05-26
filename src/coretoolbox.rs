use directories;
use failure::{bail, Fallible};
use lazy_static::lazy_static;
use serde::{Deserialize, Serialize};
use serde_json;
use signal_hook;
use std::io::prelude::*;
use std::os::unix::process::CommandExt;
use std::path::Path;
use std::process::{Command, Stdio};
use structopt::StructOpt;

lazy_static! {
    static ref APPDIRS: directories::ProjectDirs =
        directories::ProjectDirs::from("com", "coreos", "toolbox").expect("creating appdirs");
}

static CONTAINER_NAME: &str = "coreos-toolbox";
static MAX_UID_COUNT: u32 = 65536;

/// Set of statically known paths to files/directories
/// that we redirect inside the container to /host.
static STATIC_HOST_FORWARDS: &[&str] = &["/run/dbus", "/run/libvirt"];

static PRESERVED_ENV: &[&str] = &[
    "COLORTERM",
    "DBUS_SESSION_BUS_ADDRESS",
    "DESKTOP_SESSION",
    "DISPLAY",
    "USER",
    "LANG",
    "SHELL",
    "SSH_AUTH_SOCK",
    "TERM",
    "VTE_VERSION",
    "XDG_CURRENT_DESKTOP",
    "XDG_DATA_DIRS",
    "XDG_MENU_PREFIX",
    "XDG_RUNTIME_DIR",
    "XDG_SEAT",
    "XDG_SESSION_DESKTOP",
    "XDG_SESSION_ID",
    "XDG_SESSION_TYPE",
    "XDG_VTNR",
    "WAYLAND_DISPLAY",
];

trait CommandRunExt {
    fn run(&mut self) -> Fallible<()>;
}

impl CommandRunExt for Command {
    fn run(&mut self) -> Fallible<()> {
        let r = self.status()?;
        if !r.success() {
            bail!("Child [{:?}] exited: {}", self, r);
        }
        Ok(())
    }
}

#[derive(Debug, StructOpt)]
struct RunOpts {
    #[structopt(
        short = "I",
        long = "image",
        default_value = "registry.fedoraproject.org/f30/fedora-toolbox:30"
    )]
    /// Use a different base image
    image: String,

    #[structopt(short = "N", long = "nested")]
    /// Allow running inside a container
    nested: bool,

    #[structopt(short = "D", long = "destroy")]
    /// Destroy any existing container
    destroy: bool,
}

#[derive(Debug, StructOpt)]
#[structopt(name = "coretoolbox", about = "Toolbox")]
#[structopt(rename_all = "kebab-case")]
enum Opt {
    /// Enter the toolbox
    Run(RunOpts),
    /// Delete the toolbox container
    Rm,
    /// Internal implementation detail; do not use
    RunPid1,
    /// Internal implementation detail; do not use
    Exec,
}

fn cmd_podman() -> Command {
    if let Some(podman) = std::env::var_os("podman") {
        Command::new(podman)
    } else {
        Command::new("podman")
    }
}

/// Returns true if the host is OSTree based
fn is_ostree_based_host() -> bool {
    std::path::Path::new("/run/ostree-booted").exists()
}

#[allow(dead_code)]
enum InspectType {
    Container,
    Image,
}

/// Returns true if an image or container is in the podman
/// storage.
fn podman_has(t: InspectType, name: &str) -> Fallible<bool> {
    let typearg = match t {
        InspectType::Container => "container",
        InspectType::Image => "image",
    };
    Ok(cmd_podman()
        .args(&["inspect", "--type", typearg, name])
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .status()?
        .success())
}

/// Pull a container image if not present
fn ensure_image(name: &str) -> Fallible<()> {
    if !podman_has(InspectType::Image, name)? {
        cmd_podman().args(&["pull", name]).run()?;
    }
    Ok(())
}

/// Parse an extant environment variable as UTF-8
fn getenv_required_utf8(n: &str) -> Fallible<String> {
    if let Some(v) = std::env::var_os(n) {
        Ok(v.to_str()
            .ok_or_else(|| failure::format_err!("{} is invalid UTF-8", n))?
            .to_string())
    } else {
        bail!("{} is unset", n)
    }
}

#[derive(Serialize, Deserialize, Debug)]
struct EntrypointState {
    username: String,
    uid: u32,
    home: String,
    ostree_based_host: bool,
}

fn append_preserved_env(c: &mut Command) -> Fallible<()> {
    for n in PRESERVED_ENV.iter() {
        let v = match std::env::var_os(n) {
            Some(v) => v,
            None => continue,
        };
        let v = v
            .to_str()
            .ok_or_else(|| failure::format_err!("{} contains invalid UTF-8", n))?;
        c.arg(format!("--env={}={}", n, v));
    }
    Ok(())
}

fn create(opts: &RunOpts) -> Fallible<()> {
    ensure_image(&opts.image)?;

    if podman_has(InspectType::Container, CONTAINER_NAME)? {
        return Ok(());
    }

    // exec ourself as the entrypoint.  In the future this
    // would be better with podman fd passing.
    let self_bin = std::fs::read_link("/proc/self/exe")?;
    let self_bin = self_bin
        .as_path()
        .to_str()
        .ok_or_else(|| failure::err_msg("non-UTF8 self"))?;

    let runtime_dir = getenv_required_utf8("XDG_RUNTIME_DIR")?;
    let statefile = "coreos-toolbox.initdata";

    let mut podman = cmd_podman();
    podman.args(&[
        "create",
        "--interactive",
        "--tty",
        "--hostname=toolbox",
        "--network=host",
        "--privileged",
        "--security-opt=label=disable",
    ]);
    podman.arg(format!("--name={}", CONTAINER_NAME));
    podman.arg(format!("--volume={}:/usr/bin/toolbox:ro", self_bin));
    let real_uid: u32 = nix::unistd::getuid().into();
    // In true privileged mode we don't use userns
    if real_uid != 0 {
        let uid_plus_one = real_uid + 1;
        let max_minus_uid = MAX_UID_COUNT - real_uid;
        podman.args(&[
            format!("--uidmap={}:0:1", real_uid),
            format!("--uidmap=0:1:{}", real_uid),
            format!(
                "--uidmap={}:{}:{}",
                uid_plus_one, uid_plus_one, max_minus_uid
            ),
        ]);
    }
    // TODO: Detect what devices are accessible
    for p in &["/dev/bus", "/dev/dri", "/dev/fuse"] {
        if Path::new(p).exists() {
            podman.arg(format!("--volume={}:{}:rslave", p, p));
        }
    }
    for p in &["/usr", "/var", "/etc", "/run", "/tmp"] {
        podman.arg(format!("--volume={}:/host{}:rslave", p, p));
    }
    if is_ostree_based_host() {
        podman.arg(format!("--volume=/sysroot:/host/sysroot:rslave"));
    } else {
        for p in &["/media", "/mnt", "/home", "/srv"] {
            podman.arg(format!("--volume={}:/host{}:rslave", p, p));
        }
    }
    append_preserved_env(&mut podman)?;
    podman.arg(format!("--env=TOOLBOX_STATEFILE={}", statefile));

    {
        let state = EntrypointState {
            username: getenv_required_utf8("USER")?,
            uid: real_uid,
            home: getenv_required_utf8("HOME")?,
            ostree_based_host: is_ostree_based_host(),
        };
        let w = std::fs::File::create(format!("{}/{}", runtime_dir, statefile))?;
        let mut w = std::io::BufWriter::new(w);
        serde_json::to_writer(&mut w, &state)?;
        w.flush()?;
    }

    podman.arg(&opts.image);
    podman.args(&["/usr/bin/toolbox", "run-pid1"]);
    podman.stdout(Stdio::null());
    podman.run()?;
    Ok(())
}

fn in_container() -> bool {
    Path::new("/run/.containerenv").exists()
}

fn run(opts: &RunOpts) -> Fallible<()> {
    if in_container() && !opts.nested {
        bail!("Already inside a container");
    }

    if opts.destroy {
        rm()?;
    }

    create(&opts)?;

    cmd_podman()
        .args(&["start", CONTAINER_NAME])
        .stdout(Stdio::null())
        .run()?;

    let mut podman = cmd_podman();
    podman.args(&["exec", "--interactive", "--tty"]);
    append_preserved_env(&mut podman)?;
    podman.args(&[CONTAINER_NAME, "/usr/bin/toolbox", "exec"]);
    return Err(podman.exec().into());
}

fn rm() -> Fallible<()> {
    if !podman_has(InspectType::Container, CONTAINER_NAME)? {
        return Ok(());
    }
    let mut podman = cmd_podman();
    podman
        .args(&["rm", "-f", CONTAINER_NAME])
        .stdout(Stdio::null());
    Err(podman.exec().into())
}

fn run_pid1(_opts: Opt) -> Fallible<()> {
    unsafe {
        signal_hook::register(signal_hook::SIGCHLD, waitpid_all)?;
        signal_hook::register(signal_hook::SIGTERM, || std::process::exit(0))?;
    };
    loop {
        std::thread::sleep(std::time::Duration::from_secs(1_000_000));
    }
}

fn waitpid_all() {
    use nix::sys::wait::WaitStatus;
    loop {
        match nix::sys::wait::waitpid(None, Some(nix::sys::wait::WaitPidFlag::WNOHANG)) {
            Ok(status) => {
                match status {
                    WaitStatus::StillAlive => break,
                    _ => {},
                }
            }
            Err(_) => break,
        }
    }
}

mod entrypoint {
    use super::CommandRunExt;
    use super::EntrypointState;
    use failure::{bail, Fallible, ResultExt};
    use fs2::FileExt;
    use rayon::prelude::*;
    use std::io::prelude::*;
    use std::os::unix;
    use std::os::unix::process::CommandExt;
    use std::path::Path;
    use std::process::Command;

    static CONTAINER_INITIALIZED_LOCK: &str = "/run/coreos-toolbox.lock";
    static CONTAINER_INITIALIZED_STAMP: &str = "/run/coreos-toolbox.initialized";

    /// Update /etc/passwd with the same user from the host,
    /// and bind mount the homedir.
    fn adduser(state: &EntrypointState) -> Fallible<()> {
        if state.uid == 0 {
            return Ok(());
        }
        let uidstr = format!("{}", state.uid);
        Command::new("useradd")
            .args(&[
                "--no-create-home",
                "--home-dir",
                &state.home,
                "--uid",
                &uidstr,
                "--groups",
                "wheel",
                state.username.as_str(),
            ])
            .run()?;

        // Bind mount the homedir rather than use symlinks
        // as various software is unhappy if the path isn't canonical.
        std::fs::create_dir_all(&state.home)?;
        let uid = nix::unistd::Uid::from_raw(state.uid);
        let gid = nix::unistd::Gid::from_raw(state.uid);
        nix::unistd::chown(state.home.as_str(), Some(uid), Some(gid))?;
        let host_home = format!("/host{}", state.home);
        Command::new("mount")
            .args(&["--bind", host_home.as_str(), state.home.as_str()])
            .run()?;
        Ok(())
    }

    /// Symlink a path e.g. /run/dbus/system_bus_socket to the
    /// /host equivalent, creating any necessary parent directories.
    fn host_symlink<P: AsRef<Path> + std::fmt::Display>(p: P) -> Fallible<()> {
        let path = p.as_ref();
        std::fs::create_dir_all(path.parent().unwrap())?;
        match std::fs::remove_dir_all(path) {
            Ok(_) => Ok(()),
            Err(ref e) if e.kind() == std::io::ErrorKind::NotFound => Ok(()),
            Err(e) => Err(e),
        }?;
        unix::fs::symlink(format!("/host{}", p), path)?;
        Ok(())
    }

    /// Podman unprivileged mode has a bug where it exposes the host
    /// selinuxfs which is bad because it can make e.g. librpm
    /// think it can do domain transitions to rpm_exec_t, which
    /// isn't actually permitted.
    fn workaround_podman_selinux() -> Fallible<()> {
        let sysfs_selinux = "/sys/fs/selinux";
        if Path::new(sysfs_selinux).join("status").exists() {
            Command::new("mount")
                .args(&["--bind", "/usr/share/empty", sysfs_selinux])
                .run()?;
        }
        Ok(())
    }

    fn init_container() -> Fallible<()> {
        let initstamp = Path::new(CONTAINER_INITIALIZED_STAMP);
        if initstamp.exists() {
            return Ok(());
        }

        let lockf = std::fs::OpenOptions::new()
            .read(true)
            .write(true)
            .create(true)
            .open(CONTAINER_INITIALIZED_LOCK)?;
        lockf.lock_exclusive()?;

        workaround_podman_selinux()?;

        let runtime_dir = super::getenv_required_utf8("XDG_RUNTIME_DIR")?;
        let state: EntrypointState = {
            let p = format!("/host/{}/{}", runtime_dir, "coreos-toolbox.initdata");
            let f =
                std::fs::File::open(&p).with_context(|e| format!("Opening statefile: {}", e))?;
            std::fs::remove_file(p)?;
            serde_json::from_reader(std::io::BufReader::new(f))?
        };

        let var_mnt_dirs = ["/srv", "/mnt"];
        if state.ostree_based_host {
            var_mnt_dirs.par_iter().chain(["/home"].par_iter())
            .try_for_each(|d| -> Fallible<()> {
                let hostd = format!("/host{}", d);
                let vard = format!("var{}", d);
                unix::fs::symlink(vard, hostd)?;
                Ok(())
            })?;
        }

        // Remove anaconda cruft
        std::fs::read_dir("/tmp")?.try_for_each(|e| -> Fallible<()> {
            let e = e?;
            if let Some(name) = e.file_name().to_str() {
                if name.starts_with("ks-script-") {
                    std::fs::remove_file(e.path())?;
                }
            }
            Ok(())
        })?;

        // Propagate data and temporary directories to the host
        var_mnt_dirs.par_iter().chain(["/tmp", "/var/tmp"].par_iter())
            .try_for_each(|d| -> Fallible<()> {
                std::fs::remove_dir(d)?;
                let hostd = format!("/host{}", d);
                unix::fs::symlink(hostd, d)?;
                Ok(())
            })
            .with_context(|e| format!("Symlinking host dir: {}", e))?;

        // And forward the runtime dir
        host_symlink(runtime_dir).with_context(|e| format!("Forwarding runtime dir: {}", e))?;

        // These symlinks into /host are our set of default forwarded APIs/state
        // directories.
        super::STATIC_HOST_FORWARDS
            .par_iter()
            .try_for_each(host_symlink)
            .with_context(|e| format!("Enabling static host forwards: {}", e))?;

        // Allow sudo
        || -> Fallible<()> {
            let f = std::fs::File::create(format!("/etc/sudoers.d/toolbox-{}", state.username))?;
            let mut f = std::io::BufWriter::new(f);
            writeln!(&mut f, "{} ALL=(ALL) NOPASSWD: ALL", state.username)?;
            f.flush()?;
            Ok(())
        }()
        .with_context(|e| format!("Enabling sudo: {}", e))?;

        adduser(&state)?;
        let _ = std::fs::File::create(&initstamp)?;

        Ok(())
    }

    pub(crate) fn exec() -> Fallible<()> {
        if !super::in_container() {
            bail!("Not inside a container");
        }
        init_container().with_context(|e| format!("Initializing container: {}", e))?;
        let initstamp = Path::new(CONTAINER_INITIALIZED_STAMP);
        if !initstamp.exists() {
            bail!("toolbox not initialized");
        }
        let username = super::getenv_required_utf8("USER")?;
        let su_preserved_env_arg =
            format!("--whitelist-environment={}", super::PRESERVED_ENV.join(","));
        Err(Command::new("setpriv")
            .args(&[
                "--inh-caps=-all",
                "su",
                su_preserved_env_arg.as_str(),
                "-",
                &username,
            ])
            .env_remove("TOOLBOX_STATEFILE")
            .exec()
            .into())
    }
}

/// Primary entrypoint
fn main() {
    || -> Fallible<()> {
        let opts = Opt::from_args();
         match opts {
            Opt::Run(ref runopts) => run(runopts),
            Opt::Exec => entrypoint::exec(),
            Opt::Rm => rm(),
            Opt::RunPid1 => run_pid1(opts),
        }
    }()
    .unwrap_or_else(|e| {
        eprintln!("error: {}", e);
        std::process::exit(1)
    })
}
