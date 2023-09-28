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

func search_messages(query string, limit int) []*tg.User {
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
	var users []*tg.User

	if err := client.Run(ctx, func(ctx context.Context) error {
		raw := tg.NewClient(client)

		req := tg.MessagesSearchGlobalRequest{
			Q:          query,
			Limit:      limit,
			Filter:     &tg.InputMessagesFilterEmpty{},
			OffsetPeer: &tg.InputPeerEmpty{},
		}

		res, e := raw.MessagesSearchGlobal(ctx, &req)

		for _, user := range (*res.(*tg.MessagesMessages)).Users {
			users = append(users, user.(*tg.User))
		}

		fmt.Println(res)
		fmt.Println(e)

		return err
	}); err != nil {
		panic(err)
	}

	return users
}
