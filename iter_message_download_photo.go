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
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/telegram/query"
	"github.com/gotd/td/telegram/query/dialogs"
	"github.com/gotd/td/telegram/query/messages"
	
	"github.com/gotd/td/tg"
)

func downloader_photo() {
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
				// if msg.Media != nil {
				// 	fmt.Println(msg.Media.TypeName())
				// }
				if msg.Media != nil && msg.Media.TypeName() == "messageMediaPhoto" {
					media := msg.Media.(*tg.MessageMediaPhoto)
					photo, _ := media.GetPhoto()

					originPhoto := photo.(*tg.Photo)
					downloaderInstance := downloader.NewDownloader()
					file := &tg.InputPhotoFileLocation{
						ID:            originPhoto.GetID(),
						AccessHash:    originPhoto.GetAccessHash(),
						FileReference: originPhoto.GetFileReference(),
						ThumbSize: "500",
					}
					fmt.Println(file)
					_, err = downloaderInstance.Download(client.API(), file).ToPath(ctx, filepath.Join(fmt.Sprintf("%d.png", originPhoto.GetID())))
					if err!=nil {
						return  err
					}
				}
				time.Sleep(200 * time.Millisecond)
				return nil
			})
		}

		return query.GetDialogs(raw).ForEach(ctx, cb)
	}); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(data.DC, data.Addr)
}
