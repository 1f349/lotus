package main

import (
	"github.com/1f349/lotus/imap"
	"github.com/1f349/lotus/sendmail"
)

type Conf struct {
	Listen   string         `yaml:"listen"`
	Audience string         `yaml:"audience"`
	Smtp     *sendmail.Smtp `yaml:"sendmail"`
	Imap     *imap.Imap     `yaml:"imap"`
}
