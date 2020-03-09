package httpx

import (
	"net"
	"net/http"
	"strings"
)

// AccessConfig configures what can be accessed
type AccessConfig struct {
	DisallowedIPs []net.IP
}

// NewAccessConfig creates a new access config
func NewAccessConfig(disallowedIPs []net.IP) *AccessConfig {
	return &AccessConfig{DisallowedIPs: disallowedIPs}
}

// Allow determines whether the given request should be allowed
func (c *AccessConfig) Allow(request *http.Request) (bool, error) {
	host := strings.ToLower(request.URL.Hostname())

	addrs, err := net.LookupIP(host)
	if err != nil {
		return false, err
	}

	// if any of the host's addresses appear in the disallowed list, deny the request
	for _, addr := range addrs {
		for _, disallowed := range c.DisallowedIPs {
			if addr.Equal(disallowed) {
				return false, nil
			}
		}
	}
	return true, nil
}
