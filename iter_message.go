package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/query"
	"github.com/gotd/td/telegram/query/dialogs"
	"github.com/gotd/td/telegram/query/messages"
	"github.com/gotd/td/tg"
)

func iter_message() {
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
		//raw.Query(ctx, &query.Query{}
		cb := func(ctx context.Context, dlg dialogs.Elem) error {
			// Skip deleted dialogs.
			if dlg.Deleted() {
				return nil
			}
			return dlg.Messages(raw).ForEach(ctx, func(ctx context.Context, elem messages.Elem) error {
				msg, ok := elem.Msg.(*tg.Message)
				if !ok {
					return nil
				}
				peer := msg
				fmt.Println(peer.Message)
				time.Sleep(200 * time.Millisecond)
				return nil
			})
		}

		return query.GetDialogs(raw).ForEach(ctx, cb)
	}); err != nil {
		panic(err)
	}

	fmt.Println(data.DC, data.Addr)
}
