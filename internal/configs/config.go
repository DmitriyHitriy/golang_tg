package configs

import (
	"github.com/gookit/ini"
)

type Configs struct {
	channel_name         string
	channel_desc         string
	channel_photo        string
	offer_text           string
	offer_photo          string
	parser_auditory_list string
	parser_content_list  string
}

func (c *Configs) New() *Configs {
	cfg, err := ini.LoadFiles("config.ini")
	if err != nil {
		panic(err)
	}
	channel_name, _ := cfg.String("channel_name")
	channel_desc, _ := cfg.String("channel_desc")
	channel_photo, _ := cfg.String("channel_photo")
	offer_text, _ := cfg.String("offer_text")
	offer_photo, _ := cfg.String("offer_photo")
	parser_auditory_list, _ := cfg.String("parser_auditory_list")
	parser_content_list, _ := cfg.String("parser_content_list")

	if channel_name == "" || channel_desc == "" {
		panic("channel_name and channel_desc")
	}

	c.setChannelName(channel_name)
	c.setChannelDesc(channel_desc)
	c.setChannelPhoto(channel_photo)
	c.setOfferText(offer_text)
	c.setOfferPhoto(offer_photo)
	c.setParserAuditoryList(parser_auditory_list)
	c.setParserContentList(parser_content_list)

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

func (c *Configs) GetOfferText() string {
	return c.offer_text
}

func (c *Configs) GetOfferPhoto() string {
	return c.offer_photo
}

func (c *Configs) GetParserAuditoryList() string {
	return c.parser_auditory_list
}

func (c *Configs) GetParserContentList() string {
	return c.parser_content_list
}

func (c *Configs) setChannelName(channel_name string) {
	c.channel_name = channel_name
}

func (c *Configs) setChannelDesc(channel_desc string) {
	c.channel_desc = channel_desc
}

func (c *Configs) setChannelPhoto(channel_photo string) {
	c.channel_photo = channel_photo
}

func (c *Configs) setOfferText(offer_text string) {
	c.offer_text = offer_text
}

func (c *Configs) setOfferPhoto(offer_photo string) {
	c.offer_photo = offer_photo
}

func (c *Configs) setParserAuditoryList(parser_auditory_list string) {
	c.parser_auditory_list = parser_auditory_list
}

func (c *Configs) setParserContentList(parser_content_list string) {
	c.parser_content_list = parser_content_list
}