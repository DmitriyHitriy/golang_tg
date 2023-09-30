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

// Найти список пользователей которые когда либо писали в чате.
func get_users_in_messages_from_channel(channel_id int64, access_hash int64) []*tg.User {
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

		//ch := &tg.InputPeerChannel{ChannelID: 1055678063, AccessHash: 2971438347507658668}
		ch := &tg.InputPeerChannel{ChannelID: channel_id, AccessHash: access_hash}
		var offset int

		for {
			if len(users) <= 50 {
				req := tg.MessagesGetHistoryRequest{
					Peer:      ch,
					Limit:     100,
					AddOffset: offset,
				}

				res, e := raw.MessagesGetHistory(ctx, &req)
				res_user_list := (*res.(*tg.MessagesChannelMessages)).Users

				for _, user := range res_user_list {
					if user.((*tg.User)).Bot == true {
						continue
					}

					if elementExists(users, user.((*tg.User)).ID) == false {
						users = append(users, user.((*tg.User)))
					}
				}

				if len((*res.(*tg.MessagesChannelMessages)).Messages) == 100 {
					offset += 100
				} else {
					break
				}

				fmt.Println(res)
				fmt.Println(e)

			} else {
				break
			}
		}

		return err
	}); err != nil {
		panic(err)
	}

	return users
}

func elementExists(haystack []*tg.User, user_id int64) bool {
	for _, v := range haystack {
		if v.ID == user_id {
			return true
		}
	}
	return false
}
