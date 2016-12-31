extern crate hyper;
extern crate serde_json;

use hyper::client::Client;
use hyper::header::ContentType;
use std::env;
use std::io::Read;
use serde_json::Value;

fn main() {

    let api = "https://getpocket.com/v3/get";
    let consumer_key = env::var("CONSUMER_KEY").unwrap();
    let access_token = env::var("ACCESS_TOKEN").unwrap();

    let client = Client::new();

    let get_json = format!(
            "{{\"access_token\":\"{}\",
            \"consumer_key\":\"{}\",
            \"count\":\"{}\",
            \"favorite\":1}}",
            access_token, consumer_key, 10);

    let mut res = client.post(api)
        .body(&get_json)
        .header(ContentType("application/json".parse().unwrap())) // mime::Mime
        .send()
        .unwrap();
    assert_eq!(res.status, hyper::Ok);

    let mut raw_json = String::new();
    let _ = res.read_to_string(&mut raw_json).unwrap();

    let json: Value = serde_json::from_str(&raw_json).unwrap();
    if let &Value::Object(ref o) = json.find_path(&["list"]).unwrap() {
        for (item_id, vals) in o {
            on_item(item_id, vals);
        }
    }

    assert_eq!(res.status, hyper::Ok);
}

fn on_item(item_id: &String, val: &Value) {
    let url = val.find_path(&["resolved_url"]).unwrap();
    println!("{} {:?}", item_id, url);
}
