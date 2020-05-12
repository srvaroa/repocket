Repocket
========

**Repocket** is a (very, very raw) tool to manage a local copy of articles
from [GetPocket](https://getpocket.com).  Because, you know, The
Cloud is not a reliable storage.  Also, grep.  Local copies contain a
header in plain YAML with the [metadata provided by
Pocket](https://getpocket.com/developer/docs/v3/retrieve)

**Repocket** maintains a simple folder structure to represent states.
Articles are rendered as text files into two folders for unread and
favourite articles.  You can move articles to a deleted and archived
folders to sync those statuses back to Pocket.

**Caveat emptor**

* Local copies are plain text, rendered with
  [w3m](http://w3m.sourceforge.net/), so expect images to be lost.  I'm
  generally OK with this, and the original URLs are still there.
* This workflow is not meant to be particularly user friendly, it just
  fits me well.  I use the "favourite" state to mark articles that I
  want to keep a copy for.
* It works in Linux and OSX, and will very likely not in Windows
  (although it should be trivial to fix.)

Building
--------

Go >= 1.11.

You *MUST* have [w3m](http://w3m.sourceforge.net/) installed. Repocket
uses [w3m](http://w3m.sourceforge.net/) to render the articles into
plain text.

Configuration
-------------

You must create the config file before running `Repocket`.  It's a very
simple yaml stored at `~/.config/repocket/config` with these contents:

    consumer_key: 85480-9793dd8ed508561cb941d987
    favs_dir: <absolute_path_dir_for_favourited_articles>
    unread_dir: <absolute_path_dir_for_unread_articles>
    deleted_dir: <absolute_path_dir_for_deleted_articles>
    archived_dir: <absolute_path_dir_for_archived_articles>

* The `consumer_key` indicates the GetPocket Application.  You can leave
  the sample value above (my own app), but if you prefer to use your own
  just [create a new
  application](https://getpocket.com/developer/apps/new).

The rest are directories to store articles.  All expect an absolute
path.  Two of these are synced up and downstream:

* `unread_dir` is the directory where unread, non archived articles will
  be downloaded.  If a downloaded article is marked as archived in
  Pocket, then the file will be deleted.
* `favs_dir` contains favourited articles.  You move articles here from
  the `unread_dir`, in the next sync they will be marked as favourite
  and archived.  Articles in your local fav directory are *never*
  deleted (even if you unfav them upstream.)

These directories *only sync upstream*:

* `deleted_dir` contains articles that should be deleted in the next run
  of the `delete` command.  You move articles here when you want to
  delete them from Pocket.
* `archived_dir` contains articles that should be archived in the next run
  of the `archive` command.  You move articles here when you want to
  archive them in Pocket, but don't want to keep a local copy.

Files in both `deleted_dir` and `archived_dir` are removed after a sync.

Set up
------

When you first run `Repocket`, it will authenticate against the Pocket
API.  It will ask you to browse to a URL where you can grant permissions
to read your list n articles.  The message looks something like this:

    2019/09/09 20:40:12 Browse to this URL, you may ignore errors:
    https://getpocket.com/auth/authorize?request_token=62074b8c-ed8a-b5e5-71f3-586bcf&redirect_uri=localhost

Click on the link and accept the authorisation.  Once you do this the
first time, you simply need to click the link and ignore the browser.

This step will write a new `access_token` property to your config file
so you don't need to auth again.

Synchronizing
-------------

Run

    GO111MODULE=on go run ./cmd/repocket [favs|delete|archive|unread|sync]

And **Repocket** will sync the folder associated to the given action.

* `delete`: will delete articles added to `favs_dir` as deleted in
  Pocket. 
* `archive`: will delete articles added to `favs_dir` as deleted in
  Pocket. 
* `favs`: will mark articles added to `favs_dir` as favourited in
  Pocket.
* `unread`: will download all unread articles to `unread_dir`.
* `sync`: will execute the previous three actions, in the same order as
  shown in this list.

TODO
----

Some things I'd like to implement:

* Prepend the source URL to the file.
