package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/telegram/message"
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

		// req := tg.ChannelsCreateChannelRequest{
		// 	Title:     "test2",
		// 	About:     "test about",
		// 	Broadcast: true,
		// 	Megagroup: false,
		// }

		// res, e := raw.ChannelsCreateChannel(ctx, &req)

		// fmt.Println(e)
		// channel := ((*res.(*tg.Updates)).Chats[0]).(*tg.Channel)
		// fmt.Println(channel.ID, channel.AccessHash)

		// ch_input := tg.InputChannel{
		// 	ChannelID:  channel.ID,
		// 	AccessHash: channel.AccessHash,
		// }

		ch_input := tg.InputChannel{ChannelID: 1905046891, AccessHash: 5725182504979867499}
		
		//fmt.Println(ch_input)
		//f, _ := os.ReadFile("skull.jpg")
		upl := uploader.NewUploader(raw)
		//upl.WithPartSize(1024)
		//upl.WithThreads(8)
		fl, err_f := upl.FromPath(ctx, "skull2.jpg")

		in_photo := tg.InputPhoto{ID: (*fl.(*tg.InputFile)).ID}
		
		photo := tg.InputChatPhoto{ID: &in_photo}

		req_photo := tg.ChannelsEditPhotoRequest{
			Channel: &ch_input, 
			Photo: &photo,
		}
		
		change, err_upl_file := raw.ChannelsEditPhoto(ctx, &req_photo)
		fmt.Println(change)
		fmt.Println(err_upl_file)
		fmt.Println(fl)
		fmt.Println(err_f)
		fmt.Println(raw)

		return err
	}); err != nil {
		panic(err)
	}
}
