package configs

import (
	"github.com/gookit/ini"
)

type Configs struct {
	channel_name string
	channel_desc string
	channel_photo string
}

func (c *Configs) New() *Configs {
	cfg, err := ini.LoadFiles("config.ini")
	if err != nil {
		panic(err)
	}
	channel_name, _ := cfg.String("channel_name")
	channel_desc, _ := cfg.String("channel_desc")
	channel_photo, _ := cfg.String("channel_photo")

	if channel_name == "" || channel_desc == "" {
		panic("channel_name and channel_desc")
	}

	c.setChannelName(channel_name)
	c.setChannelDesc(channel_desc)
	c.serChannelPhoto(channel_photo)

	return c
}

func (c *Configs) GetChannelName() string {
	return c.channel_name
}

func (c *Configs) GetChannelDesc() string {
	return c.channel_desc
}

func (c *Configs) GetChannelPhoto() string {
	return c.channel_photo
}

func (c *Configs) setChannelName(channel_name string) {
	c.channel_name = channel_name
}

func (c *Configs) setChannelDesc(channel_desc string) {
	c.channel_desc = channel_desc
}

func (c *Configs) serChannelPhoto(channel_photo string) {
	c.channel_photo = channel_photo
}