package dyatl

import (
	"context"
	"net"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/idna"
)

var DefaultDNS = "8.8.8.8:53"

func NewCheckedURL(s string) *CheckedURL {
	var u CheckedURL
	u.URL, u.err = url.Parse(s)
	return &u
}

var dnsResolver = &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: time.Second,
		}
		return d.DialContext(ctx, "udp", DefaultDNS)
	},
}

type CheckedURL struct {
	*url.URL
	err error
}

func (u CheckedURL) LooksCorrect() bool {
	if u.err != nil {
		return false
	}

	host := u.Hostname()
	if hasMajorTLD(host) && (u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "ftp") && !strings.HasPrefix(host, ".") {
		return true
	}

	return net.ParseIP(host) != nil
}

func (u CheckedURL) HostAvailable() bool {
	if u.err != nil {
		return false
	}

	ips, _ := dnsResolver.LookupIPAddr(context.Background(), u.Hostname())
	if len(ips) > 0 {
		return true
	}

	if asciiHost := u.asciiHost(); asciiHost != "" {
		ips, _ := dnsResolver.LookupIPAddr(context.Background(), asciiHost)
		if len(ips) > 0 {
			return true
		}
	}

	return false
}

func (u CheckedURL) asciiHost() string {
	host := u.Hostname()
	asciiHost, err := idna.ToASCII(host)
	if err != nil || asciiHost == host {
		return ""
	}
	return asciiHost
}

func hasMajorTLD(s string) bool {
	s = strings.TrimSuffix(s, ".")
	if s == "" {
		return false
	}
	bits := strings.Split(s, ".")
	domain := strings.ToUpper(bits[len(bits)-1])
	_, ok := majorTLDs[domain]
	return ok
}
