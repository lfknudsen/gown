# Gown

Download utility. Supports multiple URLs at once. Multithreaded.

If `--overwrite` is supplied, will automatically
overwrite files that already exist with the same name. Otherwise, the new file
will be given a unique name.

## Dependencies
To build, only `Go` is required to be installed.

To support cloning Git repositories, `git` needs to be installed and on
your PATH.

To support downloading from DeviantArt, `gallery-dl` needs to be installed and
on your PATH.

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
