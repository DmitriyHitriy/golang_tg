package channel

import (
	"fmt"
	"context"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/telegram"
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