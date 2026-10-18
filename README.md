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
Downloading from Mega requires `megatools`.

## Build

The project can build three separate executables, though all functionality can
be reached with the main `gown`, located in `cmd\gown`:

```
go build -C cmd/gown [-o <output file-path>]
```

To build the other utilities:
```
go build -C cmd/queue [-o <output file-path>]
go build -C cmd/gown-srv [-o <output file-path>]
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

gown [OPTIONS] interactive
gown-srv [OPTIONS]
            Launches an interactive shell which allows you to repeatedly
            input and download URLs.

Options:
-h | --help         Print help message and exit.

-overwrite          If a file with the name already exists, delete it. If this is
                    false, the new file will be given a unique name instead.
                    Default: false.

-generate-lists     Generates/re-generates the lists which contain the gallery
                    and video websites supported for download. If your version
                    of `gallery-dl` or `yt-dlp` is newer than when you last
                    used this programme to download something using one of those,
                    include this argument to update the lists, so that this
                    programme can correctly delegate URLs to those services.
                    Default: false.
```

Whether you prefix options with one or two dashes is irrelevant.

## Supported websites

GitLab.com, GitHub.com, any website that ends in .git\
Mediafire.com\
DeviantArt, and anything else supported by `gallery-dl`\
YouTube, and anything else supported by `yt-dlp`\
Mega.nz
