pub mod std;

pub enum ActionArg<'a> {
    Required { name: &'a str, description: &'a str },
    Optional { name: &'a str, description: &'a str },
}

pub struct InputArg<'a> {
    name: &'a str,
    value: &'a str,
}

pub struct ActionArgs<'a> {
    pub name: &'a str,
    pub args: Vec<ActionArg<'a>>,
}

pub trait Action {
    fn action<'a>(&self) -> ActionArgs<'a>;
    fn execute(&self, args: &[&InputArg]) -> eyre::Result<()>;
}

pub trait Plugin {}

pub struct CharBuilder {
    actions: Vec<Box<dyn Action>>,
}

impl CharBuilder {
    pub fn new() -> Self {
        CharBuilder {
            actions: Vec::new(),
        }
    }

    pub fn add_context<C>(mut self, context: C) -> Self {
        self
    }

    pub fn add_action(mut self, action: impl Into<Box<dyn Action>>) -> Self {
        self.actions.push(action.into());
        self
    }

    pub fn add_plugin(mut self, plugin: impl Plugin) -> Self {
        self
    }

    pub fn execute(mut self) -> eyre::Result<()> {
        //clap::Command::new("").arg(a)
        Ok(())
    }
}
