Repocket
========

Repocket is a (raw) tool to download favourited articles on
[Pocket](https://getpocket.com).  Because, you know, The Cloud is not a
reliable storage.  Also, grep.

Repocket will retrieve the list of starred links, download a plain text
copy of the link, and store it in disk.  It works in Linux, might work
in OSX (untested) and will very likely not in Windows (although it
should be trivial to fix.)

Dependencies
------------

Go >= 1.11.

You *MUST* have `links2` installed (links should also work, but you'll
need to create an alias.)  Repocket uses `links2` to render the articles
into plain text.

Exporting your list
-------------------

    make

Repocket will first authenticate against the Pocket API.  It will ask
you to browse to a URL where you can grant permissions to read your
list of articles.  The message looks something like this:

    2019/09/09 20:40:12 Browse to this URL, you may ignore errors:
    https://getpocket.com/auth/authorize?request_token=62074b8c-ed8a-b5e5-71f3-586bcf&redirect_uri=localhost

Click on the link and accept the authorisation.  Once you do this the
first time, you simply need to click the link and ignore the browser.

Once you're done, click enter on the console and watch articles being
downloaded.

You can run the same command several times and it'll skip through
articles that are already downloaded, adding only new ones.

Configuration
-------------

You may (and probably want to) set two environment variables:

- `CONSUMER_KEY` is obtained by creating a new application in the
  GetPocket web interface at https://getpocket.com/developer/apps/new.
  Feel free to use mine: 85480-9793dd8ed508561cb941d987.  This is the
  default value.

- `OUTPUT_DIR` is the directory where all files with downloaded articles
  will be stored.  The directory *must* exist.  The default value is
  ./repocket

For example

    OUTPUT_DIR=/some/other/dir make

Would download articles to `/some/other/dir`.

TODO
----

Some things I'd like to implement:

* Don't download the entire list, using the "since" parameter to
  retrieve new articles from the last run. 
* Store credentials in a file to avoid doing the auth ritual on every
  run.
