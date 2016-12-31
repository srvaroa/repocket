extern crate hyper;

use std::env;
use hyper::client::Client;
use hyper::header::ContentType;
use std::io::Read;

// TODO: move these to config

fn main() {

    let api = "https://getpocket.com/v3/get";
    let consumer_key = env::var("CONSUMER_KEY").unwrap();
    let access_token = env::var("ACCESS_TOKEN").unwrap();

    let client = Client::new();

    let get_json = format!("{{\"access_token\":\"{}\", \"consumer_key\":\"{}\",\"count\":\"{}\"}}",
                           access_token, consumer_key, 10);

    let mut res = client.post(api)
        .body(&get_json)
        .header(ContentType("application/json".parse().unwrap())) // mime::Mime
        .send()
        .unwrap();
    assert_eq!(res.status, hyper::Ok);

    let mut text = String::new();
    let _ = res.read_to_string(&mut text).unwrap();

    println!("The text is: {}", text);
    assert_eq!(res.status, hyper::Ok);
}
