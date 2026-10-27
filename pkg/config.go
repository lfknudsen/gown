package main

type Config struct {
	Version string `toml:"version"`
	DownloadDirectory string `toml:"download_dir"`

}
