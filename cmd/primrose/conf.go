package main

import (
	"github.com/1f349/primrose/api"
	"github.com/1f349/primrose/imap"
	"github.com/1f349/primrose/smtp"
)

type Conf struct {
	Smtp smtp.Smtp `yaml:"smtp"`
	Imap imap.Imap `yaml:"imap"`
	Api  api.Conf  `yaml:"api"`
}
