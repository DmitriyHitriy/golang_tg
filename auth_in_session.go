package main

import (
	"context"
	"fmt"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func main4() {
	app_id := 29917381
	app_acceshash := "be04216d4206c2eaad612101f3e07013"

	ctx := context.Background()
	var (
		storage = new(session.FileStorage)
	)

	storage.Path = "aaa.session"
	storage.LoadSession(ctx)

	client := telegram.NewClient(app_id, app_acceshash, telegram.Options{SessionStorage: storage})

	client.Run(ctx, func(ctx context.Context) error {
		var err error
		
		fmt.Println(client.Self(ctx))
		raw := tg.NewClient(client)

		message.NewSender(raw).Resolve("@hitriydima").Text(ctx, "Ha ha ha in session file auth")

		return err
	})
	
	
}
