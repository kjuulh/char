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

    pub fn set_path(&self, path: PathBuf) {
        let mut writer = self.path.write().unwrap();
        *writer = Some(path);
    }

    pub fn get_path(&self) -> eyre::Result<PathBuf> {
        let read_path = self.path.read().unwrap();
        let path = match read_path.clone() {
            Some(p) => p,
            None => todo!(), // find using git later on
        };

        Ok(path)
    }

    pub fn parse(&self) -> eyre::Result<models::Char> {
        let mut path = self.get_path()?;
        if !path.ends_with("char.toml") {
            path.push("char.toml")
        }
        let contents =
            std::fs::read_to_string(&path).context("char.toml doesn't exist at that path")?;

        contents.parse::<Char>()
    }

    pub fn parse_from(&self, path: &PathBuf) -> eyre::Result<models::Char> {
        let contents =
            std::fs::read_to_string(path).context("char.toml doesn't exist at that path")?;

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
