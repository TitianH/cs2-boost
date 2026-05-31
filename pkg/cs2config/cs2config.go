package cs2config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cs2-boost/pkg/backup"
	"cs2-boost/pkg/steam"
)

const (
	cs2BoostMarkerStart = "// === CS2-BOOST START - DO NOT EDIT THIS SECTION ==="
	cs2BoostMarkerEnd   = "// === CS2-BOOST END ==="
)

var optimizedSettings = map[string]string{
	"rate":                                       "786432",
	"cl_interp":                                  "0",
	"cl_interp_ratio":                            "1",
	"cl_updaterate":                              "128",
	"cl_cmdrate":                                 "128",
	"fps_max":                                    "0",
	"engine_low_latency_sleep_after_client_tick": "true",
	"m_rawinput":                                 "1",
	"m_mousespeed":                               "0",
}

func GetConfigFilePath() (string, error) {
	// Try userdata first (user-specific config)
	userConfigPath, err := steam.GetUserDataPath()
	if err == nil {
		autoexecPath := filepath.Join(userConfigPath, "autoexec.cfg")
		// Create directory if it doesn't exist
		os.MkdirAll(userConfigPath, 0755)
		return autoexecPath, nil
	}

	// Fallback to game directory
	configPath, err := steam.GetCS2ConfigPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(configPath, "autoexec.cfg"), nil
}

func ApplyOptimizedConfig() error {
	steamPath, err := steam.GetSteamPath()
	if err != nil {
		return fmt.Errorf("failed to locate Steam: %w", err)
	}

	userdataPath := filepath.Join(steamPath, "userdata")
	entries, err := os.ReadDir(userdataPath)
	if err != nil {
		return fmt.Errorf("failed to read userdata: %w", err)
	}

	successCount := 0
	var lastError error

	// Build cs2-boost section once
	var cs2BoostSection strings.Builder
	cs2BoostSection.WriteString("\n" + cs2BoostMarkerStart + "\n")
	cs2BoostSection.WriteString("// Network settings\n")
	cs2BoostSection.WriteString(fmt.Sprintf("rate %s\n", optimizedSettings["rate"]))
	cs2BoostSection.WriteString(fmt.Sprintf("cl_interp %s\n", optimizedSettings["cl_interp"]))
	cs2BoostSection.WriteString(fmt.Sprintf("cl_interp_ratio %s\n", optimizedSettings["cl_interp_ratio"]))
	cs2BoostSection.WriteString(fmt.Sprintf("cl_updaterate %s\n", optimizedSettings["cl_updaterate"]))
	cs2BoostSection.WriteString(fmt.Sprintf("cl_cmdrate %s\n", optimizedSettings["cl_cmdrate"]))
	cs2BoostSection.WriteString("\n// Performance\n")
	cs2BoostSection.WriteString(fmt.Sprintf("fps_max %s\n", optimizedSettings["fps_max"]))
	cs2BoostSection.WriteString(fmt.Sprintf("engine_low_latency_sleep_after_client_tick %s\n", optimizedSettings["engine_low_latency_sleep_after_client_tick"]))
	cs2BoostSection.WriteString("\n// Mouse\n")
	cs2BoostSection.WriteString(fmt.Sprintf("m_rawinput %s\n", optimizedSettings["m_rawinput"]))
	cs2BoostSection.WriteString(fmt.Sprintf("m_mousespeed %s\n", optimizedSettings["m_mousespeed"]))
	cs2BoostSection.WriteString(cs2BoostMarkerEnd + "\n")
	cs2BoostContent := cs2BoostSection.String()

	// Apply to all Steam accounts
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "0" || entry.Name() == "ac" {
			continue
		}

		configPath := filepath.Join(userdataPath, entry.Name(), "730", "local", "cfg", "autoexec.cfg")
		configDir := filepath.Dir(configPath)

		// Check if CS2 config directory exists
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			continue
		}

		fmt.Printf("\n      User %s: %s", entry.Name(), configPath)

		// Create directory if needed
		os.MkdirAll(configDir, 0755)

		// Read existing config
		var existingContent string
		data, err := os.ReadFile(configPath)
		if err == nil {
			existingContent = string(data)
		}

		// Remove old cs2-boost section if it exists
		existingContent = removeCS2BoostSection(existingContent)

		// Append cs2-boost section to existing content
		newContent := strings.TrimSpace(existingContent) + "\n" + cs2BoostContent

		// Write merged config
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
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

func removeCS2BoostSection(content string) string {
	startIdx := strings.Index(content, cs2BoostMarkerStart)
	if startIdx == -1 {
		return content
	}

	endIdx := strings.Index(content, cs2BoostMarkerEnd)
	if endIdx == -1 {
		return content
	}

	// Remove everything from start marker to end marker (including end marker line)
	endOfLine := strings.Index(content[endIdx:], "\n")
	if endOfLine != -1 {
		endIdx += endOfLine + 1
	} else {
		endIdx += len(cs2BoostMarkerEnd)
	}

	return content[:startIdx] + content[endIdx:]
}

func RemoveOptimizations() error {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to locate CS2 config: %w", err)
	}

	// Read current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Remove cs2-boost section
	newContent := removeCS2BoostSection(string(content))
	newContent = strings.TrimSpace(newContent)

	// If file is empty after removal, delete it
	if newContent == "" {
		os.Remove(configPath)
		return nil
	}

	// Write back without cs2-boost section
	if err := os.WriteFile(configPath, []byte(newContent+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func GetStatus() (bool, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return false, nil
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return false, nil
	}

	// Check if our config section is present
	return strings.Contains(string(content), cs2BoostMarkerStart), nil
}

func BackupCurrent() *backup.CS2ConfigBackup {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return &backup.CS2ConfigBackup{
			Existed: false,
		}
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return &backup.CS2ConfigBackup{
			Existed: false,
		}
	}

	contentStr := string(content)
	return &backup.CS2ConfigBackup{
		Existed:         true,
		PreviousContent: &contentStr,
		ConfigPath:      configPath,
	}
}

func RestoreFromBackup(b *backup.CS2ConfigBackup) error {
	if b == nil {
		return nil
	}

	configPath := b.ConfigPath
	if configPath == "" {
		var err error
		configPath, err = GetConfigFilePath()
		if err != nil {
			return err
		}
	}

	if !b.Existed {
		// Config didn't exist before, remove our section
		return RemoveOptimizations()
	}

	if b.PreviousContent == nil {
		return RemoveOptimizations()
	}

	// Restore exact previous content
	return os.WriteFile(configPath, []byte(*b.PreviousContent), 0644)
}
