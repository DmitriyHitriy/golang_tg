package parser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	//"github.com/gotd/td/telegram/downloader"
)

func get_tg_channel_from_string(channel_name string) *tg.Channel {
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
	var searched_channel *tg.Channel

	if err := client.Run(ctx, func(ctx context.Context) error {
		raw := tg.NewClient(client)
		req := tg.ContactsSearchRequest{
			Q:     channel_name,
			Limit: 5,
		}

		res, _ := raw.ContactsSearch(ctx, &req)

		chats := (*res).Chats

		for _, chat := range chats {
			if chat.(*tg.Channel).Username == channel_name {
				searched_channel = chat.(*tg.Channel)
			}
		}

		return err
	}); err != nil {
		panic(err)
	}

	return searched_channel
}

func Channel_parser_post(channel_name string, limit int) {
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

	tg_channel := get_tg_channel_from_string(channel_name)
	input_channel := &tg.InputPeerChannel{ChannelID: tg_channel.ID, AccessHash: tg_channel.AccessHash}

	var offset int

	if err := client.Run(ctx, func(ctx context.Context) error {
		raw := tg.NewClient(client)
		req := tg.MessagesGetHistoryRequest{
			Peer:      input_channel,
			Limit:     100,
			AddOffset: offset,
		}

		res, _ := raw.MessagesGetHistory(ctx, &req)
		res_post_list := (*res.(*tg.MessagesChannelMessages)).Messages
		//fmt.Println(res_post_list)
		for _, post := range res_post_list {
			if post.TypeName() == "message" {
				post_data := post.(*tg.Message)
				fmt.Println(post_data.Message)
				if post_data.GroupedID == 0 {
					//d := downloader.NewDownloader()
					// media := post_data.Media
					// d.Download()
					// raw.Dow
					if post_data.Media.TypeName() == "messageMediaPhoto" {
						photo := post_data.Media.(*tg.MessageMediaPhoto).Photo.(*tg.Photo)
						// photo.AsInput()
						// photo_class := *&tg.InputPhoto{ID: photo.ID, AccessHash: photo.AccessHash}
						// tg.InputFileLocation
						// d.Download(raw, photo_class)
						fmt.Println(photo)
					}
					fmt.Println(post_data.Media.TypeName())
				}
			}
		}

		return err
	}); err != nil {
		panic(err)
	}
}
