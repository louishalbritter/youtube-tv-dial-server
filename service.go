package main

import (
	"log/slog"
	"os/exec"
	"time"

	"github.com/koron/go-ssdp"
)

type Service struct {
	BaseUrl      string
	FriendlyName string
	Manufacture  string
	ModelName    string
	Uuid         string
	Proc         *exec.Cmd
}

func (s *Service) ssdp() {
	aliveTick := time.Tick(300 * time.Second)

	// Start SSDP advertisement
	ad, err := ssdp.Advertise(
		"urn:dial-multiscreen-org:service:dial:1", // USN (Unique Service Name)
		"uuid:"+s.Uuid,                    // UUID
		s.BaseUrl+"/ssdp/device-desc.xml", // Location of the description XML
		"ssdp:alive",                      // SSDP message type
		1800,                              // Duration for advertisement in seconds
	)
	if err != nil {
		slog.Error("Failed to advertise via SSDP", slog.Any("error", err))
	}

	defer ad.Close()

	for {
		select {
		case <-aliveTick:
			err := ad.Alive()
			if err != nil {
				slog.Warn("Failed keep ssdp alive", slog.Any("error", err))
			}
		}
	}
}

func (s *Service) start(args string) {
	if s.Proc != nil {
		if err := s.Proc.Cancel(); err != nil {
			slog.Error("Failed to Cancel process", slog.Any("error", err))
		}
	}

	p := exec.Command("xset", "dpms", "force", "on")
	if err := p.Run(); err != nil {
		slog.Error("Failed to turn on display", slog.Any("error", err))
	}

	s.Proc = exec.Command("brave-beta",
		"--incognito",
		"--user-agent=\"Mozilla/5.0 (X11; Linux i686) AppleWebKit/534.24 (KHTML, like Gecko) Chrome/11.0.696.77 Large Screen Safari/534.24 GoogleTV/092754\"",
		"--start-fullscreen",
		"--bwsi", // browse without signing
		"https://www.youtube.com/tv?"+args,
	)

	err := s.Proc.Start()
	if err != nil {
		slog.Error("Failed to launch Service", slog.Any("error", err))
	}
}

func (s *Service) stop() {
	if err := s.Proc.Cancel(); err != nil {
		slog.Error("Failed to Cancel process", slog.Any("error", err))
	}

	s.Proc = nil

	p := exec.Command("xset", "dpms", "force", "off")
	if err := p.Run(); err != nil {
		slog.Error("Failed to turn off display", slog.Any("error", err))
	}
}
