package main

import (
	"log/slog"
	"net"
	"runtime"
	"time"

	"github.com/micmonay/keybd_event"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		slog.Warn("cannot get InterfaceAddrs", slog.Any("error", err))
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// For some reason YouTube is not responding without user interaction, so we simulate 2x F12 keypress
func hackKeyboadEvent() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		slog.Error("error on keyboard input", slog.Any("error", err))
		return
	}

	// For linux, it is very important to wait 2 seconds - that's what docs say, but in my case 0.5 is enough
	if runtime.GOOS == "linux" {
		time.Sleep(500 * time.Millisecond)
	}

	kb.SetKeys(keybd_event.VK_DOT)

	t := time.NewTicker(5 * time.Second)

	for {
		<-t.C
		err = kb.Launching()
		if err != nil {
			slog.Error("error on keyboard input", slog.Any("error", err))
		}

	}
}
