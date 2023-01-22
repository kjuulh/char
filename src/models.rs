use std::collections::BTreeMap;

use serde::{Deserialize, Serialize};

type Overrides = BTreeMap<String, String>;
type Dependencies = Vec<String>;

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Conf {
    pub plan: String,
    pub dependencies: Option<Dependencies>,
    pub overrides: Option<Overrides>,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Application {
    name: String,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Plan {
    name: String,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
#[serde(untagged)]
pub enum Char {
    Application {
        char: Option<Conf>,
        application: Application,
        config: BTreeMap<String, toml::Value>,
    },
    Plan {
        char: Option<Conf>,
        plan: Plan,
        config: BTreeMap<String, toml::Value>,
    },
}

impl std::str::FromStr for Char {
    type Err = eyre::Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let t: Char = toml::from_str(s)?;

        Ok(t)
    }
}
