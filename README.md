# Gown

Download utility. Supports multiple URLs at once. Concurrent.

For most websites, the process is delegated to specialised programmes
which will need to be installed.

## Dependencies
To build, only `Go` is required to be installed.

The following dependencies are only required if you wish to download from the
websites they support. You will need the executables/scripts to be on your PATH:

Cloning Git repositories requires `git`.\
Downloading from DeviantArt (and other gallery websites) requires `gallery-dl`.\
Downloading from YouTube (and other video websites) requires `yt-dlp`.\
Downloading from Mega requires `megatools`.\
Downloading from Instagram requires `instaloader`\
Downloading from Google Drive requires `gdown`

## Build

The project can build several separate executables, though all functionality can
be reached with the main `gown`, located in `cmd\gown`:

```
go build -C cmd/gown [-o <output directory or file>]
```

To build the other utilities:
```
go build -C cmd/queue [-o <output directory or file>]
go build -C cmd/gown-srv [-o <output directory or file>]
go build -C cmd/gown-client [-o <output directory or file>]
go build -C cmd/worker [-o <output directory or file>]
```

A Makefile exists to ease the build process. Some of the recipes (including `all`)
rely on my "[install](https://gitlab.com/lfknudsen/install)" programme.

```
make
```

## Usage
```
gown [OPTIONS] <URLs...>
            Downloads the given URL(s).

gown [OPTIONS] queue <URLs...>
queue [OPTIONS] <URLs...>
            Saves the URL(s) to timestamped files in the 'gown/queue' directory.

gown queue download [<number>]
queue download [<number>]
            Downloads the files in the queue. If <number> is included, will download
            the first <number> of entries in the queue.

gown queue list
queue list
            Prints the URLs in the queue to the standard output.

gown [OPTIONS] interactive
gown-srv [OPTIONS]
            Launches an interactive shell which allows you to repeatedly
            input and download URLs.
            Exit by entering (without quotiation marks) "exit" or "quit".

gown [OPTIONS] server
gown-srv [OPTIONS]
            Launches a GRPC server which will allow clients to request it to
            perform download of URLs.

gown [OPTIONS] client <URL(s)...>
gown-client [OPTIONS] <URL(s)...>
            Requests a GRPC server to download the given URL(s).

Options:
-h, --help                    Print help message and exit.

-o, --overwrite               If a file with the name already exists, delete it. If this is
                              false, the new file will be given a unique name instead.
                              Default: false.

-g, --generate-lists          Generates/re-generates the lists which contain the gallery
                              and video websites supported for download. If your version
                              of `gallery-dl` or `yt-dlp` is newer than when you last
                              used this programme to download something using one of those,
                              include this argument to update the lists, so that this
                              programme can correctly delegate URLs to those services.
                              Default: false.

-s, --sequential              Perform downloads sequentially (by default multiple URLs are
                              downloaded concurrently).

-v, --verbose <integer>       Verbosity level. Higher = more info is printed to console.
                              0 is silent. Default: 1.

-r, --retries <integer>       The number of retries to attempt when downloading through
                              gallery-dl. Less than 0 = infinite. Default: 5.

-p, --port <integer>          The port to connect/listen to when in server or client mode.
                              Default: 8080
```

## Supported websites

GitLab.com, GitHub.com, any website that ends in .git\
Mediafire.com\
DeviantArt, and anything else supported by `gallery-dl`\
YouTube, and anything else supported by `yt-dlp`\
Mega.nz\
Instagram.com (individual posts)\
Google Drive
