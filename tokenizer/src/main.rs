use std::env;
extern crate rustc_serialize;
use rustc_serialize::json;


fn tokenize(pattern: String) -> Vec<String> {
	let mut tokens = vec![];
	let mut token = "".to_string();
	let mut openGroup = false;
	for c in pattern.chars() {
		if openGroup {
			token.push(c);
			if c.to_string() == "]" {
				tokens.push(token.clone());
				openGroup = false;
			}
		}
		else if c.to_string() == "[" {
			openGroup = true;
			token = "[".to_string();
		}
		else {
			tokens.push(c.to_string());
		}
	}
	return tokens;
}

fn main() {
	let args: Vec<String> = env::args().collect();

	let tokens = tokenize(args[1].clone());
	println!("{}", json::encode(&tokens).unwrap());
}
