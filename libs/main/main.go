package WysbWifi

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type WifiManager struct{}

func NewWifiManager() *WifiManager {
	return &WifiManager{}
}

func (wm *WifiManager) ListNetworks() ([]string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("netsh", "wlan", "show", "network")
	} else {
		cmd = exec.Command("nmcli", "dev", "wifi")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseNetworkList(string(output)), nil
}

func (wm *WifiManager) KillNetwork(ssid string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("netsh", "wlan", "disconnect", "name="+ssid)
	} else {
		cmd = exec.Command("nmcli", "dev", "disconnect", ssid)
	}

	_, err := cmd.Output()
	return err
}

func parseNetworkList(output string) []string {

	lines := strings.Split(output, "\n")
	var networks []string
	for _, line := range lines {
		if strings.Contains(line, "SSID") {

			network := strings.TrimSpace(line)
			networks = append(networks, network)
		}
	}
	return networks
}

func Disclaimer() {
	fmt.Println("Esta biblioteca é para fins educacionais e éticos. O criador não se responsabiliza por mau uso.")
}