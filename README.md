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
simple yaml stored at `~/.repocket/config` with these contents:

    consumer_key: 85480-9793dd8ed508561cb941d987
    favs_dir: <absolute_path_dir_for_favourited_links>
    unread_dir: <absolute_path_dir_for_unread_links>

The `consumer_key` indicates the GetPocket Application.  You can leave
the sample value above (my own app), but of you prefer to use your own
just [create a new application](https://getpocket.com/developer/apps/new).

The `favs_dir` is the directory where favourited articles will be
downloaded.  Expects an absolute path.

The `unread_dir` is the directory where unread, non archived articles
will be downloaded.  Expects an absolute path. 

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

Run

    GO111MODULE=on go run ./cmd/repocket [favs|unread|archive]

Use `favs` to download only favourited articles, `unread` for queued
ones.

You can run any command several times and it'll skip through
articles that are already downloaded, adding only new ones.

TODO
----

Some things I'd like to implement:

* Paginate and download only changes since the last sync.  Not sure if
  this plays well with the point below.
* I'm converging to the idea that an ideal workflow is to sync both my
  favs and unread folders regularly and move files from unread -> favs
  whenever I want to keep it.  When this happens, I want a additional
  sync back to GetPocket (e.g. mark the item as favourited).  But: how
  does a deletion work?
* Prepend the source URL to the file.
