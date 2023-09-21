package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message/unpack"
	"github.com/gotd/td/tg"
)

func check3() {
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
		req := tg.ChannelsCreateChannelRequest{
			Title: "test",
			About: "test about",
		}
		res, e := raw.ChannelsCreateChannel(ctx, &req)
		msg, err := unpack.Message(res, e)
		fmt.Println(msg)
		//req_photo := tg.ChannelsEditPhotoRequest{}

		fmt.Println(res, e)
		//raw.ChannelsEditPhoto()

		return err
	}); err != nil {
		panic(err)
	}
}
