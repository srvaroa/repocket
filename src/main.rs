extern crate getopts;
extern crate hyper;
extern crate serde_json;

use getopts::Options;
use hyper::client::Client;
use hyper::client::Response;
use hyper::header::ContentType;
use serde_json::Value;
use std::env;
use std::fs::File;
use std::fs;
use std::io::Read;
use std::io::{Write, BufWriter};
use std::path::Path;
use std::process::Command;

struct Config {
    api: String,
    consumer_key: String,
    access_token: String,
    output_dir: String,
}

fn main() {

    let args: Vec<String> = env::args().collect();
    let mut opts = Options::new();
    opts.reqopt("k", "consumer_key", "Pocket's consumer key", "NAME");
    opts.reqopt("t", "access_token", "Pocket's access token", "NAME");
    opts.reqopt("o", "base_dir", "Base directory to store articles", "NAME");
    let matches = match opts.parse(&args[1..]) {
        Ok(m) => { m }
        Err(f) => { panic!(f.to_string()) }
    };

    let cfg = Config{
        api: String::from("https://getpocket.com/v3/get"),
        consumer_key: matches.opt_str("k").unwrap(),
        access_token: matches.opt_str("t").unwrap(),
        output_dir: matches.opt_str("o").unwrap(),
    };

    println!("Ensuring that {} exists..", cfg.output_dir);
    fs::create_dir_all(&cfg.output_dir).expect("Failed to create output dir");

    println!("Querying Pocket..");
    let client = Client::new();
    let mut res = query_favourites(&client, &cfg);

    println!("Processing results..");
    process_batch(&mut res, &cfg);

    println!("Done!");
}

fn query_favourites(client: &Client, cfg: &Config) -> Response {
    let get_json = format!(
            "{{\"access_token\":\"{}\",
            \"consumer_key\":\"{}\",
            \"favorite\":1}}",
            cfg.access_token, cfg.consumer_key);

    /*
     * Add:
     *
     *  \"count\":\"{}\",
     *
     * to limit the number of items pulled.  a Since or similar would make more
     * sense.
     */

    let res = client.post(&cfg.api)
        .body(&get_json)
        .header(ContentType("application/json".parse().unwrap())) // mime::Mime
        .send()
        .unwrap();

    assert_eq!(res.status, hyper::Ok);
    return res;
}

fn process_batch(res: &mut Response, cfg: &Config) {
    let mut raw_json = String::new();
    let _ = res.read_to_string(&mut raw_json).unwrap();
    let json: Value = serde_json::from_str(&raw_json).unwrap();
    if let &Value::Object(ref o) = json.find_path(&["list"]).unwrap() {
        for (item_id, vals) in o {
            store(item_id, vals, &cfg.output_dir);
        }
    }
}

fn store(item_id: &String, val: &Value, output_dir: &str) {

    let url = val.find_path(&["resolved_url"]).unwrap();
    let title = val.find_path(&["given_title"]).unwrap().as_str().unwrap();

    let file_name = clean_title(&title);
    let out_path = Path::new(output_dir).join(&file_name);
    if out_path.exists() {
        println!("File already exists {:?}", out_path);
        return;
    }

    let f = File::create(out_path).expect("Unable to open output file");
    let out = Command::new("links")
        .arg("-dump")
        .arg(url.as_str().unwrap())
        .output()
        .expect("Unable to dump url");

    let mut f = BufWriter::new(f);
    f.write_all(&out.stdout).expect("Unable to write data");

    println!("Saved {} {} at {}/{:?}", item_id, url, output_dir, file_name);

}

fn clean_title(t: &str) -> String {
    return str::replace(t, " ", "_")
        .replace(":", "_")
        .replace(",", "_")
        .replace("|", "_by_")
        .replace("/", "_")
}
