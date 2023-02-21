use char_sdk_derive::CharAction;

#[derive(CharAction)]
struct Run;

impl char_sdk::Action for Run {
    fn action<'a>(&self) -> char_sdk::ActionArgs<'a> {
        char_sdk::ActionArgs {
            name: "run",
            args: vec![char_sdk::ActionArg::Required {
                name: "profile",
                description: "which release profile to use, takes ['release', 'test']",
            }],
        }
    }

    fn execute(&self, args: &[&char_sdk::InputArg]) -> eyre::Result<()> {
        todo!()
    }
}

#[derive(CharAction)]
struct Build;

impl char_sdk::Action for Build {
    fn action<'a>(&self) -> char_sdk::ActionArgs<'a> {
        todo!()
    }

    fn execute(&self, args: &[&char_sdk::InputArg]) -> eyre::Result<()> {
        todo!()
    }
}

fn main() {
    char_sdk::CharBuilder::new()
        .add_action(Run {})
        .add_action(Build {})
        .add_plugin(char_sdk::std::dagger::Plugin::default())
        .add_plugin(char_sdk::std::k8s::Plugin::default())
        .execute()
        .unwrap();
}
