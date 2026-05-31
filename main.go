package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"cs2-boost/pkg/backup"
	"cs2-boost/pkg/cs2config"
	"cs2-boost/pkg/gamedvr"
	"cs2-boost/pkg/gpu"
	"cs2-boost/pkg/hpet"
	"cs2-boost/pkg/launchopts"
	"cs2-boost/pkg/network"
	"cs2-boost/pkg/power"
	"cs2-boost/pkg/priority"
	"cs2-boost/pkg/utils"
	"cs2-boost/pkg/visual"
)

const (
	version = "1.0.0"
)

func main() {
	if !utils.IsAdmin() {
		fmt.Println("Error: This tool requires administrator privileges")
		fmt.Println("Please run as administrator")
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	installCmd := flag.NewFlagSet("install", flag.ExitOnError)
	uninstallCmd := flag.NewFlagSet("uninstall", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)

	if len(os.Args) < 2 {
		runInteractiveMode()
		return
	}

	switch os.Args[1] {
	case "install":
		installCmd.Parse(os.Args[2:])
		handleInstall()
	case "uninstall":
		uninstallCmd.Parse(os.Args[2:])
		handleUninstall()
	case "status":
		statusCmd.Parse(os.Args[2:])
		handleStatus()
	case "version":
		fmt.Printf("cs2-boost v%s\n", version)
	default:
		printUsage()
		os.Exit(1)
	}
}

func runInteractiveMode() {
	fmt.Println("===========================================")
	fmt.Println("  CS2 Performance Boost Tool v" + version)
	fmt.Println("===========================================")
	fmt.Println()
	fmt.Println("What would you like to do?")
	fmt.Println()
	fmt.Println("  1. Install optimizations")
	fmt.Println("  2. Uninstall (restore previous settings)")
	fmt.Println("  3. Check status")
	fmt.Println("  4. Exit")
	fmt.Println()
	fmt.Print("Enter your choice (1-4): ")

	var choice string
	fmt.Scanln(&choice)
	fmt.Println()

	switch choice {
	case "1":
		handleInstall()
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
	case "2":
		handleUninstall()
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
	case "3":
		handleStatus()
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
	case "4":
		fmt.Println("Exiting...")
		return
	default:
		fmt.Println("Invalid choice. Exiting...")
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
	}
}

func printUsage() {
	fmt.Println("CS2 Performance Boost Tool")
	fmt.Println("===========================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  cs2-boost install     - Apply all optimizations")
	fmt.Println("  cs2-boost uninstall   - Remove all optimizations")
	fmt.Println("  cs2-boost status      - Check current status")
	fmt.Println("  cs2-boost version     - Show version")
	fmt.Println()
	fmt.Println("Optimizations:")
	fmt.Println("  • Set cs2.exe process priority to High")
	fmt.Println("  • Disable HPET (High Precision Event Timer)")
	fmt.Println("  • Disable Xbox Game DVR and Game Bar")
	fmt.Println("  • Enable Hardware-accelerated GPU scheduling")
	fmt.Println("  • Optimize network settings (disable Nagle's Algorithm)")
	fmt.Println("  • Set High Performance power plan")
	fmt.Println("  • Optimize visual effects for performance")
	fmt.Println("  • Apply CS2 launch options (Steam)")
	fmt.Println("  • Apply CS2 config tweaks (autoexec.cfg)")
	fmt.Println()
	fmt.Println("Note: Requires administrator privileges")
}

func handleInstall() {
	if backup.Exists() {
		fmt.Println("Warning: A backup file already exists.")
		fmt.Println("This means optimizations may already be installed.")
		fmt.Print("Continue anyway? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Installation cancelled.")
			return
		}
	}

	fmt.Println("Installing CS2 optimizations...")
	fmt.Println()

	fmt.Print("[0/9] Creating backup of current settings... ")
	backupData := &backup.BackupData{
		Priority:   priority.BackupCurrent(),
		HPET:       hpet.BackupCurrent(),
		GameDVR:    gamedvr.BackupCurrent(),
		GPU:        gpu.BackupCurrent(),
		Network:    network.BackupCurrent(),
		Power:      power.BackupCurrent(),
		Visual:     visual.BackupCurrent(),
		CS2Config:  cs2config.BackupCurrent(),
		LaunchOpts: launchopts.BackupCurrent(),
	}
	if err := backup.Save(backupData); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
		fmt.Println("\nAborting installation for safety.")
		return
	}
	fmt.Println("OK")

	fmt.Print("[1/9] Setting cs2.exe priority to High... ")
	if err := priority.SetCS2Priority(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[2/9] Disabling HPET... ")
	if err := hpet.Disable(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[3/9] Disabling Game DVR/Game Bar... ")
	if err := gamedvr.Disable(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[4/9] Enabling Hardware GPU scheduling... ")
	if err := gpu.EnableHardwareScheduling(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[5/9] Optimizing network settings... ")
	if err := network.OptimizeForGaming(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[6/9] Setting High Performance power plan... ")
	if err := power.SetHighPerformance(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[7/9] Optimizing visual effects... ")
	if err := visual.OptimizeForPerformance(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[8/9] Applying CS2 launch options... ")
	if err := launchopts.ApplyLaunchOptions(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
		if strings.Contains(err.Error(), "Steam is currently running") {
			fmt.Println("      Please close Steam and run install again to apply launch options.")
		}
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[9/9] Applying CS2 config tweaks... ")
	if err := cs2config.ApplyOptimizedConfig(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n      (CS2 not found)\n", err)
	} else {
		fmt.Println("OK")
	}

	backupPath, _ := backup.GetBackupPath()
	fmt.Println()
	fmt.Println("Installation complete!")
	fmt.Printf("Backup saved to: %s\n", backupPath)
	fmt.Println("Please restart your computer for changes to take effect.")
}

func handleUninstall() {
	if !backup.Exists() {
		fmt.Println("Warning: No backup file found.")
		fmt.Println("Cannot restore previous settings safely.")
		fmt.Print("Continue with default uninstall? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Uninstallation cancelled.")
			return
		}
		handleUninstallDefault()
		return
	}

	fmt.Println("Restoring previous settings from backup...")
	fmt.Println()

	backupData, err := backup.Load()
	if err != nil {
		fmt.Printf("Error loading backup: %v\n", err)
		fmt.Println("Falling back to default uninstall.")
		handleUninstallDefault()
		return
	}

	fmt.Print("[1/9] Restoring cs2.exe priority... ")
	if err := priority.RestoreFromBackup(backupData.Priority); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[2/9] Restoring HPET setting... ")
	if err := hpet.RestoreFromBackup(backupData.HPET); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[3/9] Restoring Game DVR/Game Bar... ")
	if err := gamedvr.RestoreFromBackup(backupData.GameDVR); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[4/9] Restoring GPU scheduling... ")
	if err := gpu.RestoreFromBackup(backupData.GPU); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[5/9] Restoring network settings... ")
	if err := network.RestoreFromBackup(backupData.Network); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[6/9] Restoring power plan... ")
	if err := power.RestoreFromBackup(backupData.Power); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[7/9] Restoring visual effects... ")
	if err := visual.RestoreFromBackup(backupData.Visual); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[8/9] Restoring CS2 launch options... ")
	if err := launchopts.RestoreFromBackup(backupData.LaunchOpts); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[9/9] Restoring CS2 config... ")
	if err := cs2config.RestoreFromBackup(backupData.CS2Config); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("\nDeleting backup file... ")
	if err := backup.Delete(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Println()
	fmt.Println("Restoration complete!")
	fmt.Println("All settings have been restored to their previous state.")
	fmt.Println("Please restart your computer for changes to take effect.")
}

func handleUninstallDefault() {
	fmt.Println("Removing CS2 optimizations (default mode)...")
	fmt.Println()

	fmt.Print("[1/7] Removing cs2.exe priority setting... ")
	if err := priority.RemoveCS2Priority(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[2/7] Re-enabling HPET... ")
	if err := hpet.Enable(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[3/7] Re-enabling Game DVR/Game Bar... ")
	if err := gamedvr.Enable(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[4/7] Disabling Hardware GPU scheduling... ")
	if err := gpu.DisableHardwareScheduling(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[5/7] Resetting network settings... ")
	if err := network.ResetOptimizations(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Print("[6/7] Resetting power plan... ")
	fmt.Println("SKIPPED (manual reset recommended)")

	fmt.Print("[7/7] Resetting visual effects... ")
	if err := visual.ResetToDefault(); err != nil {
		fmt.Printf("FAILED\n      Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fmt.Println()
	fmt.Println("Uninstallation complete!")
	fmt.Println("Please restart your computer for changes to take effect.")
}

func handleStatus() {
	fmt.Println("CS2 Optimization Status")
	fmt.Println("=======================")
	fmt.Println()

	priorityStatus, err := priority.GetStatus()
	if err != nil {
		fmt.Printf("CS2 Priority:     ERROR (%v)\n", err)
	} else {
		if priorityStatus {
			fmt.Println("CS2 Priority:     ✓ ENABLED (High)")
		} else {
			fmt.Println("CS2 Priority:     ✗ DISABLED")
		}
	}

	hpetStatus, err := hpet.GetStatus()
	if err != nil {
		fmt.Printf("HPET:             ERROR (%v)\n", err)
	} else {
		if hpetStatus {
			fmt.Println("HPET:             ✗ ENABLED")
		} else {
			fmt.Println("HPET:             ✓ DISABLED (Optimized)")
		}
	}

	gameDVRStatus, err := gamedvr.GetStatus()
	if err != nil {
		fmt.Printf("Game DVR:         ERROR (%v)\n", err)
	} else {
		if gameDVRStatus {
			fmt.Println("Game DVR:         ✗ ENABLED")
		} else {
			fmt.Println("Game DVR:         ✓ DISABLED (Optimized)")
		}
	}

	gpuStatus, err := gpu.GetStatus()
	if err != nil {
		fmt.Printf("GPU Scheduling:   ERROR (%v)\n", err)
	} else {
		if gpuStatus {
			fmt.Println("GPU Scheduling:   ✓ ENABLED (Hardware)")
		} else {
			fmt.Println("GPU Scheduling:   ✗ DISABLED")
		}
	}

	networkStatus, err := network.GetStatus()
	if err != nil {
		fmt.Printf("Network:          ERROR (%v)\n", err)
	} else {
		if networkStatus {
			fmt.Println("Network:          ✓ OPTIMIZED (Gaming)")
		} else {
			fmt.Println("Network:          ✗ DEFAULT")
		}
	}

	powerStatus, err := power.GetStatus()
	if err != nil {
		fmt.Printf("Power Plan:       ERROR (%v)\n", err)
	} else {
		if powerStatus {
			fmt.Println("Power Plan:       ✓ HIGH PERFORMANCE")
		} else {
			plan, _ := power.GetActivePlan()
			fmt.Printf("Power Plan:       ✗ %s\n", plan)
		}
	}

	visualStatus, err := visual.GetStatus()
	if err != nil {
		fmt.Printf("Visual Effects:   ERROR (%v)\n", err)
	} else {
		if visualStatus {
			fmt.Println("Visual Effects:   ✓ OPTIMIZED (Performance)")
		} else {
			fmt.Println("Visual Effects:   ✗ DEFAULT")
		}
	}

	launchOptsStatus, err := launchopts.GetStatus()
	if err != nil {
		fmt.Printf("Launch Options:   ERROR (%v)\n", err)
	} else {
		if launchOptsStatus {
			fmt.Println("Launch Options:   ✓ OPTIMIZED")
		} else {
			fmt.Println("Launch Options:   ✗ DEFAULT")
		}
	}

	cs2ConfigStatus, err := cs2config.GetStatus()
	if err != nil {
		fmt.Printf("CS2 Config:       ERROR (%v)\n", err)
	} else {
		if cs2ConfigStatus {
			fmt.Println("CS2 Config:       ✓ OPTIMIZED")
		} else {
			fmt.Println("CS2 Config:       ✗ DEFAULT")
		}
	}
}
