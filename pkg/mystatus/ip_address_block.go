package mystatus

import (
	"fmt"
	"os"
	"strings"

	"github.com/shirou/gopsutil/net"
)

type ipAddressBlock struct {
	interfaceName  string
	interfaceNames []string
}

func (cb *ipAddressBlock) hadInterface(name string) bool {
	for _, i := range cb.interfaceNames {
		if i == name {
			return true
		}
	}
	return false
}

func (cb *ipAddressBlock) Render() barBlockData {
	fullText := ""
	interfaces, err := net.Interfaces()

	if err != nil {
		fullText = "IP error: " + err.Error()
	} else {
		for _, iface := range interfaces {
			text := ""
			if iface.Name == "lo" {
				continue
			}
			if _, err := os.Stat("/sys/class/net/" + iface.Name); err == nil {
				var emoji string

				if _, err := os.Stat("/sys/class/net/" + iface.Name + "/wireless"); err == nil { // check is is wiereless
					emoji += "📡"
				} else {
					emoji += "🔌"
				}
				isVirtual := false
				linkDest, err := os.Readlink(fmt.Sprintf("/sys/class/net/%s", iface.Name))
				if err == nil {
					if strings.Contains(linkDest, "/virtual/") {
						isVirtual = true
					}
				}
				if isVirtual {
					continue
				}
				if len(iface.Addrs) == 0 {
					text = "---"
				} else {
					text = strings.Split(iface.Addrs[0].Addr, "/")[0]
					text = strings.Replace(text, "192.168", `<span size="6000">192.168</span>`, -1)
				}

				fullText += emoji + text + " "
			}
		}
	}

	return barBlockData{
		Block:    cb,
		FullText: fullText,
		Markup:   "pango",
	}
}
