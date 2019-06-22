package astislack

import (
	"flag"

	"github.com/asticode/go-astitools/http"
)

// Flags
var (
	LegacyToken = flag.String("slack-legacy-token", "", "the Slack legacy token")
)

// Configuration represents the configuration of the logger
type Configuration struct {
	LegacyToken string `toml:"legacy_token"`
	Sender      astihttp.SenderOptions
}

// FlagConfig generates a Configuration based on flags
func FlagConfig() Configuration {
	return Configuration{
		LegacyToken: *LegacyToken,
	}
}
