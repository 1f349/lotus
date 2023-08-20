package main

import (
	"github.com/1f349/lotus/imap"
	"github.com/1f349/lotus/smtp"
)

type Conf struct {
	Listen   string     `yaml:"listen"`
	Audience string     `yaml:"audience"`
	Smtp     *smtp.Smtp `yaml:"smtp"`
	Imap     *imap.Imap `yaml:"imap"`
}
