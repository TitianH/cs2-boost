package network

import (
	"fmt"

	"cs2-boost/pkg/backup"

	"golang.org/x/sys/windows/registry"
)

func OptimizeForGaming() error {
	interfacesPath := `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces`

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, interfacesPath, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return fmt.Errorf("failed to open interfaces key: %w", err)
	}
	defer key.Close()

	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return fmt.Errorf("failed to read subkeys: %w", err)
	}

	for _, subkey := range subkeys {
		interfaceKey, err := registry.OpenKey(registry.LOCAL_MACHINE, interfacesPath+`\`+subkey, registry.SET_VALUE)
		if err != nil {
			continue
		}

		interfaceKey.SetDWordValue("TcpAckFrequency", 1)
		interfaceKey.SetDWordValue("TCPNoDelay", 1)
		interfaceKey.Close()
	}

	if err := setGlobalTCPSettings(); err != nil {
		return err
	}

	return nil
}

func ResetOptimizations() error {
	interfacesPath := `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces`

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, interfacesPath, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return fmt.Errorf("failed to open interfaces key: %w", err)
	}
	defer key.Close()

	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return fmt.Errorf("failed to read subkeys: %w", err)
	}

	for _, subkey := range subkeys {
		interfaceKey, err := registry.OpenKey(registry.LOCAL_MACHINE, interfacesPath+`\`+subkey, registry.SET_VALUE)
		if err != nil {
			continue
		}

		interfaceKey.DeleteValue("TcpAckFrequency")
		interfaceKey.DeleteValue("TCPNoDelay")
		interfaceKey.Close()
	}

	if err := resetGlobalTCPSettings(); err != nil {
		return err
	}

	return nil
}

func setGlobalTCPSettings() error {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create TCP parameters key: %w", err)
	}
	defer key.Close()

	key.SetDWordValue("TcpAckFrequency", 1)
	key.SetDWordValue("TCPNoDelay", 1)

	return nil
}

func resetGlobalTCPSettings() error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`, registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return fmt.Errorf("failed to open TCP parameters key: %w", err)
	}
	defer key.Close()

	key.DeleteValue("TcpAckFrequency")
	key.DeleteValue("TCPNoDelay")

	return nil
}

func GetStatus() (bool, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetIntegerValue("TcpAckFrequency")
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to read TcpAckFrequency value: %w", err)
	}

	return value == 1, nil
}

func BackupCurrent() *backup.NetworkBackup {
	backupData := &backup.NetworkBackup{
		Interfaces: make(map[string]backup.NetInterface),
	}

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`, registry.QUERY_VALUE)
	if err == nil {
		defer key.Close()
		if value, _, err := key.GetIntegerValue("TcpAckFrequency"); err == nil {
			backupData.GlobalTcpAckFrequency = &value
		}
		if value, _, err := key.GetIntegerValue("TCPNoDelay"); err == nil {
			backupData.GlobalTCPNoDelay = &value
		}
	}

	interfacesPath := `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces`
	interfacesKey, err := registry.OpenKey(registry.LOCAL_MACHINE, interfacesPath, registry.ENUMERATE_SUB_KEYS)
	if err == nil {
		defer interfacesKey.Close()
		subkeys, _ := interfacesKey.ReadSubKeyNames(-1)
		for _, subkey := range subkeys {
			interfaceKey, err := registry.OpenKey(registry.LOCAL_MACHINE, interfacesPath+`\`+subkey, registry.QUERY_VALUE)
			if err != nil {
				continue
			}
			netInterface := backup.NetInterface{}
			if value, _, err := interfaceKey.GetIntegerValue("TcpAckFrequency"); err == nil {
				netInterface.TcpAckFrequency = &value
			}
			if value, _, err := interfaceKey.GetIntegerValue("TCPNoDelay"); err == nil {
				netInterface.TCPNoDelay = &value
			}
			if netInterface.TcpAckFrequency != nil || netInterface.TCPNoDelay != nil {
				backupData.Interfaces[subkey] = netInterface
			}
			interfaceKey.Close()
		}
	}

	return backupData
}

func RestoreFromBackup(b *backup.NetworkBackup) error {
	if b == nil {
		return nil
	}

	if b.GlobalTcpAckFrequency != nil || b.GlobalTCPNoDelay != nil {
		key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`, registry.SET_VALUE)
		if err == nil {
			defer key.Close()
			if b.GlobalTcpAckFrequency != nil {
				key.SetDWordValue("TcpAckFrequency", uint32(*b.GlobalTcpAckFrequency))
			} else {
				key.DeleteValue("TcpAckFrequency")
			}
			if b.GlobalTCPNoDelay != nil {
				key.SetDWordValue("TCPNoDelay", uint32(*b.GlobalTCPNoDelay))
			} else {
				key.DeleteValue("TCPNoDelay")
			}
		}
	}

	interfacesPath := `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces`
	for subkey, netInterface := range b.Interfaces {
		interfaceKey, _, err := registry.CreateKey(registry.LOCAL_MACHINE, interfacesPath+`\`+subkey, registry.SET_VALUE)
		if err != nil {
			continue
		}
		if netInterface.TcpAckFrequency != nil {
			interfaceKey.SetDWordValue("TcpAckFrequency", uint32(*netInterface.TcpAckFrequency))
		} else {
			interfaceKey.DeleteValue("TcpAckFrequency")
		}
		if netInterface.TCPNoDelay != nil {
			interfaceKey.SetDWordValue("TCPNoDelay", uint32(*netInterface.TCPNoDelay))
		} else {
			interfaceKey.DeleteValue("TCPNoDelay")
		}
		interfaceKey.Close()
	}

	return nil
}
