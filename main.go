package main

import (
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/gliderlabs/ssh"
)

var addr = ":2222"

func init() {
	flag.StringVar(&addr, "addr", addr, "address to listen")
}

func main() {
	flag.Parse()

	slog.Info("start listening",
		"addr", addr,
	)

	sever := &ssh.Server{
		Version: "OpenSSH_8.0",
		Addr:    addr,
		Handler: func(session ssh.Session) {
			_, _ = fmt.Fprintln(session, "Fuck you")
		},
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			sessionId := ctx.SessionID()
			if len(sessionId) > 8 {
				sessionId = sessionId[:8]
			}
			slog.Info("new login",
				"session_id", sessionId,
				"user", ctx.User(),
				"password", password,
				"client_version", ctx.ClientVersion(),
				"remote_addr", ctx.RemoteAddr(),
			)
			select {
			case <-ctx.Done():
				return false
			case <-time.After(2 * time.Second):
				return false
			}
		},
	}
	if err := sever.ListenAndServe(); err != nil {
		slog.Error("listen and serve failed",
			"error", err,
		)
	}
}