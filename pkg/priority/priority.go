package priority

import (
	"fmt"

	"cs2-boost/pkg/backup"
	"golang.org/x/sys/windows/registry"
)

const (
	registryPath = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Image File Execution Options\cs2.exe\PerfOptions`
	priorityHigh = 5
)

func SetCS2Priority() error {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, registryPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create registry key: %w", err)
	}
	defer key.Close()

	err = key.SetDWordValue("CpuPriorityClass", priorityHigh)
	if err != nil {
		return fmt.Errorf("failed to set priority value: %w", err)
	}

	return nil
}

func RemoveCS2Priority() error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Image File Execution Options\cs2.exe`, registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	err = registry.DeleteKey(key, "PerfOptions")
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to delete registry key: %w", err)
	}

	return nil
}

func GetStatus() (bool, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, registryPath, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetIntegerValue("CpuPriorityClass")
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to read priority value: %w", err)
	}

	return value == priorityHigh, nil
}

func BackupCurrent() *backup.PriorityBackup {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, registryPath, registry.QUERY_VALUE)
	if err != nil {
		return &backup.PriorityBackup{Existed: false}
	}
	defer key.Close()

	value, _, err := key.GetIntegerValue("CpuPriorityClass")
	if err != nil {
		return &backup.PriorityBackup{Existed: false}
	}

	return &backup.PriorityBackup{
		Existed:       true,
		PreviousValue: value,
	}
}

func RestoreFromBackup(b *backup.PriorityBackup) error {
	if b == nil {
		return nil
	}

	if !b.Existed {
		return RemoveCS2Priority()
	}

	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, registryPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create registry key: %w", err)
	}
	defer key.Close()

	err = key.SetDWordValue("CpuPriorityClass", uint32(b.PreviousValue))
	if err != nil {
		return fmt.Errorf("failed to restore priority value: %w", err)
	}

	return nil
}
