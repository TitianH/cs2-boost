package hpet

import (
	"fmt"
	"os/exec"
	"strings"

	"cs2-boost/pkg/backup"
)

func Disable() error {
	cmd := exec.Command("bcdedit", "/deletevalue", "useplatformclock")
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "not currently exist") || strings.Contains(outputStr, "nicht vorhanden") {
			return nil
		}
		return fmt.Errorf("bcdedit failed: %w - %s", err, outputStr)
	}
	return nil
}

func Enable() error {
	cmd := exec.Command("bcdedit", "/set", "useplatformclock", "true")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("bcdedit failed: %w - %s", err, string(output))
	}
	return nil
}

func GetStatus() (bool, error) {
	cmd := exec.Command("bcdedit", "/enum", "{current}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("bcdedit failed: %w", err)
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "useplatformclock") && strings.Contains(outputStr, "Yes") {
		return true, nil
	}

	return false, nil
}

func BackupCurrent() *backup.HPETBackup {
	status, err := GetStatus()
	if err != nil {
		return &backup.HPETBackup{WasEnabled: false}
	}
	return &backup.HPETBackup{WasEnabled: status}
}

func RestoreFromBackup(b *backup.HPETBackup) error {
	if b == nil {
		return nil
	}

	if b.WasEnabled {
		return Enable()
	}
	return Disable()
}
