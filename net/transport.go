package net

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"

	"region-exporter/config"
)

func CreateTransport(cfg *config.CLI) (*http.Transport, error) {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	if cfg.Network.Interface != "" {
		iface, err := net.InterfaceByName(cfg.Network.Interface)
		if err != nil {
			return nil, fmt.Errorf("failed to find interface %s: %w", cfg.Network.Interface, err)
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("failed to get addresses for interface %s: %w", cfg.Network.Interface, err)
		}

		if len(addrs) == 0 {
			return nil, fmt.Errorf("interface %s has no addresses", cfg.Network.Interface)
		}

		addr := addrs[0]
		if ipNet, ok := addr.(*net.IPNet); ok {
			dialer.LocalAddr = &net.TCPAddr{IP: ipNet.IP}
			if cfg.Verbose {
				fmt.Printf("Using local address from interface %s: %s\n", cfg.Network.Interface, ipNet.IP)
			}
		}
	}

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		if cfg.Network.IPv4Only {
			network = "tcp4"
		} else if cfg.Network.IPv6Only {
			network = "tcp6"
		}
		return dialer.DialContext(ctx, network, addr)
	}

	if cfg.Network.Proxy != "" {
		if cfg.Verbose {
			fmt.Printf("Configuring SOCKS5 proxy: %s\n", cfg.Network.Proxy)
		}

		socksDialer, err := proxy.SOCKS5("tcp", cfg.Network.Proxy, nil, dialer)
		if err != nil {
			return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
		}

		if contextDialer, ok := socksDialer.(proxy.ContextDialer); ok {
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				if cfg.Network.IPv4Only {
					network = "tcp4"
				} else if cfg.Network.IPv6Only {
					network = "tcp6"
				}
				return contextDialer.DialContext(ctx, network, addr)
			}
		} else {
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				if cfg.Network.IPv4Only {
					network = "tcp4"
				} else if cfg.Network.IPv6Only {
					network = "tcp6"
				}
				return socksDialer.Dial(network, addr)
			}
		}
	} else {
		transport.DialContext = dialContext
	}

	return transport, nil
}
