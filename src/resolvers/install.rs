use crate::deps;

use super::{DynResolver, Resolver};

#[derive(Debug, Clone)]
pub struct Install {
    deps: deps::Deps,
}

impl Install {
    pub fn new(deps: deps::Deps) -> DynResolver {
        Box::new(Self { deps })
    }
}

impl Resolver for Install {
    fn cmd(&self) -> eyre::Result<clap::Command> {
        let install = clap::Command::new("install");

        Ok(install)
    }

    fn matches(&self, _args: &clap::ArgMatches) -> eyre::Result<()> {
        let char = self.deps.parser.parse()?;

        dbg!(char);

        Ok(())
    }
}
