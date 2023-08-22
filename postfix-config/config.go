package postfix_config

import mapProvider "github.com/1f349/lotus/postfix-config/map-provider"

type Config struct {
	// same
	VirtualMailboxDomains mapProvider.MapProvider
	VirtualAliasMaps      mapProvider.MapProvider
	VirtualMailboxMaps    mapProvider.MapProvider
	AliasMaps             mapProvider.MapProvider
	LocalRecipientMaps    mapProvider.MapProvider
	SmtpdSenderLoginMaps  mapProvider.MapProvider
}

var parseProviderData = map[string]string{
	"virtual_mailbox_domains": "comma",
	"virtual_alias_maps":      "comma",
	"virtual_mailbox_maps":    "comma",
	"alias_maps":              "comma",
	"local_recipient_maps":    "comma",
	"smtpd_sender_login_maps": "union",
}

func (c *Config) ParseProvider(k string) string {
	return parseProviderData[k]
}

func (c *Config) SetKey(k string, m mapProvider.MapProvider) {
	switch k {
	case "virtual_mailbox_domains":
		c.VirtualMailboxDomains = m
	case "virtual_alias_maps":
		c.VirtualAliasMaps = m
	case "virtual_mailbox_maps":
		c.VirtualMailboxMaps = m
	case "alias_maps":
		c.AliasMaps = m
	case "local_recipient_maps":
		c.LocalRecipientMaps = m
	case "smtpd_sender_login_maps":
		c.SmtpdSenderLoginMaps = m
	}
}
