use crate::models::Char;

pub struct Context {}

impl Context {
    pub fn new(chars: Vec<Char>) -> Self {
        Self {}
    }
}

pub struct ContextBuilder {
    chars: Vec<Char>,
}

impl ContextBuilder {
    pub fn new() -> Self {
        Self { chars: vec![] }
    }

    pub fn add_char(mut self, char: &Char) -> Self {
        self.chars.push(char.clone());

        self
    }

    pub fn build(self) -> Context {
        Context::new(self.chars)
    }
}
