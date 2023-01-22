use std::{
    path::PathBuf,
    sync::{Arc, RwLock},
};

use eyre::Context;

use crate::models::{self, Char};

#[derive(Debug, Clone)]
pub struct Parser {
    path: Arc<RwLock<Option<PathBuf>>>,
}

impl Parser {
    pub fn new(path: PathBuf) -> Self {
        Self {
            path: Arc::new(RwLock::new(Some(path))),
        }
    }

    pub fn set_path(&self, _path: PathBuf) {
        let _writer = self.path.write().unwrap();
    }

    pub fn parse(&self) -> eyre::Result<models::Char> {
        let read_path = self.path.read().unwrap();
        let path = match read_path.clone() {
            Some(p) => p,
            None => todo!(), // find using git later on
        };

        let contents =
            std::fs::read_to_string(&path).context("char.toml doesn't exist at that path")?;

        contents.parse::<Char>()
    }
}

impl Default for Parser {
    fn default() -> Self {
        Self {
            path: Arc::new(RwLock::new(None)),
        }
    }
}
