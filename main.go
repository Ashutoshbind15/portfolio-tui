package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
)

const defaultSSHPort = "23235"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "ssh" {
		if err := runSSH(); err != nil {
			log.Fatal("ssh server", "err", err)
		}
		return
	}
	if err := runLocal(); err != nil {
		log.Fatal("tui", "err", err)
	}
}

func runLocal() error {
	p := tea.NewProgram(newAppModel())
	_, err := p.Run()
	return err
}

func sshListenHost() string {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("TUI_PF_ENV"))) {
	case "prod", "production":
		return "0.0.0.0"
	default:
		return "127.0.0.1"
	}
}

func sshPort() string {
	if p := strings.TrimSpace(os.Getenv("TUI_PF_PORT")); p != "" {
		return p
	}
	return defaultSSHPort
}

func runSSH() error {
	host := sshListenHost()
	port := sshPort()

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				return newAppModel(), nil
			}),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		return err
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		return err
	}
	return nil
}
