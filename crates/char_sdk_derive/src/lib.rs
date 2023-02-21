extern crate proc_macro;

use proc_macro::TokenStream;
use quote::quote;

#[proc_macro_derive(CharAction)]
pub fn derive_char_action_into_box(input: TokenStream) -> TokenStream {
    let ast = syn::parse(input).unwrap();

    // Build the trait implementation
    impl_hello_macro(&ast)
}

fn impl_hello_macro(ast: &syn::DeriveInput) -> TokenStream {
    let name = &ast.ident;
    let gen = quote! {
        impl Into<Box<dyn char_sdk::Action>> for #name {
            fn into(self) -> Box<dyn char_sdk::Action> {
                Box::new(self)
            }
        }
    };
    gen.into()
}
