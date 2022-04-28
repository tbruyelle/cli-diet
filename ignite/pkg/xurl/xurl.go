package xurl

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
)

const (
	schemeTCP   = "tcp"
	schemeHTTP  = "http"
	schemeHTTPS = "https"
	schemeWS    = "ws"
)

// TCP ensures that a URL contains a TCP scheme.
func TCP(s string) (string, error) {
	u, err := parseURL(s)
	if err != nil {
		return "", err
	}

	u.Scheme = schemeTCP

	return u.String(), nil
}

// HTTP ensures that a URL contains an HTTP scheme.
func HTTP(s string) (string, error) {
	u, err := parseURL(s)
	if err != nil {
		return "", err
	}

	u.Scheme = schemeHTTP

	return u.String(), nil
}

// HTTPS ensures that a URL contains an HTTPS scheme.
func HTTPS(s string) (string, error) {
	u, err := parseURL(s)
	if err != nil {
		return "", err
	}

	u.Scheme = schemeHTTPS

	return u.String(), nil
}

// MightHTTPS ensures that a URL contains an HTTPS scheme when the current scheme is not HTTP.
// When the URL contains an HTTP scheme it is not modified.
func MightHTTPS(s string) (string, error) {
	if strings.HasPrefix(strings.ToLower(s), "http://") {
		return s, nil
	}

	return HTTPS(s)
}

// WS ensures that a URL contains a WS scheme.
func WS(s string) (string, error) {
	u, err := parseURL(s)
	if err != nil {
		return "", err
	}

	u.Scheme = schemeWS

	return u.String(), nil
}

// HTTPEnsurePort ensures that url has a port number suits with the connection type.
func HTTPEnsurePort(s string) string {
	u, err := url.Parse(s)
	if err != nil || u.Port() != "" {
		return s
	}

	port := "80"

	if u.Scheme == schemeHTTPS {
		port = "443"
	}

	u.Host = fmt.Sprintf("%s:%s", u.Hostname(), port)

	return u.String()
}

// CleanPath cleans path from the url.
func CleanPath(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}

	u.Path = ""

	return u.String()
}

// Address ensures that address contains localhost as host if non specified.
func Address(address string) string {
	if strings.HasPrefix(address, ":") {
		return "localhost" + address
	}
	return address
}

// IsLocalPath checks if given address is a local fs path or a URL.
func IsLocalPath(address string) bool {
	for _, pattern := range []string{
		"http://",
		"https://",
		"git@",
	} {
		if strings.HasPrefix(address, pattern) {
			return false
		}
	}
	return true
}

func IsHTTP(address string) bool {
	return strings.HasPrefix(address, "http")
}

func parseURL(s string) (*url.URL, error) {
	if s == "" {
		return nil, errors.New("url is empty")
	}

	// Handle the case where the URI is an IP:PORT or HOST:PORT
	// without scheme prefix because that case can't be URL parsed.
	// When the URI has not scheme it is parsed as a path by "url.Parse"
	// placing the colon within the path, which is invalid.
	if isAddressPort(s) {
		return &url.URL{Host: s}, nil
	}

	return url.Parse(Address(s))
}

func isAddressPort(s string) bool {
	// Check that the value doesn't contain a URI path
	if strings.Index(s, "/") != -1 {
		return false
	}

	// Use the net split function to support IPv6 addresses
	if _, _, err := net.SplitHostPort(s); err != nil {
		return false
	}
	return true
}
