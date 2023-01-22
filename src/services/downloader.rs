use std::{fs::canonicalize, path::PathBuf};

use crate::{
    context::{Context, ContextBuilder},
    models::{Char, Conf},
    parser::Parser,
};

#[derive(Debug, Clone)]
pub struct Downloader {
    parser: Parser,
}
#[allow(dead_code)]
impl Downloader {
    pub fn new(parser: Parser) -> Self {
        Self { parser }
    }

    /// Unfolds char
    /// 1. Download path
    /// 2. Parse char in downloaded path
    /// 3. Repeat from 1. until there are no more parents
    pub fn download(&self) -> eyre::Result<Context> {
        let mut context_builder = ContextBuilder::new();

        let char = self.parser.parse()?;
        context_builder = context_builder.add_char(&char);
        let first_char_path = self.parser.get_path()?;

        let mut root = std::env::current_dir()?;
        root = root.join(&first_char_path);
        root.push(".char");
        let output = self.create_output_dir(&root)?;

        let mut parent_char = char;
        let path = first_char_path;
        loop {
            parent_char = match &parent_char {
                Char::Application {
                    char,
                    application: _,
                    config: _,
                } => match char {
                    Some(c) => self.download_plan(c, &path, &output)?,
                    None => {
                        break;
                    }
                },
                Char::Plan {
                    char,
                    plan: _,
                    config: _,
                } => match char {
                    Some(_c) => todo!(),
                    None => {
                        break;
                    }
                },
            }
        }

        Ok(context_builder.build())
    }

    fn download_plan(
        &self,
        conf: &Conf,
        path: &PathBuf,
        output_path: &PathBuf,
    ) -> eyre::Result<Char> {
        let plan = &conf.plan;

        // TODO: decide whether it is a file or a git repo
        // TODO: Starting with files only, as such implement git repo later

        let path_buf = std::path::PathBuf::from(plan);
        let path = path.join(path_buf);
        if !path.exists() {
            eyre::bail!("path doesn't exist: {}", path.to_string_lossy())
        }
        let path = canonicalize(path)?;

        dbg!(path);
        dbg!(output_path);

        todo!()
    }

    fn create_output_dir(&self, root: &PathBuf) -> eyre::Result<PathBuf> {
        let mut output = root.clone();
        output.push("plans");
        std::fs::create_dir_all(&output)?;

        Ok(output)
    }
}
