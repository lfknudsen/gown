# Gown

Download utility. Supports multiple URLs at once. Multithreaded.

If `--overwrite` is supplied, will automatically
overwrite files that already exist with the same name. Otherwise, the new file
will be given a unique name.

## Dependencies
To build, only `Go` is required to be installed.

Cloning Git repositories requires `git`.\
Downloading from DeviantArt requires `gallery-dl`.\
Downloading from YouTube requires `yt-dlp`.

## Build

```
go build
```

## Usage

```
gown [--overwrite] <URLs...>
```

## Supported websites

Gitlab.com, GitHub.com, any website that ends in .git\
Mediafire.com
DeviantArt.com
YouTube.com
