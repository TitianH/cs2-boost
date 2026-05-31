# CS2 Performance Boost Tool

Modular Go tool for optimizing Counter-Strike 2 performance on Windows.

## Quick Start

**Build:**
```powershell
.\build.bat
```

**Install optimizations (run as Administrator):**
```powershell
.\cs2-boost.exe install
```

Then restart your PC.

**Important Notes:**
- **Close Steam before running install** - Launch options can only be applied when Steam is not running
- Steam must be installed and CS2 must be in your library for launch options and config tweaks to work

## Features

- **Process Priority**: Sets cs2.exe to High priority via Windows registry
- **HPET Disable**: Disables High Precision Event Timer to reduce input lag
- **Game DVR/Game Bar**: Disables Xbox overlay features that impact performance
- **GPU Scheduling**: Enables Hardware-accelerated GPU scheduling (Windows 10 2004+)
- **Network Optimization**: Disables Nagle's Algorithm for lower latency
- **Power Plan**: Sets Windows to High Performance mode
- **Visual Effects**: Optimizes Windows visual effects for performance
- **CS2 Launch Options**: Auto-detects Steam and applies optimized launch options
- **CS2 Config Tweaks**: Creates optimized autoexec.cfg with network/performance settings
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
```batch
.\build.bat
```

## Project Structure

```
cs2-boost/
├── main.go                 # CLI entry point
├── pkg/
│   ├── backup/
│   │   └── backup.go      # Backup/restore system
│   ├── cs2config/
│   │   └── cs2config.go   # CS2 autoexec.cfg management
│   ├── gamedvr/
│   │   └── gamedvr.go     # Xbox Game DVR/Game Bar management
│   ├── gpu/
│   │   └── gpu.go         # GPU scheduling optimization
│   ├── hpet/
│   │   └── hpet.go        # HPET management
│   ├── launchopts/
│   │   └── launchopts.go  # Steam launch options management
│   ├── network/
│   │   └── network.go     # Network latency optimization
│   ├── power/
│   │   └── power.go       # Power plan management
│   ├── priority/
│   │   └── priority.go    # Process priority management
│   ├── steam/
│   │   └── steam.go       # Steam installation detection
│   ├── utils/
│   │   └── admin.go       # Admin privilege checks
│   └── visual/
│       └── visual.go      # Visual effects optimization
├── go.mod
└── .gitignore
```

Each module includes `BackupCurrent()` and `RestoreFromBackup()` functions for safe state management.

### CS2-Specific Features
- **Auto-detection**: Finds Steam installation and CS2 game files automatically
- **Launch Options**: Applies optimized flags (`-high`, `-nojoy`, `-novid`, auto-detected `-threads`)
- **Config Tweaks**: Creates `autoexec.cfg` with network settings (rate 786432, cl_interp 0, etc.)

## What Gets Applied

### Launch Options
```
-high -nojoy -novid +fps_max 0 +cl_forcepreload 1 -threads [auto-detected] +exec autoexec
```

### autoexec.cfg (Merged, not replaced)
The tool adds a marked section to your existing autoexec.cfg:
```
// === CS2-BOOST START - DO NOT EDIT THIS SECTION ===
// Network settings
rate 786432
cl_interp 0
cl_interp_ratio 1
cl_updaterate 128
cl_cmdrate 128

// Performance
fps_max 0
engine_low_latency_sleep_after_client_tick true

// Mouse
m_rawinput 1
m_mousespeed 0
// === CS2-BOOST END ===
```

**Your existing settings are preserved!** The tool only manages the section between the markers.
