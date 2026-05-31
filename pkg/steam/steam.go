package steam

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	cs2AppID = "730"
)

func GetSteamPath() (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		key, err = registry.OpenKey(registry.LOCAL_MACHINE, `Software\Valve\Steam`, registry.QUERY_VALUE)
		if err != nil {
			return "", fmt.Errorf("Steam not found in registry")
		}
	}
	defer key.Close()

	steamPath, _, err := key.GetStringValue("SteamPath")
	if err != nil {
		return "", fmt.Errorf("failed to read Steam path: %w", err)
	}

	return filepath.FromSlash(steamPath), nil
}

func GetCS2InstallPath() (string, error) {
	steamPath, err := GetSteamPath()
	if err != nil {
		return "", err
	}

	// Check default library
	defaultPath := filepath.Join(steamPath, "steamapps", "common", "Counter-Strike Global Offensive")
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath, nil
	}

	// Check library folders
	libraryFoldersPath := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
	data, err := os.ReadFile(libraryFoldersPath)
	if err != nil {
		return "", fmt.Errorf("CS2 installation not found")
	}

	// Parse VDF to find library paths
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "\"path\"") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 4 {
				libPath := filepath.FromSlash(parts[3])
				cs2Path := filepath.Join(libPath, "steamapps", "common", "Counter-Strike Global Offensive")
				if _, err := os.Stat(cs2Path); err == nil {
					return cs2Path, nil
				}
			}
		}
	}

	return "", fmt.Errorf("CS2 installation not found in any Steam library")
}

func GetCS2ConfigPath() (string, error) {
	cs2Path, err := GetCS2InstallPath()
	if err != nil {
		return "", err
	}

	// CS2 config is in game/csgo/cfg
	configPath := filepath.Join(cs2Path, "game", "csgo", "cfg")
	if _, err := os.Stat(configPath); err != nil {
		return "", fmt.Errorf("CS2 config directory not found: %w", err)
	}

	return configPath, nil
}

func GetUserDataPath() (string, error) {
	steamPath, err := GetSteamPath()
	if err != nil {
		return "", err
	}

	// Get Steam user ID (use most recent)
	userdataPath := filepath.Join(steamPath, "userdata")
	entries, err := os.ReadDir(userdataPath)
	if err != nil {
		return "", fmt.Errorf("failed to read userdata: %w", err)
	}

	var latestUserID string
	var latestTime int64
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "0" && entry.Name() != "ac" {
			info, err := entry.Info()
			if err == nil && info.ModTime().Unix() > latestTime {
				latestTime = info.ModTime().Unix()
				latestUserID = entry.Name()
			}
		}
	}

	if latestUserID == "" {
		return "", fmt.Errorf("no Steam user found")
	}

	cs2ConfigPath := filepath.Join(userdataPath, latestUserID, cs2AppID, "local", "cfg")
	return cs2ConfigPath, nil
}
