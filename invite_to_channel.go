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

func check7() {
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
		ch := &tg.InputChannel{ChannelID: 1905046891, AccessHash: 5725182504979867499}
		us := &tg.InputUser{UserID: 5419321810, AccessHash: -5698706159390738840}

		//users := []tg.InputUser{*us}

		req := tg.ChannelsInviteToChannelRequest{ch, []tg.InputUserClass{us}}
		res, e := raw.ChannelsInviteToChannel(ctx, &req)

		fmt.Println(res)
		fmt.Println(e)

		return err
	}); err != nil {
		panic(err)
	}
}
