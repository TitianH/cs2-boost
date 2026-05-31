package power

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"cs2-boost/pkg/backup"
)

const (
	highPerformanceGUID = "8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c"
)

func SetHighPerformance() error {
	cmd := exec.Command("powercfg", "/setactive", highPerformanceGUID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set high performance mode: %w - %s", err, string(output))
	}
	return nil
}

func DisableUSBPowerSaving() error {
	cmd := exec.Command("powercfg", "/setacvalueindex", highPerformanceGUID, "2a737441-1930-4402-8d77-b2bebba308a3", "48e6b7a6-50f5-4782-a5d4-53bb8f07e226", "0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable USB power saving: %w - %s", err, string(output))
	}

	cmd = exec.Command("powercfg", "/setdcvalueindex", highPerformanceGUID, "2a737441-1930-4402-8d77-b2bebba308a3", "48e6b7a6-50f5-4782-a5d4-53bb8f07e226", "0")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable USB power saving (DC): %w - %s", err, string(output))
	}

	cmd = exec.Command("powercfg", "/setactive", highPerformanceGUID)
	cmd.Run()

	return nil
}

func GetActivePlan() (string, error) {
	cmd := exec.Command("powercfg", "/getactivescheme")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get active power plan: %w", err)
	}

	outputStr := string(output)
	if strings.Contains(outputStr, highPerformanceGUID) || strings.Contains(strings.ToLower(outputStr), "high performance") || strings.Contains(strings.ToLower(outputStr), "höchstleistung") {
		return "High Performance", nil
	}

	lines := strings.Split(outputStr, "\n")
	if len(lines) > 0 {
		parts := strings.Split(lines[0], ":")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "Unknown", nil
}

func GetStatus() (bool, error) {
	plan, err := GetActivePlan()
	if err != nil {
		return false, err
	}
	return plan == "High Performance", nil
}

func GetActivePlanGUID() (string, error) {
	cmd := exec.Command("powercfg", "/getactivescheme")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get active power plan: %w", err)
	}

	re := regexp.MustCompile(`([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) > 0 {
		return matches[1], nil
	}

	return "", fmt.Errorf("could not extract GUID from output")
}

func BackupCurrent() *backup.PowerBackup {
	guid, err := GetActivePlanGUID()
	if err != nil {
		return &backup.PowerBackup{}
	}

	name, _ := GetActivePlan()

	return &backup.PowerBackup{
		PreviousPlanGUID: guid,
		PreviousPlanName: name,
	}
}

func RestoreFromBackup(b *backup.PowerBackup) error {
	if b == nil || b.PreviousPlanGUID == "" {
		return nil
	}

	cmd := exec.Command("powercfg", "/setactive", b.PreviousPlanGUID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restore power plan: %w - %s", err, string(output))
	}

	return nil
}
