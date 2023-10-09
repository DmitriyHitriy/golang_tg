package channel

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/html"
	"github.com/gotd/td/tg"
)

type Channel struct {
	InputChannel *tg.InputChannel
}

func (c *Channel) GetChannel() *tg.InputChannel {
	return c.InputChannel
}

func (c *Channel) SetChannel(ch *tg.InputChannel) {
	c.InputChannel = ch
}

func (c *Channel) InviteToChannel(ctx context.Context, client *telegram.Client, channel *tg.InputChannel, user *tg.User) {
	if err := client.Run(ctx, func(ctx context.Context) error {
		raw := tg.NewClient(client)

		us := &tg.InputUser{UserID: user.GetID(), AccessHash: user.AccessHash}

		req := tg.ChannelsInviteToChannelRequest{
			Channel: channel,
			Users:   []tg.InputUserClass{us},
		}

		res, e := raw.ChannelsInviteToChannel(ctx, &req)

		if e != nil {
			fmt.Println("Ошибка добавления в канал ", e)
		}

		fmt.Println(res)
		fmt.Println(e)

		return e
	}); err != nil {
		fmt.Println("Ошибка добавления в канал ", err)
	}

}

func (c *Channel) CreatePost(ctx context.Context, client *telegram.Client, channel *tg.InputChannel) {

}

// TODO: Написать нормальный обработчик получения информации
func (c *Channel) GetChannelInfo(ctx context.Context, client *telegram.Client, channel *tg.InputChannel){
	if err := client.Run(ctx, func(ctx context.Context) error {
		var err error
		raw := tg.NewClient(client)
		
		// req := tg.ChannelsGetChannelsRequest{
		// 	ID: []tg.InputChannelClass{channel},
		// }

		res, e := raw.ChannelsGetChannels(ctx, []tg.InputChannelClass{channel})

		fmt.Println(res)
		fmt.Println(e)
		return err
	}); err != nil {
		panic(err)
	}
}

func (c *Channel) ChannelSendMessage(ctx context.Context, client *telegram.Client, channel string) {
	if err := client.Run(ctx, func(ctx context.Context) error {
		var err error
		raw := tg.NewClient(client)

		nm, _ := message.NewSender(raw).Resolve(channel).Upload(message.Upload(func(ctx context.Context, b message.Uploader) (tg.InputFileClass, error) {
			r, err := b.FromPath(ctx, "skull.jpg")
			if err != nil {
				return nil, err
			}

			return r, nil
		})).Photo(ctx, html.String(nil, "<b>Hello</b>, world! This is <a href=\"http://google.com\">address</a>)"))

		fmt.Println(nm)
		fmt.Println(err)

		return err
	}); err != nil {
		panic(err)
	}
	
}