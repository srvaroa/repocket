Repocket
========

Repocket is a (raw) tool to download favourited articles on
[Pocket](https://getpocket.com).  Because, you know, The Cloud is not a
reliable storage.  Also, grep.

Building
--------

Rust, openssl and OSX woes:

https://stackoverflow.com/questions/34612395/openssl-crate-fails-compilation-on-mac-os-x-10-11

    $ export OPENSSL_INCLUDE_DIR=$(brew --prefix openssl)/include
    $ export OPENSSL_LIB_DIR=$(brew --prefix openssl)/lib

    $ cargo build

Run
---

    cargo run -- -k $CONSUMER_KEY \
                 -t $ACCESS_TOKEN \
                 -o $TARGET_DIR

The $CONSUMER_KEY is obtained by creating a new application in the
GetPocket web interface at https://getpocket.com/developer/apps/new.
Feel free to use mine: 61972-a63858eb2fdd9aa63e2dcc76.

Once you have it, run:

    ./retrieve_token 61972-a63858eb2fdd9aa63e2dcc76

And follow the instructions.  You'll need to open a link in your browser
(where you are assumed to be logged into Pocket) and authorize the
application (this will redirect to a non existent page, just ignore it).

Now go back to the script and press any key to continue.  You'll see the
$ACCESS_TOKEN on screen.

TODO
----

At the moment each run pulls the entire list of favourites and stores
those that are not present in the target folder.  Ideally this thing
would run on cron, and should use `since` or a similar param to request
new items.

