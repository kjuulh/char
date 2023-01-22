use std::{ops::Deref, sync::Arc};

use crate::{parser::Parser, services::downloader::Downloader};

#[derive(Debug, Clone)]
pub struct Deps {
    inner: Arc<InnerDeps>,
}

#[derive(Debug, Clone)]
pub struct InnerDeps {
    pub parser: Parser,
    pub downloader: Downloader,
}

impl Default for Deps {
    fn default() -> Self {
        let parser = Parser::default();
        let downloader = Downloader::new(parser.clone());

        Self {
            inner: Arc::new(InnerDeps { parser, downloader }),
        }
    }
}

impl Deref for Deps {
    type Target = Arc<InnerDeps>;

    fn deref(&self) -> &Self::Target {
        &self.inner
    }
}
