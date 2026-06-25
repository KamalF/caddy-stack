package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/lucaslorentz/caddy-docker-proxy/v2"
	_ "github.com/caddy-dns/ovh"
	_ "github.com/hslatman/caddy-crowdsec-bouncer"
)

func main() { caddycmd.Main() }
