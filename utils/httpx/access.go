package httpx

import (
	"net"
	"net/http"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

// AccessConfig configures what can be accessed
type AccessConfig struct {
	DisallowedIPs []string
}

func NewAccessConfig(disallowedIPs []string) *AccessConfig {
	return &AccessConfig{DisallowedIPs: disallowedIPs}
}

func (c *AccessConfig) Allow(request *http.Request) bool {
	host := strings.ToLower(request.URL.Hostname())

	// if host looks like an IP address, normalize it
	asIP := net.ParseIP(host)
	if asIP != nil {
		host = asIP.String()
	}

	return !utils.StringSliceContains(c.DisallowedIPs, host, true)
}
