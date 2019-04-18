use structopt::StructOpt;
use std::path::Path;
use std::process::{Command, Stdio};
use std::os::unix::process::CommandExt;
use std::io::prelude::*;
use directories;
use failure::{Fallible, bail};
use lazy_static::lazy_static;
use serde_json;
use serde::{Serialize, Deserialize};

lazy_static! {
    static ref APPDIRS : directories::ProjectDirs = directories::ProjectDirs::from("com", "coreos", "toolbox").expect("creating appdirs");
}

static CONTAINER_NAME : &str = "coreos-toolbox";
static MAX_UID_COUNT : u32 = 65536;

static PRESERVED_ENV : &[&str] = &["COLORTERM", 
        "DBUS_SESSION_BUS_ADDRESS",
        "DESKTOP_SESSION",
        "DISPLAY",
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
];

trait CommandRunExt {
    fn run(&mut self) -> Fallible<()>;
}

impl CommandRunExt for Command {
    fn run(&mut self) -> Fallible<()> {
        let r = self.status()?;
        if !r.success() {
            bail!("Child exited: {}", r);
        }
        Ok(())
    }
}

#[derive(Debug, StructOpt)]
#[structopt(name = "coretoolbox", about = "Toolbox")]
#[structopt(rename_all = "kebab-case")]
/// Main options struct
struct Opt {
    #[structopt(short = "I", long = "image", default_value = "registry.fedoraproject.org/f30/fedora-toolbox:30")]
    /// Use a versioned installer binary
    image: String,

    #[structopt(subcommand)]
    cmd: Option<Cmd>,
}

#[derive(Debug, StructOpt)]
#[structopt(rename_all = "kebab-case")]
enum Cmd {
    Entrypoint,
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
    Ok(cmd_podman().args(&["inspect", "--type", typearg, name])
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .status()?.success())
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
        Ok(v.to_str().ok_or_else(|| failure::format_err!("{} is invalid UTF-8", n))?.to_string())
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

fn create(opts: &Opt) -> Fallible<()> {
    ensure_image(&opts.image)?;

    if podman_has(InspectType::Container, CONTAINER_NAME)? {
        return Ok(())
    }

    // exec ourself as the entrypoint.  In the future this
    // would be better with podman fd passing.
    let self_bin = std::fs::read_link("/proc/self/exe")?;
    let self_bin = self_bin.as_path().to_str().ok_or_else(|| failure::err_msg("non-UTF8 self"))?;

    let runtime_dir = getenv_required_utf8("XDG_RUNTIME_DIR")?;
    let statefile = "coreos-toolbox.initdata";

    let mut podman = cmd_podman();
    podman.args(&["create", "--interactive", "--tty", "--hostname=toolbox",
                  "--name=coreos-toolbox", "--network=host",
                  "--privileged", "--security-opt=label=disable"]);
    podman.arg(format!("--volume={}:/usr/libexec/toolbox.entrypoint:rslave", self_bin));
    let real_uid : u32 = nix::unistd::getuid().into();
    let uid_plus_one = real_uid + 1;             
    let max_minus_uid = MAX_UID_COUNT - real_uid;     
    podman.args(&[format!("--uidmap={}:0:1", real_uid),
                  format!("--uidmap=0:1:{}", real_uid),
                  format!("--uidmap={}:{}:{}", uid_plus_one, uid_plus_one, max_minus_uid)]);
    // TODO: Detect what devices are accessible
    for p in &["/dev/bus", "/dev/dri", "/dev/fuse"] {
        if Path::new(p).exists() {
            podman.arg(format!("--volume={}:{}:rslave", p, p));
        }
    }
    for p in &["/usr", "/var", "/etc", "/run"] {
        podman.arg(format!("--volume={}:/host{}:rslave", p, p));
    }    
    if is_ostree_based_host() {
        podman.arg(format!("--volume=/sysroot:/host/sysroot:rslave"));
    } else {
        for p in &["/media", "/mnt", "/home", "/srv"] {
            podman.arg(format!("--volume={}:/host{}:rslave", p, p));
        }           
    }
    for n in PRESERVED_ENV.iter() {
        let v = match std::env::var_os(n) {
            Some(v) => v,
            None => continue, 
        };
        let v = v.to_str().ok_or_else(|| failure::format_err!("{} contains invalid UTF-8", n))?;
        podman.arg(format!("--env={}={}", n, v));
    }
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

    podman.arg("--entrypoint=/usr/libexec/toolbox.entrypoint");
    podman.arg(&opts.image);
    podman.stdout(Stdio::null());
    podman.run()?;
    Ok(())
}

fn run(opts: Opt) -> Fallible<()> {
    create(&opts)?;

    cmd_podman().args(&["start", CONTAINER_NAME]).run()?;

    let mut podman = cmd_podman();
    podman.args(&["exec", "--interactive", "--tty", CONTAINER_NAME, "/usr/bin/toolbox", "exec"]);
    return Err(podman.exec().into());
}

mod entrypoint {
    use failure::{Fallible, bail, ResultExt};
    use std::process::Command;
    use std::io::prelude::*;
    use std::path::Path;
    use std::os::unix::process::CommandExt;
    use std::os::unix;
    use super::CommandRunExt;
    use super::EntrypointState;

    static CONTAINER_INITIALIZED_STAMP : &str = "/run/coreos-toolbox.initialized";

    fn remove_file_if_exists<P: AsRef<std::path::Path>>(p: P) -> Fallible<()> {
        let p = p.as_ref();
        match std::fs::remove_file(p) {
            Ok(_) => Ok(()),
            Err(ref e) if e.kind() == std::io::ErrorKind::NotFound => Ok(()),
            Err(e) => Err(e.into()),
        }
    }

    /// Update /etc/passwd with the same user from the host,
    /// and bind mount the homedir.
    fn adduser(state: &EntrypointState) -> Fallible<()> {
        let uidstr = format!("{}", state.uid);
        Command::new("useradd")
            .args(&["--no-create-home", "--home-dir", &state.home,
                    "--uid", &uidstr,
                    "--groups", "wheel",
                    state.username.as_str()])
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

    fn init_container() -> Fallible<()> {
        let initstamp = Path::new(CONTAINER_INITIALIZED_STAMP);
        if initstamp.exists() {
            return Ok(())
        }

        let runtime_dir = super::getenv_required_utf8("XDG_RUNTIME_DIR")?;
        let state : EntrypointState = {
            let p = format!("/host/{}/{}", runtime_dir, "coreos-toolbox.initdata");
            let f = std::fs::File::open(&p).with_context(|e| format!("Opening statefile: {}", e))?;
            std::fs::remove_file(p)?;
            serde_json::from_reader(std::io::BufReader::new(f))?
        };

        if state.ostree_based_host {
            std::fs::remove_dir("/home")?;
            unix::fs::symlink("../var/home", "/home")?;
            std::fs::create_dir("/var/home")?;
        }

        // Propagate "data" directories to the host
        for d in &["/srv", "/media", "/mnt"] {
            std::fs::remove_dir(d)?;
            let hostd = format!("host{}", d);
            unix::fs::symlink(hostd, d)?;
        }

        // Symlink ourself
        remove_file_if_exists("/usr/bin/toolbox")?;
        unix::fs::symlink("../libexec/toolbox.entrypoint", "/usr/bin/toolbox")?;

        // Set up runtime dir for this user
        let runtime_dir_path = std::path::Path::new(&runtime_dir);
        std::fs::create_dir_all(runtime_dir_path.parent().unwrap())?;
        unix::fs::symlink(format!("/host{}", &runtime_dir), runtime_dir_path)?;

        // Allow sudo
        || -> Fallible<()> {
            let f = std::fs::File::create(format!("/etc/sudoers.d/toolbox-{}", state.username))?;
            let mut f = std::io::BufWriter::new(f);
            writeln!(&mut f, "{} ALL=(ALL) NOPASSWD: ALL", state.username)?;
            f.flush()?;
            Ok(())
        }().with_context(|e| format!("Enabling sudo: {}", e))?;

        adduser(&state)?;
        let _ = std::fs::File::create(&initstamp)?;

        Ok(())
    }

    /// Called to initialize the container
    pub(crate) fn entrypoint() -> Fallible<()> {
        if !Path::new("/run/.containerenv").exists() {
            bail!("Not running in a container");
        }

        init_container().with_context(|e| format!("Initializing container: {}", e))?;
        // And now we just wait for other processes to exec
        Err(Command::new("sleep").arg("infinity").exec().into())
    }


    pub(crate) fn exec() -> Fallible<()> {
        let initstamp = Path::new(CONTAINER_INITIALIZED_STAMP);        
        if !initstamp.exists() {
            bail!("toolbox not initialized");
        }
        let username = super::getenv_required_utf8("USER")?;        
        Err(Command::new("setpriv")
            .args(&["--ruid", &username, "--init-groups", "--inh-caps=-all", "/bin/bash"])
            .env_remove("TOOLBOX_STATEFILE")
            .exec().into())
    }
}

/// Primary entrypoint
fn main() {
    || -> Fallible<()> {
        let argv0 = std::env::args().next().expect("argv0");
        if argv0.ends_with(".entrypoint") {
            return entrypoint::entrypoint();
        }
        let opts = Opt::from_args();
        if let Some(cmd) = opts.cmd.as_ref() {
            match cmd {
                Cmd::Entrypoint => entrypoint::entrypoint(),
                Cmd::Exec => entrypoint::exec(),
            }
        } else {
            run(opts)
        }
    }().unwrap_or_else(|e| {
        eprintln!("{}", e);
        std::process::exit(1)
    })
}
