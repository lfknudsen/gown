package pkg

import (
	"os"
	"path/filepath"
)

const envGownDir = "GOWN_DIR"
const gownDirName = "gown"
const outputDirName = "downloads"
const dataDirName = "data"
const queueDirName = "queue"
const failedDirName = "failed"

func determineGownDir() string {
	gownDir, found := os.LookupEnv(envGownDir)
	if !found {
		INFO("GOWN_DIR environment variable not found.")
		ex, err := os.Executable()
		FATAL(err)
		ExecutableDir := filepath.Dir(ex)
		gownDir = filepath.Join(ExecutableDir, gownDirName)
	}
	INFO("Gown dir is " + gownDir)
	return gownDir
}

var GownDirectory = determineGownDir()
var OutputDir = filepath.Join(GownDirectory, outputDirName)
var DataDir = filepath.Join(GownDirectory, dataDirName)
var QueueDir = filepath.Join(GownDirectory, queueDirName)
var FailedDir = filepath.Join(GownDirectory, failedDirName)

// InitialiseDirectories ensures that the given directory or directories exist by
// creating them if they do not.
func InitialiseDirectories(dir ...string) {
	for _, d := range dir {
		INFO("Ensuring that " + d + " exists.")
		err := os.MkdirAll(d, os.ModeDir)
		FATAL(err)
	}
}

// SetupDirs creates the major subdirectories needed. Individual handlers
// may initialise and use their own if they so wish, but they can rely on
// e.g. OutputDir existing.
func SetupDirs() {
	InitialiseDirectories(OutputDir, QueueDir, DirMega, RepoDir, FailedDir)
}
