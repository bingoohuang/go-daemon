package autoload

import (
	"github.com/bingoohuang/godaemon"
	"os"
	"strings"
)

func init() {
	if env := os.Getenv("DAEMON"); env != "" {
		switch strings.ToLower(env) {
		case "true", "1", "t", "yes", "y", "on":
			godaemon.Daemonize()
		}
	}
}
