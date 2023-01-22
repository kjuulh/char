use std::{path::PathBuf};

use crate::{
    deps,
    resolvers::{install, Resolver},
};

pub struct Cli {
    deps: deps::Deps,
    install: Box<dyn Resolver + Send + Sync>,
}

impl Cli {
    pub fn new(deps: deps::Deps) -> eyre::Result<Self> {
        Ok(Self {
            deps: deps.clone(),
            install: install::Install::new(deps),
        })
    }

    pub fn matches(self, args: &[&str]) -> eyre::Result<()> {
        let mut cli = clap::Command::new("char")
            .arg(clap::Arg::new("path").long("path").short('p'))
            .subcommand(self.install.cmd()?);

        let matches = cli.clone().get_matches_from(args);

        let path = matches.get_one::<String>("path");
        if let Some(p) = path {
            let validated_path = PathBuf::from(p);
            if !validated_path.exists() {
                eyre::bail!("no char.toml exists at --path")
            }
            self.deps.parser.set_path(validated_path);
        }

        match matches.subcommand() {
            Some(("install", args)) => self.install.matches(args)?,
            _ => cli.print_help()?,
        }

        Ok(())
    }
}
