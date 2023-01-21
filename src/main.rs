fn main() -> eyre::Result<()> {
    color_eyre::install()?;

    let args = std::env::args();

    Char::new().execute_from(
        args.collect::<Vec<String>>()
            .iter()
            .map(|s| s.as_str())
            .collect::<Vec<&str>>()
            .as_slice(),
    )?;

    Ok(())
}

mod library {
    use std::path::PathBuf;

    pub struct Shell {
        path: PathBuf,
    }

    impl Shell {
        pub fn new(path: Option<PathBuf>) -> eyre::Result<Self> {
            Ok(Self {
                path: path.unwrap_or(std::env::current_dir()?),
            })
        }
        pub fn execute_shell(&self, args: &[&str]) -> eyre::Result<String> {
            let output = std::process::Command::new(
                args.get(0)
                    .ok_or(eyre::anyhow!("no first arg in shell command"))?,
            )
            .args(args.get(1..).unwrap_or(&[]))
            .current_dir(&self.path)
            .output()?;

            Ok(String::from_utf8(output.stdout)?)
        }
    }
}

fn into_static_str(s: String) -> &'static str {
    Box::leak(s.into_boxed_str())
}

trait Pipeline {
    fn name(&self) -> String;
    fn cmd(&self) -> clap::Command {
        let name = self.name();

        clap::Command::new(into_static_str(name))
    }
    fn execute(&self, args: &clap::ArgMatches) -> eyre::Result<()>;
}

struct BuildPipeline {}

impl Pipeline for BuildPipeline {
    fn name(&self) -> String {
        "build_pipeline".into()
    }

    fn cmd(&self) -> clap::Command {
        let name = self.name();

        clap::Command::new(into_static_str(name)).arg(clap::Arg::new("some").long("some"))
    }

    fn execute(&self, args: &clap::ArgMatches) -> eyre::Result<()> {
        let some = args
            .get_one::<String>("some")
            .ok_or(eyre::anyhow!("some is missing"))?;

        let output = library::Shell::new(None)?.execute_shell(&["echo", some])?;

        println!("{output}");

        Ok(())
    }
}

type DynPipeline = Box<dyn Pipeline + Send + Sync>;
struct Pipelines {
    pipelines: Vec<DynPipeline>,
}

impl Pipelines {
    fn list(&self) -> eyre::Result<Vec<String>> {
        let cmds = self.pipelines.iter().map(|p| p.name()).collect();

        Ok(cmds)
    }
}

#[allow(dead_code)]
pub struct Char {
    pipelines: Pipelines,
}

impl Char {
    pub fn new() -> Self {
        let build_pipeline = BuildPipeline {};
        Self {
            pipelines: Pipelines {
                pipelines: vec![Box::new(build_pipeline)],
            },
        }
    }

    fn execute_from(self, args: &[&str]) -> eyre::Result<()> {
        let matches = self.main_cmd().get_matches_from(args);

        match matches.subcommand() {
            Some(("list", _)) => self.execute_list()?,
            Some(("run", sub)) => self.execute_run(sub)?,
            _ => eyre::bail!("no command matches, please [char --help] to see available commands"),
        }

        Ok(())
    }

    fn main_cmd(&self) -> clap::Command {
        clap::Command::new("char")
            .subcommand(self.list_cmd())
            .subcommand(self.run_cmd())
    }

    fn list_cmd(&self) -> clap::Command {
        clap::Command::new("list")
    }

    fn run_cmd(&self) -> clap::Command {
        clap::Command::new("run").subcommands(self.pipelines.pipelines.iter().map(|p| p.cmd()))
    }

    fn execute_list(&self) -> eyre::Result<()> {
        println!("list");
        self.pipelines
            .pipelines
            .iter()
            .for_each(|p| println!("  - {}", p.name()));

        Ok(())
    }

    fn execute_run(&self, sub: &clap::ArgMatches) -> eyre::Result<()> {
        if let Some((name, subm)) = sub.subcommand() {
            let pipelines: Vec<&DynPipeline> = self
                .pipelines
                .pipelines
                .iter()
                .filter(|p| p.name() == name)
                .collect();

            if let Some(pipeline) = pipelines.get(0) {
                println!("{}", pipeline.name());

                pipeline.execute(subm)?;

                return Ok(());
            }
        }
        eyre::bail!("no command matches, please [char list] to see available commands");
    }
}

#[cfg(test)]
mod tests {
    use crate::Char;

    #[test]
    fn execute_list() {
        Char::new().execute_from(&["char", "list"]).unwrap();
    }

    #[test]
    fn execute_run() {
        Char::new()
            .execute_from(&["char", "run", "build_pipeline", "--some", "arg"])
            .unwrap();
    }
}
