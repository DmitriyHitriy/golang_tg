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

func invite_to_channel(ch *tg.InputChannel, us *tg.InputUser) bool {
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

	//ch := &tg.InputChannel{ChannelID: 1905046891, AccessHash: 5725182504979867499}
	//us := &tg.InputUser{UserID: 5419321810, AccessHash: -5698706159390738840}
	//user_info, _ := client.API().UsersGetUsers(ctx, []tg.InputUserClass{us})
	//fmt.Println(user_info)

	if err := client.Run(ctx, func(ctx context.Context) error {
		raw := tg.NewClient(client)

		//users := []tg.InputUser{*us}

		req := tg.ChannelsInviteToChannelRequest{
			Channel: ch,
			Users:   []tg.InputUserClass{us},
		}

		res, e := raw.ChannelsInviteToChannel(ctx, &req)

		if e != nil {
			fmt.Println("Ошибка добавления в канал ", e)
			//fmt.Println((*user_info[0].(*tg.User)).FirstName, (*user_info[0].(*tg.User)).LastName, " ошибка при приглашении в канал", e)
		}

		fmt.Println(res)
		fmt.Println(e)

		return err
	}); err != nil {
		panic(err)
	}
	fmt.Println("Успешно пригласил в канал")
	//fmt.Println((*user_info[0].(*tg.User)).FirstName, (*user_info[0].(*tg.User)).LastName, " успешно пригласил в канал")
	return true
}
