struct Run;
impl char::Action for Run {}

struct Build;
impl char::Action for Build {}

fn main() {
    char::new()
        .add_context(char::dagger::Context::default())
        .add_action(Run {})
        .add_action(Build {})
        .add_plugin(char::std::k8s::Context::default())
        .execute();
}
