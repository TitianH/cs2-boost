package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type BackupData struct {
	Priority       *PriorityBackup       `json:"priority,omitempty"`
	HPET           *HPETBackup           `json:"hpet,omitempty"`
	GameDVR        *GameDVRBackup        `json:"gamedvr,omitempty"`
	GPU            *GPUBackup            `json:"gpu,omitempty"`
	Network        *NetworkBackup        `json:"network,omitempty"`
	Power          *PowerBackup          `json:"power,omitempty"`
	Visual         *VisualBackup         `json:"visual,omitempty"`
	CS2Config      *CS2ConfigBackup      `json:"cs2config,omitempty"`
	LaunchOpts     *LaunchOptsBackup     `json:"launchopts,omitempty"`
	BackupVersion  string                `json:"backup_version"`
}

type PriorityBackup struct {
	Existed       bool   `json:"existed"`
	PreviousValue uint64 `json:"previous_value,omitempty"`
}

type HPETBackup struct {
	WasEnabled bool `json:"was_enabled"`
}

type GameDVRBackup struct {
	GameDVREnabled    *uint64 `json:"gamedvr_enabled,omitempty"`
	AppCaptureEnabled *uint64 `json:"appcapture_enabled,omitempty"`
}

type GPUBackup struct {
	HwSchMode *uint64 `json:"hwschmode,omitempty"`
}

type NetworkBackup struct {
	GlobalTcpAckFrequency *uint64            `json:"global_tcp_ack_frequency,omitempty"`
	GlobalTCPNoDelay      *uint64            `json:"global_tcp_nodelay,omitempty"`
	Interfaces            map[string]NetInterface `json:"interfaces,omitempty"`
}

type NetInterface struct {
	TcpAckFrequency *uint64 `json:"tcp_ack_frequency,omitempty"`
	TCPNoDelay      *uint64 `json:"tcp_nodelay,omitempty"`
}

type PowerBackup struct {
	PreviousPlanGUID string `json:"previous_plan_guid"`
	PreviousPlanName string `json:"previous_plan_name"`
}

type VisualBackup struct {
	VisualFXSetting *uint64 `json:"visualfx_setting,omitempty"`
}

type CS2ConfigBackup struct {
	Configs map[string]CS2ConfigAccount `json:"configs,omitempty"` // key: Steam user ID
}

type CS2ConfigAccount struct {
	Existed         bool    `json:"existed"`
	PreviousContent *string `json:"previous_content,omitempty"`
	ConfigPath      string  `json:"config_path,omitempty"`
}

type LaunchOptsBackup struct {
	Accounts map[string]string `json:"accounts,omitempty"` // key: Steam user ID, value: previous launch options
}

func GetBackupPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	
	dir := filepath.Dir(exePath)
	return filepath.Join(dir, "cs2-boost-backup.json"), nil
}

func Save(data *BackupData) error {
	data.BackupVersion = "1.0"
	
	backupPath, err := GetBackupPath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup data: %w", err)
	}

	err = os.WriteFile(backupPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

func Load() (*BackupData, error) {
	backupPath, err := GetBackupPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no backup file found")
	}

	jsonData, err := os.ReadFile(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %w", err)
	}

	var data BackupData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup data: %w", err)
	}

	return &data, nil
}

func Exists() bool {
	backupPath, err := GetBackupPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(backupPath)
	return err == nil
}

func Delete() error {
	backupPath, err := GetBackupPath()
	if err != nil {
		return err
	}

	if !Exists() {
		return nil
	}

	err = os.Remove(backupPath)
	if err != nil {
		return fmt.Errorf("failed to delete backup file: %w", err)
	}

	return nil
}
