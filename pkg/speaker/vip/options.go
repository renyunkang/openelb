package vip

import (
	"github.com/spf13/pflag"
)

type VipOptions struct {
	EnableVIP bool
	LogPath   string
}

func NewVipOptions() *VipOptions {
	return &VipOptions{
		EnableVIP: false,
		LogPath:   "",
	}
}

func (v *VipOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&v.EnableVIP, "enable-keepalived-vip", v.EnableVIP, "specify whether to start keepalived-vip")
	fs.StringVar(&v.LogPath, "log-path", v.LogPath, "specify the path of the keepalived log file")
}
