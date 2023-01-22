pub mod cli;
mod deps;
mod models;
mod parser;
mod resolvers;

fn main() -> eyre::Result<()> {
    color_eyre::install()?;

    let args = std::env::args();

    let deps = deps::Deps::default();

    let c = cli::Cli::new(deps)?;
    c.matches(
        args.collect::<Vec<String>>()
            .iter()
            .map(|s| s.as_str())
            .collect::<Vec<&str>>()
            .as_slice(),
    )?;

    let p = std::path::PathBuf::from("examples/service/char.toml");

    let char = std::fs::read_to_string(p)?.parse::<models::Char>()?;

    dbg!(char);

    Ok(())
}
