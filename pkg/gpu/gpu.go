package gpu

import (
	"fmt"

	"cs2-boost/pkg/backup"

	"golang.org/x/sys/windows/registry"
)

func EnableHardwareScheduling() error {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\GraphicsDrivers`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create GraphicsDrivers key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("HwSchMode", 2); err != nil {
		return fmt.Errorf("failed to enable hardware scheduling: %w", err)
	}

	return nil
}

func DisableHardwareScheduling() error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\GraphicsDrivers`, registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return fmt.Errorf("failed to open GraphicsDrivers key: %w", err)
	}
	defer key.Close()

	err = key.DeleteValue("HwSchMode")
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to disable hardware scheduling: %w", err)
	}

	return nil
}

func GetStatus() (bool, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\GraphicsDrivers`, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetIntegerValue("HwSchMode")
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to read HwSchMode value: %w", err)
	}

	return value == 2, nil
}

func BackupCurrent() *backup.GPUBackup {
	backupData := &backup.GPUBackup{}

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\GraphicsDrivers`, registry.QUERY_VALUE)
	if err == nil {
		defer key.Close()
		if value, _, err := key.GetIntegerValue("HwSchMode"); err == nil {
			backupData.HwSchMode = &value
		}
	}

	return backupData
}

func RestoreFromBackup(b *backup.GPUBackup) error {
	if b == nil || b.HwSchMode == nil {
		return DisableHardwareScheduling()
	}

	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\GraphicsDrivers`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create GraphicsDrivers key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("HwSchMode", uint32(*b.HwSchMode)); err != nil {
		return fmt.Errorf("failed to restore HwSchMode: %w", err)
	}

	return nil
}
