package config

import (
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	T "github.com/Dreamacro/clash/tunnel"
)

type JsonConfig struct {
	Port               int          `json:"port"`
	SocksPort          int          `json:"socks-port"`
	RedirPort          int          `json:"redir-port"`
	TProxyPort         int          `json:"tproxy-port"`
	MixedPort          int          `json:"mixed-port"`
	Authentication     []string     `json:"authentication"`
	AllowLan           bool         `json:"allow-lan"`
	BindAddress        string       `json:"bind-address"`
	Mode               T.TunnelMode `json:"mode"`
	LogLevel           log.LogLevel `json:"log-level"`
	IPv6               bool         `json:"ipv6"`
	ExternalController string       `json:"external-controller"`
	ExternalUI         string       `json:"external-ui"`
	Secret             string       `json:"secret"`
	Interface          string       `json:"interface-name"`
	RoutingMark        int          `json:"routing-mark"`

	ProxyProvider map[string]map[string]any `json:"proxy-providers"`
	Hosts         map[string]string         `json:"hosts"`
	DNS           JsonRawDNS                `json:"dns"`
	Experimental  Experimental              `json:"experimental"`
	Profile       JsonProfile               `json:"profile"`
	Proxy         []map[string]any          `json:"proxies"`
	ProxyGroup    []map[string]any          `json:"proxy-groups"`
	Rule          []string                  `json:"rules"`
}

type JsonRawDNS struct {
	Enable            bool                  `json:"enable"`
	IPv6              bool                  `json:"ipv6"`
	UseHosts          bool                  `json:"use-hosts"`
	NameServer        []string              `json:"nameserver"`
	Fallback          []string              `json:"fallback"`
	FallbackFilter    JsonRawFallbackFilter `json:"fallback-filter"`
	Listen            string                `json:"listen"`
	EnhancedMode      C.DNSMode             `json:"enhanced-mode"`
	FakeIPRange       string                `json:"fake-ip-range"`
	FakeIPFilter      []string              `json:"fake-ip-filter"`
	DefaultNameserver []string              `json:"default-nameserver"`
	NameServerPolicy  map[string]string     `json:"nameserver-policy"`
}

type JsonProfile struct {
	StoreSelected bool   `json:"store-selected"`
	StoreFakeIP   bool   `json:"store-fake-ip"`
	UiStorage     string `json:"ui-storage"`
}

type JsonRawFallbackFilter struct {
	GeoIP     bool     `json:"geoip"`
	GeoIPCode string   `json:"geoip-code"`
	IPCIDR    []string `json:"ipcidr"`
	Domain    []string `json:"domain"`
}

func (cfg *RawConfig) BuildJson() *JsonConfig {
	config := &JsonConfig{
		Port:               cfg.Port,
		SocksPort:          cfg.SocksPort,
		RedirPort:          cfg.RedirPort,
		TProxyPort:         cfg.TProxyPort,
		MixedPort:          cfg.MixedPort,
		Authentication:     cfg.Authentication,
		AllowLan:           cfg.AllowLan,
		BindAddress:        cfg.BindAddress,
		Mode:               cfg.Mode,
		LogLevel:           cfg.LogLevel,
		IPv6:               cfg.IPv6,
		ExternalController: cfg.ExternalController,
		ExternalUI:         cfg.ExternalUI,
		Secret:             cfg.Secret,
		Interface:          cfg.Interface,
		RoutingMark:        cfg.RoutingMark,

		ProxyProvider: cfg.ProxyProvider,
		Hosts:         cfg.Hosts,
		DNS: JsonRawDNS{
			Enable:     cfg.DNS.Enable,
			IPv6:       cfg.DNS.IPv6,
			UseHosts:   cfg.DNS.UseHosts,
			NameServer: cfg.DNS.NameServer,
			Fallback:   cfg.DNS.Fallback,
			FallbackFilter: JsonRawFallbackFilter{
				GeoIP:     cfg.DNS.FallbackFilter.GeoIP,
				GeoIPCode: cfg.DNS.FallbackFilter.GeoIPCode,
				IPCIDR:    cfg.DNS.FallbackFilter.IPCIDR,
				Domain:    cfg.DNS.FallbackFilter.Domain,
			},
			Listen:            cfg.DNS.Listen,
			EnhancedMode:      cfg.DNS.EnhancedMode,
			FakeIPRange:       cfg.DNS.FakeIPRange,
			FakeIPFilter:      cfg.DNS.FakeIPFilter,
			DefaultNameserver: cfg.DNS.DefaultNameserver,
			NameServerPolicy:  cfg.DNS.NameServerPolicy,
		},
		Experimental: cfg.Experimental,
		Profile: JsonProfile{
			StoreSelected: cfg.Profile.StoreSelected,
			StoreFakeIP:   cfg.Profile.StoreFakeIP,
			UiStorage:     cfg.Profile.UiStorage,
		},
		Proxy:      cfg.Proxy,
		ProxyGroup: cfg.ProxyGroup,
		Rule:       cfg.Rule,
	}
	return config
}
