package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func save_post(query string, limit int) []*tg.Channel {
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

	var chats_results []*tg.Channel

	client := telegram.NewClient(1, "s", telegram.Options{SessionStorage: storage})

	if err := client.Run(ctx, func(ctx context.Context) error {
		raw := tg.NewClient(client)

		req := tg.ContactsSearchRequest{
			Q:     query,
			Limit: limit,
		}

		res, _ := raw.ContactsSearch(ctx, &req)

		for _, chat := range (*res).Chats {
			chats_results = append(chats_results, chat.(*tg.Channel))
			//fmt.Println(chat)
		}

		

		return err
	}); err != nil {
		panic(err)
	}

	return chats_results
}
