# CS2 Performance Boost Tool

Modular Go tool for optimizing Counter-Strike 2 performance on Windows.

## Quick Start

**Build:**
```powershell
.\build.ps1
```

**Install optimizations (run as Administrator):**
```powershell
.\cs2-boost.exe install
```

Then restart your PC.

## Features

- **Process Priority**: Sets cs2.exe to High priority via Windows registry
- **HPET Disable**: Disables High Precision Event Timer to reduce input lag
- **Game DVR/Game Bar**: Disables Xbox overlay features that impact performance
- **GPU Scheduling**: Enables Hardware-accelerated GPU scheduling (Windows 10 2004+)
- **Network Optimization**: Disables Nagle's Algorithm for lower latency
- **Power Plan**: Sets Windows to High Performance mode
- **Visual Effects**: Optimizes Windows visual effects for performance
- **Safe Backup/Restore**: Automatically backs up all settings before making changes
- **Complete Restoration**: Uninstall restores exact previous settings from backup
- **Modular Design**: Clean package structure for easy extension
- **Status Checking**: View current optimization status

## Usage

```powershell
cs2-boost install     # Apply all optimizations (creates backup first)
cs2-boost uninstall   # Restore previous settings from backup
cs2-boost status      # Check current status
cs2-boost version     # Show version
```

## Safety Features

### Automatic Backup
Before applying any optimizations, the tool automatically:
1. Backs up all current settings to `cs2-boost-backup.json`
2. Stores exact registry values, power plan GUIDs, and system configurations
3. Prevents installation if backup already exists (with confirmation prompt)

### Safe Restoration
When uninstalling:
1. Loads the backup file
2. Restores **exact previous values** (not defaults)
3. Deletes backup file after successful restoration
4. Falls back to default uninstall if no backup exists (with confirmation)

The backup file is saved in the same directory as the executable.

## Building from Source

Requires Go 1.21 or later:

```powershell
go mod tidy
go build -o cs2-boost.exe .
```

Or use the build script:
```powershell
.\build.ps1
```

## Project Structure

```
cs2-boost/
├── main.go                 # CLI entry point
├── pkg/
│   ├── backup/
│   │   └── backup.go      # Backup/restore system
│   ├── gamedvr/
│   │   └── gamedvr.go     # Xbox Game DVR/Game Bar management
│   ├── gpu/
│   │   └── gpu.go         # GPU scheduling optimization
│   ├── hpet/
│   │   └── hpet.go        # HPET management
│   ├── network/
│   │   └── network.go     # Network latency optimization
│   ├── power/
│   │   └── power.go       # Power plan management
│   ├── priority/
│   │   └── priority.go    # Process priority management
│   ├── utils/
│   │   └── admin.go       # Admin privilege checks
│   └── visual/
│       └── visual.go      # Visual effects optimization
├── go.mod
└── build.ps1
```

Each module includes `BackupCurrent()` and `RestoreFromBackup()` functions for safe state management.

## PowerShell Alternative

For those who prefer PowerShell, `Set-CS2Priority.ps1` is still available for priority-only setting.
