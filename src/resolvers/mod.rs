pub(crate) mod install;

pub trait Resolver {
    fn cmd(&self) -> eyre::Result<clap::Command>;
    fn matches(&self, args: &clap::ArgMatches) -> eyre::Result<()>;
}

pub type DynResolver = Box<dyn Resolver + Send + Sync>;
