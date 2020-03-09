package httpx

import (
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"
)

// AccessConfig configures what can be accessed
type AccessConfig struct {
	ResolveTimeout time.Duration
	DisallowedIPs  []net.IP
}

// NewAccessConfig creates a new access config
func NewAccessConfig(resolveTimeout time.Duration, disallowedIPs []net.IP) *AccessConfig {
	return &AccessConfig{
		ResolveTimeout: resolveTimeout,
		DisallowedIPs:  disallowedIPs,
	}
}

// Allow determines whether the given request should be allowed
func (c *AccessConfig) Allow(request *http.Request) (bool, error) {
	host := strings.ToLower(request.URL.Hostname())

	ctx, cancel := context.WithTimeout(context.Background(), c.ResolveTimeout)
	defer cancel()

	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return false, err
	}

	// if any of the host's addresses appear in the disallowed list, deny the request
	for _, addr := range addrs {
		for _, disallowed := range c.DisallowedIPs {
			if addr.IP.Equal(disallowed) {
				return false, nil
			}
		}
	}
	return true, nil
}
