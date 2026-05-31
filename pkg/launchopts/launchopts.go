package launchopts

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"cs2-boost/pkg/backup"
	"cs2-boost/pkg/steam"
)

func getOptimizedLaunchOptions() string {
	// Auto-detect monitor refresh rate and CPU cores
	cores := runtime.NumCPU()
	
	// Use physical cores (assume hyperthreading)
	physicalCores := cores / 2
	if physicalCores < 2 {
		physicalCores = cores
	}

	return fmt.Sprintf("-high -nojoy -novid +fps_max 0 +cl_forcepreload 1 -threads %d +exec autoexec", physicalCores)
}

func GetLocalConfigPath() (string, error) {
	steamPath, err := steam.GetSteamPath()
	if err != nil {
		return "", err
	}

	localConfigPath := filepath.Join(steamPath, "userdata")
	entries, err := os.ReadDir(localConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to read userdata: %w", err)
	}

	// Find most recent user
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

	return filepath.Join(localConfigPath, latestUserID, "config", "localconfig.vdf"), nil
}

func ApplyLaunchOptions() error {
	// Check if Steam is running
	if err := steam.WarnIfSteamRunning(); err != nil {
		return err
	}

	steamPath, err := steam.GetSteamPath()
	if err != nil {
		return fmt.Errorf("failed to locate Steam: %w", err)
	}

	userdataPath := filepath.Join(steamPath, "userdata")
	entries, err := os.ReadDir(userdataPath)
	if err != nil {
		return fmt.Errorf("failed to read userdata: %w", err)
	}

	launchOpts := getOptimizedLaunchOptions()
	successCount := 0
	var lastError error

	// Apply to all Steam accounts
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "0" || entry.Name() == "ac" {
			continue
		}

		configPath := filepath.Join(userdataPath, entry.Name(), "config", "localconfig.vdf")
		
		// Check if config exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}

		fmt.Printf("\n      User %s: %s", entry.Name(), configPath)

		content, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf(" - SKIP (read error)")
			lastError = err
			continue
		}

		contentStr := string(content)

		// Check if CS2 exists in this user's library
		if !strings.Contains(contentStr, "\"730\"") {
			fmt.Printf(" - SKIP (no CS2)")
			continue
		}

		// Remove ALL existing LaunchOptions entries
		launchOptRe := regexp.MustCompile(`(?m)^\s*"LaunchOptions"\s+"[^"]*"\s*$`)
		matches := launchOptRe.FindAllStringIndex(contentStr, -1)
		
		for i := len(matches) - 1; i >= 0; i-- {
			match := matches[i]
			lineStart := match[0]
			lineEnd := match[1]
			
			if lineEnd < len(contentStr) && contentStr[lineEnd] == '\n' {
				lineEnd++
			} else if lineEnd < len(contentStr) && contentStr[lineEnd] == '\r' {
				lineEnd++
				if lineEnd < len(contentStr) && contentStr[lineEnd] == '\n' {
					lineEnd++
				}
			}
			
			contentStr = contentStr[:lineStart] + contentStr[lineEnd:]
		}
		
		// Add our LaunchOptions right after the "730" { line
		cs2StartRe := regexp.MustCompile(`("730"\s*\n\s*\{)\s*\n`)
		if !cs2StartRe.MatchString(contentStr) {
			fmt.Printf(" - SKIP (malformed)")
			continue
		}
		
		contentStr = cs2StartRe.ReplaceAllString(contentStr, "${1}\n\t\t\t\t\"LaunchOptions\"\t\t\""+launchOpts+"\"\n")

		// Write back
		if err := os.WriteFile(configPath, []byte(contentStr), 0644); err != nil {
			fmt.Printf(" - FAILED")
			lastError = err
			continue
		}

		fmt.Printf(" - OK")
		successCount++
	}

	fmt.Printf("\n      Applied to %d account(s)\n      ", successCount)

	if successCount == 0 && lastError != nil {
		return lastError
	}

	return nil
}

func GetCurrentLaunchOptions() (string, error) {
	configPath, err := GetLocalConfigPath()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`"730"\s*\{[^}]*"LaunchOptions"\s*"([^"]*)"`)
	matches := re.FindStringSubmatch(string(content))
	
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", nil
}

func RestoreLaunchOptions(opts string) error {
	configPath, err := GetLocalConfigPath()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	re := regexp.MustCompile(`("730"\s*\{[^}]*"LaunchOptions"\s*")([^"]*)("`)
	
	if re.MatchString(contentStr) {
		contentStr = re.ReplaceAllString(contentStr, `${1}`+opts+`${3}`)
		return os.WriteFile(configPath, []byte(contentStr), 0644)
	}

	return nil
}

func GetStatus() (bool, error) {
	current, err := GetCurrentLaunchOptions()
	if err != nil {
		return false, nil
	}

	return strings.Contains(current, "-high") && strings.Contains(current, "-nojoy"), nil
}

func BackupCurrent() *backup.LaunchOptsBackup {
	steamPath, err := steam.GetSteamPath()
	if err != nil {
		return &backup.LaunchOptsBackup{
			Accounts: make(map[string]string),
		}
	}

	userdataPath := filepath.Join(steamPath, "userdata")
	entries, err := os.ReadDir(userdataPath)
	if err != nil {
		return &backup.LaunchOptsBackup{
			Accounts: make(map[string]string),
		}
	}

	accounts := make(map[string]string)

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "0" || entry.Name() == "ac" {
			continue
		}

		configPath := filepath.Join(userdataPath, entry.Name(), "config", "localconfig.vdf")
		
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}

		content, err := os.ReadFile(configPath)
		if err != nil {
			continue
		}

		contentStr := string(content)
		
		// Check if CS2 exists
		if !strings.Contains(contentStr, "\"730\"") {
			continue
		}

		// Get current launch options
		re := regexp.MustCompile(`"730"\s*\{[^}]*"LaunchOptions"\s*"([^"]*)"`)
		matches := re.FindStringSubmatch(contentStr)
		
		if len(matches) > 1 {
			accounts[entry.Name()] = matches[1]
		} else {
			accounts[entry.Name()] = "" // No launch options set
		}
	}

	return &backup.LaunchOptsBackup{
		Accounts: accounts,
	}
}

func RestoreFromBackup(b *backup.LaunchOptsBackup) error {
	if b == nil || len(b.Accounts) == 0 {
		return nil
	}

	steamPath, err := steam.GetSteamPath()
	if err != nil {
		return err
	}

	userdataPath := filepath.Join(steamPath, "userdata")

	for userID, launchOpts := range b.Accounts {
		configPath := filepath.Join(userdataPath, userID, "config", "localconfig.vdf")
		
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}

		content, err := os.ReadFile(configPath)
		if err != nil {
			continue
		}

		contentStr := string(content)

		// Remove ALL existing LaunchOptions entries
		launchOptRe := regexp.MustCompile(`(?m)^\s*"LaunchOptions"\s+"[^"]*"\s*$`)
		matches := launchOptRe.FindAllStringIndex(contentStr, -1)
		
		for i := len(matches) - 1; i >= 0; i-- {
			match := matches[i]
			lineStart := match[0]
			lineEnd := match[1]
			
			if lineEnd < len(contentStr) && contentStr[lineEnd] == '\n' {
				lineEnd++
			} else if lineEnd < len(contentStr) && contentStr[lineEnd] == '\r' {
				lineEnd++
				if lineEnd < len(contentStr) && contentStr[lineEnd] == '\n' {
					lineEnd++
				}
			}
			
			contentStr = contentStr[:lineStart] + contentStr[lineEnd:]
		}

		// If there were launch options before, restore them
		if launchOpts != "" {
			cs2StartRe := regexp.MustCompile(`("730"\s*\n\s*\{)\s*\n`)
			if cs2StartRe.MatchString(contentStr) {
				contentStr = cs2StartRe.ReplaceAllString(contentStr, "${1}\n\t\t\t\t\"LaunchOptions\"\t\t\""+launchOpts+"\"\n")
			}
		}

		os.WriteFile(configPath, []byte(contentStr), 0644)
	}

	return nil
}
