package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func get_participants_in_channel(channel_id int64, access_hash int64) {
	ctx := context.Background()
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	root := filepath.Join(home, "Documents", "Telegram", "tdata")
	accounts, err := tdesktop.Read(root, nil)
	if err != nil {
		panic(err)
	}

	data, err := session.TDesktopSession(accounts[0])
	if err != nil {
		panic(err)
	}

	var (
		storage = new(session.StorageMemory)
		loader  = session.Loader{Storage: storage}
	)

	if err := loader.Save(ctx, data); err != nil {
		panic(err)
	}

	client := telegram.NewClient(1, "s", telegram.Options{SessionStorage: storage})

	if err := client.Run(ctx, func(ctx context.Context) error {
		raw := tg.NewClient(client)
		//ch := &tg.InputChannel{ChannelID: 1337318424, AccessHash: 4419570524830746652}
		ch := &tg.InputChannel{ChannelID: channel_id, AccessHash: access_hash}

		res, e := raw.ChannelsJoinChannel(ctx, ch)

		fmt.Println(res)
		fmt.Println(e)

		req := tg.ChannelsGetParticipantsRequest{
			Channel: ch,
			Limit:   100,
			Filter:  &tg.ChannelParticipantsRecent{},
		}
		res_part, ee := raw.ChannelsGetParticipants(ctx, &req)
		fmt.Println(res_part)
		fmt.Println(ee)

		return err
	}); err != nil {
		panic(err)
	}
}
