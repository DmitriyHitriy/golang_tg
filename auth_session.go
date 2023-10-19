package main

import (
	"context"
	"fmt"
	"bufio"
	"os"
	"strings"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/telegram/auth"
	//"github.com/gotd/td/tg"
)

func main3() {
	phone := "+1 940 597 3779"
	password := "Fvnh215fgrd"
	app_id := 29917381
	app_acceshash := "be04216d4206c2eaad612101f3e07013"

	ctx := context.Background()
	var (
		storage = new(session.StorageMemory)
	)

	client := telegram.NewClient(app_id, app_acceshash, telegram.Options{SessionStorage: storage})

	codeAsk := func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
		fmt.Print("code:")
		code, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return "", err
		}
		code = strings.ReplaceAll(code, "\n", "")
		return code, nil
	}

	client.Run(ctx, func(ctx context.Context) error {
		res := auth.NewFlow(
			auth.Constant(phone, password, auth.CodeAuthenticatorFunc(codeAsk)),
			auth.SendCodeOptions{},
		).Run(ctx, client.Auth())
		storage.WriteFile("aaa.session", 0644)
		fmt.Println(client.Self(ctx))
		return res
	})
	
	
}
