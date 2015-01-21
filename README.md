# GDCP

All I wanted was `scp` for Google Drive.  What I got is a lesson on
why no one else had written it yet.

### Some Assembly Required

The code you see here doesn't actually work.  It seems there's no good
way to ship an OAuth 2.0 Client ID and Client Secret in open source
software [1], so I've kept my secrets to myself.  If you want this
software to work, you'll need to mint your own by following the
instructions laid out in [2]. It's inconvenient, but it should be
simple enough for the kind of person that prefers command-line apps.

[1] http://stackoverflow.com/questions/27585412/can-i-really-not-ship-open-source-with-client-id<br/>
(The question was unanswered at the time of this writing.)<br/>
[2] https://developers.google.com/accounts/docs/OAuth2<br/>

### Sharp Edges

It turns out Google Drive doesn't work like a regular filesystem.
It's possible to have two files in the same folder with the exact same
name.  It's also possible to have multiple versions of the same file.

This leads to some ambiguity.  When the user tries to copy over a file
that already exists, should the client:

1. Report an error and exit?
2. Make a new file with the exact same name?
3. Overwrite the existing file?
4. Modify the existing file, which is like option (c) but it
   leaves history (and space consumption) intact?

I have no idea how to write a command-line UI that makes sense of
this.

It get even worse when you want to download files from Google Drive.
In that case, when you encounter two or more files with the same name,
how do you know which one the user wants?

It all makes sense behind the scenes, since everything has a unique
ID.  The web UI can make sense of it, too.  But it's hard to translate
this to the UNIX filesystem paradigm.

### Lowered Expectations

Given these limitations, I don't know what a good Google Drive command
line client would look like.  This certainly isn't it.  I abandoned
hope and decided to hack together something that would support my use
case while putting in the minimum amount of effort possible.

My use case, by the way, is to have automated backups of a single file
to Google Drive every night.

Here's how it works:

```
Usage: gdcp [options] SOURCE DEST
  -allowDupes=false: Allow this client to create files whose names shadow existing files.
  -httptest.serve="": if non-empty, httptest.NewServer serves on this address and blocks
  -keepHistory=true: Whether or not file history is preserved.
  -update=false: If a file of the same name already exists, update it.
```

(You can safely ignore that httptest flag.  It's a leftover from the
code's origins.)

By default, it will complain and do nothing if you try to copy over a
file that already exists.  The flags allow you to override that
behaviour.

Downloading files from Google Drive is not supported.  Folders are
probably not supported either; I haven't tried.

### Authentication and Security

The first time you use it, it should try to open up a web browser and
ask you to authorize the app.  You may be asked to provide your Google
login credentials.  I know that's a bit odd for a command line app,
but it appears that there's no way around it.

Once that's done, there will be an auth token stored in your OS's
cache directory.  (Mine's at $HOME/.cache/go-api-demo-tok12345).  That
token allows access to your Google Drive account, so try not to lose
it or accidentally upload it anywhere public.

### Humble Origins

In making this hack, I took copious amounts of code from the
Google API Go Client libraries.  [1]

Whatever code remains from that project retains its original license.
The code that I've added is offered under the same permissive
license, so you're free to remix this code even further.

I should be clear about one thing, though: this is not a Google
project, and it's not in any way supported by them.  It's not supported
by me, either.

[1] https://github.com/google/google-api-go-client<br/>
(Previously at code.google.com, now on GitHub)
