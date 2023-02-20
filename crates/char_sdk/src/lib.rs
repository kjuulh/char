pub mod std;

pub trait Action {}
pub trait Plugin {}

pub struct CharBuilder;

impl CharBuilder {
    pub fn new() -> Self {
        CharBuilder
    }

    pub fn add_context<C>(mut self, context: C) -> Self {
        self
    }

    pub fn add_action(mut self, action: impl Action) -> Self {
        self
    }

    pub fn add_plugin(mut self, plugin: impl Plugin) -> Self {
        self
    }

    pub fn execute(mut self) -> eyre::Result<()> {
        Ok(())
    }
}
