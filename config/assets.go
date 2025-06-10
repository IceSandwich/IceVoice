package config

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

func getAppDataDir() string {
	var appDataPath string
	switch runtime.GOOS {
	case "windows":
		// Windows 使用 APPDATA 环境变量
		appDataPath = os.Getenv("APPDATA")
		if appDataPath == "" {
			log.Fatal("Error getting current user: APPDATA environment variable is not set.")
		}
	case "darwin":
		// macOS 使用 ~/Library/Application Support
		user, err := user.Current()
		if err != nil {
			log.Fatal("Error getting current user:", err)
		}
		appDataPath = filepath.Join(user.HomeDir, "Library", "Application Support")
	case "linux":
		// Linux 使用 ~/.config
		user, err := user.Current()
		if err != nil {
			log.Fatal("Error getting current user:", err)
		}
		appDataPath = filepath.Join(user.HomeDir, ".config")
	default:
		log.Fatal("Cannot get appdata, unsupported OS.")
	}

	return appDataPath
}

var appDataPath string

func GetManifestsDir() string {
	return filepath.Join(appDataPath, "IceVoice", "models", "manifest")
}

func GetBlobsDir() string {
	return filepath.Join(appDataPath, "IceVoice", "models", "blobs")
}

func HashToBlobFile(hash string) string {
	hash = strings.TrimPrefix(hash, "sha256:")

	return filepath.Join(GetBlobsDir(), hash)
}

func init() {
	appDataPath = getAppDataDir()

	for _, d := range []string{
		GetManifestsDir(),
		GetBlobsDir(),
	} {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			log.Fatal("Error creating app data dir:", d)
		}
	}
}
