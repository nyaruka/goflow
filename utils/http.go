package utils

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

func init() {
	Validator.RegisterAlias("http_method", "eq=GET|eq=HEAD|eq=POST|eq=PUT|eq=PATCH|eq=DELETE")
}

var (
	transport *http.Transport
	client    *http.Client
	once      sync.Once
)

// NewHTTPClient creates a new hTTP cl.ient with our default options
func NewHTTPClient() *http.Client {
	once.Do(func() {
		timeout := time.Duration(15 * time.Second)
		transport = &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		}
		client = &http.Client{Transport: transport, Timeout: timeout}
	})

	return client
}

func getInsecureClient() *http.Client {
	once.Do(func() {
		timeout := time.Duration(15 * time.Second)
		transport = &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: transport, Timeout: timeout}
	})

	return client
}
