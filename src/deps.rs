use std::{
    ops::{Deref},
    sync::Arc,
};

use crate::parser::Parser;

#[derive(Debug, Clone)]
pub struct Deps {
    inner: Arc<InnerDeps>,
}

#[derive(Debug, Clone)]
pub struct InnerDeps {
    pub parser: Parser,
}

impl Default for Deps {
    fn default() -> Self {
        Self {
            inner: Arc::new(InnerDeps {
                parser: Parser::default(),
            }),
        }
    }
}

impl Deref for Deps {
    type Target = Arc<InnerDeps>;

    fn deref(&self) -> &Self::Target {
        &self.inner
    }
}
