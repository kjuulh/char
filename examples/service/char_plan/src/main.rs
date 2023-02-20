struct Run;
impl char_sdk::Action for Run {}

struct Build;
impl char_sdk::Action for Build {}

fn main() {
    char_sdk::CharBuilder::new()
        .add_context(char_sdk::std::dagger::Context::default())
        .add_action(Run {})
        .add_action(Build {})
        .add_plugin(char_sdk::std::k8s::Plugin::default())
        .execute()
        .unwrap();
}
