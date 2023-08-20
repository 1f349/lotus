package postfix_config

import mapProvider "github.com/1f349/lotus/postfix-config/map-provider"

type Config struct {
	// same
	VirtualMailboxDomains mapProvider.MapProvider
	VirtualAliasMaps      mapProvider.MapProvider
	VirtualMailboxMaps    mapProvider.MapProvider
	AliasMaps             mapProvider.MapProvider
	LocalRecipientMaps    mapProvider.MapProvider
	SmtpdSenderLoginMaps  string // TODO(melon): union map?
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
		c.SmtpdSenderLoginMaps = "<ERROR>"
	}
}

func (c *Config) NeedsMapProvider(k string) bool {
	switch k {
	case "virtual_mailbox_domains", "virtual_alias_maps", "virtual_mailbox_maps", "alias_maps", "local_recipient_maps", "smtpd_sender_login_maps":
		return true
	}
	return false
}
