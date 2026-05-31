package steam

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsSteamRunning() bool {
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq steam.exe")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToLower(string(output)), "steam.exe")
}

func WarnIfSteamRunning() error {
	if IsSteamRunning() {
		return fmt.Errorf("Steam is currently running. Please close Steam before applying launch options")
	}
	return nil
}
