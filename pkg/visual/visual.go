package visual

import (
	"fmt"

	"cs2-boost/pkg/backup"

	"golang.org/x/sys/windows/registry"
)

func OptimizeForPerformance() error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create VisualEffects key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("VisualFXSetting", 2); err != nil {
		return fmt.Errorf("failed to set visual effects: %w", err)
	}

	return nil
}

func ResetToDefault() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects`, registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return fmt.Errorf("failed to open VisualEffects key: %w", err)
	}
	defer key.Close()

	err = key.DeleteValue("VisualFXSetting")
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to reset visual effects: %w", err)
	}

	return nil
}

func GetStatus() (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects`, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetIntegerValue("VisualFXSetting")
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to read VisualFXSetting value: %w", err)
	}

	return value == 2, nil
}

func BackupCurrent() *backup.VisualBackup {
	backupData := &backup.VisualBackup{}

	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects`, registry.QUERY_VALUE)
	if err == nil {
		defer key.Close()
		if value, _, err := key.GetIntegerValue("VisualFXSetting"); err == nil {
			backupData.VisualFXSetting = &value
		}
	}

	return backupData
}

func RestoreFromBackup(b *backup.VisualBackup) error {
	if b == nil || b.VisualFXSetting == nil {
		return ResetToDefault()
	}

	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create VisualEffects key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("VisualFXSetting", uint32(*b.VisualFXSetting)); err != nil {
		return fmt.Errorf("failed to restore VisualFXSetting: %w", err)
	}

	return nil
}
