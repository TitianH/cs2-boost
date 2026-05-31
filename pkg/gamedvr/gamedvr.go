package gamedvr

import (
	"fmt"

	"cs2-boost/pkg/backup"
	"golang.org/x/sys/windows/registry"
)

func Disable() error {
	if err := disableGameDVR(); err != nil {
		return err
	}
	if err := disableGameBar(); err != nil {
		return err
	}
	return nil
}

func Enable() error {
	if err := enableGameDVR(); err != nil {
		return err
	}
	if err := enableGameBar(); err != nil {
		return err
	}
	return nil
}

func disableGameDVR() error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `System\GameConfigStore`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create GameConfigStore key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("GameDVR_Enabled", 0); err != nil {
		return fmt.Errorf("failed to disable GameDVR: %w", err)
	}

	return nil
}

func disableGameBar() error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\GameDVR`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create GameDVR key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("AppCaptureEnabled", 0); err != nil {
		return fmt.Errorf("failed to disable Game Bar: %w", err)
	}

	return nil
}

func enableGameDVR() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `System\GameConfigStore`, registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return fmt.Errorf("failed to open GameConfigStore key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("GameDVR_Enabled", 1); err != nil {
		return fmt.Errorf("failed to enable GameDVR: %w", err)
	}

	return nil
}

func enableGameBar() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\GameDVR`, registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return fmt.Errorf("failed to open GameDVR key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("AppCaptureEnabled", 1); err != nil {
		return fmt.Errorf("failed to enable Game Bar: %w", err)
	}

	return nil
}

func GetStatus() (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `System\GameConfigStore`, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return true, nil
		}
		return false, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetIntegerValue("GameDVR_Enabled")
	if err != nil {
		if err == registry.ErrNotExist {
			return true, nil
		}
		return false, fmt.Errorf("failed to read GameDVR value: %w", err)
	}

	return value == 1, nil
}

func BackupCurrent() *backup.GameDVRBackup {
	backupData := &backup.GameDVRBackup{}

	key, err := registry.OpenKey(registry.CURRENT_USER, `System\GameConfigStore`, registry.QUERY_VALUE)
	if err == nil {
		defer key.Close()
		if value, _, err := key.GetIntegerValue("GameDVR_Enabled"); err == nil {
			backupData.GameDVREnabled = &value
		}
	}

	key2, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\GameDVR`, registry.QUERY_VALUE)
	if err == nil {
		defer key2.Close()
		if value, _, err := key2.GetIntegerValue("AppCaptureEnabled"); err == nil {
			backupData.AppCaptureEnabled = &value
		}
	}

	return backupData
}

func RestoreFromBackup(b *backup.GameDVRBackup) error {
	if b == nil {
		return nil
	}

	if b.GameDVREnabled != nil {
		key, _, err := registry.CreateKey(registry.CURRENT_USER, `System\GameConfigStore`, registry.SET_VALUE)
		if err == nil {
			defer key.Close()
			key.SetDWordValue("GameDVR_Enabled", uint32(*b.GameDVREnabled))
		}
	}

	if b.AppCaptureEnabled != nil {
		key, _, err := registry.CreateKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\GameDVR`, registry.SET_VALUE)
		if err == nil {
			defer key.Close()
			key.SetDWordValue("AppCaptureEnabled", uint32(*b.AppCaptureEnabled))
		}
	}

	return nil
}
