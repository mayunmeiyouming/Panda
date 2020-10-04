package cmd

import "net/url"

var addr string
var cipher string
var password string
var tcp bool
var udp bool
var plugin string
var pluginOpts string
var socks string
var localAddr string
var remoteAddr string
var udpsocks bool

func parseURL(s string) (string, string, string, error) {
	var addr, cipher, password string
	u, err := url.Parse(s)
	if err != nil {
		return "", "", "", err
	}

	addr = u.Host
	if u.User != nil {
		cipher = u.User.Username()
		password, _ = u.User.Password()
	}
	return addr, cipher, password, nil
}
