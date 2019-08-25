use std::process::{Command, Stdio};
use std::io::prelude::*;
use serde::{Deserialize};
use serde_json;
use failure::{Fallible, bail};

#[allow(dead_code)]
pub(crate) enum InspectType {
    Container,
    Image,
}

#[derive(Deserialize, Clone, Debug)]
pub(crate) struct ImageInspect {
    pub id: String,
    pub names: Vec<String>,
}

pub(crate) fn cmd() -> Command {
    if let Some(podman) = std::env::var_os("podman") {
        Command::new(podman)
    } else {
        Command::new("podman")
    }
}

/// Returns true if an image or container is in the podman
/// storage.
pub(crate) fn has_object(t: InspectType, name: &str) -> Fallible<bool> {
    let typearg = match t {
        InspectType::Container => "container",
        InspectType::Image => "image",
    };
    Ok(cmd()
        .args(&["inspect", "--type", typearg, name])
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .status()?
        .success())
}

pub(crate) fn image_inspect<I, S>(args: I) -> Fallible<Vec<ImageInspect>>
    where I: IntoIterator<Item=S>, S: AsRef<std::ffi::OsStr>
 {
    let mut proc = cmd()
        .stdout(Stdio::piped())
        .args(&["images", "--format", "json"])
        .args(args)
        .spawn()?;
    let sout = proc.stdout.take().expect("stdout piped");
    let mut sout = std::io::BufReader::new(sout);
    let res = if sout.fill_buf()?.len() > 0 {
        serde_json::from_reader(sout)?
    } else {
        Vec::new()
    };
    if !proc.wait()?.success() {
        bail!("podman images failed")
    }
    Ok(res)
}
