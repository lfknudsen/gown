package pkg

import (
	"os"
	"path/filepath"
)

const envGownDir = "GOWN_DIR"
const gownDirName = "gown"
const outputDirName = "downloads"
const dataDirName = "data"

func determineGownDir() string {
	gownDir := os.Getenv(envGownDir)
	if gownDir == "" {
		ex, err := os.Executable()
		FATAL(err)
		ExecutableDir := filepath.Dir(ex)
		return filepath.Join(ExecutableDir, gownDirName)
	}
	return gownDir
}

var GownDirectory = determineGownDir()
var OutputDir = filepath.Join(GownDirectory, outputDirName)
var DataDir = filepath.Join(GownDirectory, dataDirName)

// InitialiseDirectories ensures that the given directory or directories exist by
// creating them if they do not.
func InitialiseDirectories(dir ...string) {
	for _, d := range dir {
		err := os.MkdirAll(d, os.ModeDir)
		FATAL(err)
	}
}

// SetupDirs creates the major subdirectories needed. Individual handlers
// may initialise and use their own if they so wish, but they can rely on
// e.g. OutputDir existing.
func SetupDirs() {
	InitialiseDirectories(OutputDir, QueueDir, DirMega, RepoDir)
}
