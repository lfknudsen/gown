# Gown

Download utility. Supports multiple URLs at once. Concurrent.

For most websites, the process is delegated to specialised programmes
which will need to be installed.

## Dependencies
To build, only `Go` is required to be installed.

Cloning Git repositories requires `git`.\
Downloading from DeviantArt (and other gallery websites) requires `gallery-dl`.\
Downloading from YouTube (and other video websites) requires `yt-dlp`.\
Downloading from Mega requires `megatools`.

## Build

```
go build
```

## Usage

```
gown [OPTIONS] <URLs...>

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

Whether you prefix arguments with one or two dashes is irrelevant.

## Supported websites

GitLab.com, GitHub.com, any website that ends in .git\
Mediafire.com\
DeviantArt, and anything else supported by `gallery-dl`\
YouTube, and anything else supported by `yt-dlp`
