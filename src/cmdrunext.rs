use failure::{Fallible, bail};
use std::process::Command;

pub(crate) trait CommandRunExt {
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
