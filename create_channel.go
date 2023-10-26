package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func TGToolsGenerateChannel(name string, about string, photo_path string) {
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

		// Создаем канал
		ch_input, err := create_channel(ctx, name, about, raw)

		// Для теста временно. Убрать потом
		//ch_input := tg.InputChannel{ChannelID: 1905046891, AccessHash: 5725182504979867499}

		// Генерируем имя канала и назначаем его каналу 223
		channel_username := change_name_channel(&ctx, ch_input, raw)

		// Устанавливаем аватарку каналу
		change_photo_channel(&ctx, ch_input, raw, channel_username, photo_path)

		return err
	}); err != nil {
		panic(err)
	}
}

func create_channel(ctx context.Context, name string, about string, raw *tg.Client) (*tg.InputChannel, error) {
	// Создаем канал
	req := tg.ChannelsCreateChannelRequest{
		Title:     name,
		About:     about,
		Broadcast: true,
		Megagroup: false,
	}

	res, err := raw.ChannelsCreateChannel(ctx, &req)
	channel := ((*res.(*tg.Updates)).Chats[0]).(*tg.Channel)

	ch_input := tg.InputChannel{
		ChannelID:  channel.ID,
		AccessHash: channel.AccessHash,
	}
	fmt.Println(ch_input)
	return &ch_input, err
}

func change_name_channel(ctx *context.Context, ch_input *tg.InputChannel, raw *tg.Client) string {
	var channel_username string
	for {
		tmp_channel_username := generate_name()
		n := tg.ChannelsUpdateUsernameRequest{
			Channel:  ch_input,
			Username: tmp_channel_username,
		}

		res_update_username, _ := raw.ChannelsUpdateUsername(*ctx, &n)

		if res_update_username {
			channel_username = tmp_channel_username
			break
		}
	}
	return channel_username
}

func change_photo_channel(ctx *context.Context, ch_input *tg.InputChannel, raw *tg.Client, channel_username string, photo_path string) error {
	nm, _ := message.NewSender(raw).Resolve(channel_username).Upload(message.Upload(func(ctx context.Context, b message.Uploader) (tg.InputFileClass, error) {
		r, err := b.FromPath(ctx, photo_path)
		if err != nil {
			return nil, err
		}

		return r, nil
	})).Photo(*ctx)

	photo := (*(*(*(*(*nm.(*tg.Updates)).Updates[2].(*tg.UpdateNewChannelMessage)).Message.(*tg.Message)).Media.(*tg.MessageMediaPhoto)).Photo.(*tg.Photo))

	// Ставим загруженное фото на аватарку
	in_photo := tg.InputPhoto{ID: photo.ID, AccessHash: photo.AccessHash, FileReference: photo.FileReference}
	chat_photo := tg.InputChatPhoto{ID: &in_photo}

	req_photo := tg.ChannelsEditPhotoRequest{
		Channel: ch_input,
		Photo:   &chat_photo,
	}

	_, err := raw.ChannelsEditPhoto(*ctx, &req_photo)

	return err
}
