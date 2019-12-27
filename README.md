Repocket
========

Repocket is a (very, very raw) tool to download favourited articles on
[Pocket](https://getpocket.com).  Because, you know, The Cloud is not a
reliable storage.  Also, grep.

Repocket has some basic operations:

* `dump`: to retrieve the list of starred links, download a plain text
  copy of the link, and store it in disk.  It works in Linux, might work
  in OSX (untested) and will very likely not in Windows (although it
  should be trivial to fix.)
* `list`: to retrieve the list of favourited articles (title, url)
* `next`: to dump the contents of the newest article in your list

Dependencies
------------

Go >= 1.11.

You *MUST* have [w3m](http://w3m.sourceforge.net/) installed. Repocket
uses [w3m](http://w3m.sourceforge.net/) to render the articles into
plain text.

Configuration
-------------

You must create the config file before running `Repocket`.  It's a very
simple yaml with these contents:

    consumer_key: 85480-9793dd8ed508561cb941d987
    output_dir: <target_directory_for_dump>

The `consumer_key` indicates the GetPocket Application.  You can leave
the sample value above (my own app), but of you prefer to use your own
just [create a new application](https://getpocket.com/developer/apps/new).

The `output_dir` is the directory where all files with downloaded
articles will be stored.  If empty, it defaults to `./repocket`.

When you first run `Repocket`, it will authenticate against the Pocket
API.  It will ask you to browse to a URL where you can grant permissions
to read your list of articles.  The message looks something like this:

    2019/09/09 20:40:12 Browse to this URL, you may ignore errors:
    https://getpocket.com/auth/authorize?request_token=62074b8c-ed8a-b5e5-71f3-586bcf&redirect_uri=localhost

Click on the link and accept the authorisation.  Once you do this the
first time, you simply need to click the link and ignore the browser.

This step will write a new `access_token` property to your config file
so you don't need to auth again.

Exporting your list
-------------------

Once you're done, click enter on the console and watch articles being
downloaded.

You can run the same command several times and it'll skip through
articles that are already downloaded, adding only new ones.

TODO
----

Some things I'd like to implement:

* Don't download the entire list, using the "since" parameter to
  retrieve new articles from the last run. 
