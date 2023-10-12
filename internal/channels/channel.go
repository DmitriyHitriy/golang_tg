package channel

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gookit/ini"
	"github.com/goombaio/namegenerator"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/html"
	"github.com/gotd/td/tg"
)

type Channel struct {
	ID int64
	AccessHash int64
	UserName string
	Title string
	ParticipantsCount int
	InputChannel *tg.InputChannel
}

func (c *Channel) GetChannel() *tg.InputChannel {
	return c.InputChannel
}

func (c *Channel) GetID() int64 {
	return c.ID
}

func (c *Channel) GetAccessHash() int64 {
	return c.AccessHash
}

func (c *Channel) GetUserName() string {
	return c.UserName
}

func (c *Channel) GetTitle() string {
	return c.Title
}

func (c *Channel) GetParticipantsCount() int {
	return c.ParticipantsCount
}

func (c *Channel) SetChannel(ch *tg.InputChannel) {
	c.InputChannel = ch
}

func (c *Channel) SetID(id int64) {
	c.ID = id
}

func (c *Channel) SetAccessHash(access_hash int64) {
	c.AccessHash = access_hash
}

func (c *Channel) SetUserName(user_name string) {
	c.UserName = user_name
}

func (c *Channel) SetTitle(title string) {
	c.Title = title
}

func (c *Channel) SetParticipantsCount(participants_count int) {
	c.ParticipantsCount = participants_count
}

func (c *Channel) InviteToChannel(ctx context.Context, client *telegram.Client, channel *tg.InputChannel, user *tg.User) (bool, error) {
	err := client.Run(ctx, func(ctx context.Context) error {
		raw := tg.NewClient(client)

		us := &tg.InputUser{UserID: user.GetID(), AccessHash: user.AccessHash}

		req := tg.ChannelsInviteToChannelRequest{
			Channel: channel,
			Users:   []tg.InputUserClass{us},
		}

		_, e := raw.ChannelsInviteToChannel(ctx, &req)

		return e
	}) 

	if err == nil {
		return true, nil
	} else {
		return false, err
	}

}

func (c *Channel) CreatePost(ctx context.Context, client *telegram.Client, channel string, post *tg.Message) {
	if err := client.Run(ctx, func(ctx context.Context) error {
		var err error
		raw := tg.NewClient(client)

		if post.Media.TypeName() == "messageMediaPhoto" {
			media := post.Media.(*tg.MessageMediaPhoto)
			photo, _ := media.GetPhoto()
			originPhoto := photo.(*tg.Photo)
			photo_file_location := tg.InputPhotoFileLocation{
				ID: originPhoto.GetID(),
				AccessHash: originPhoto.GetAccessHash(),
				FileReference: originPhoto.GetFileReference(),
				ThumbSize: "500",
			}

		message.NewSender(raw).Resolve(channel).Photo(ctx, &photo_file_location, html.String(nil, post.Message))
		}
		

		return err
	}); err != nil {
		panic(err)
	}
}

func (c *Channel) GetChannelInfo(ctx context.Context, client *telegram.Client, channel *tg.InputChannel){
	if err := client.Run(ctx, func(ctx context.Context) error {
		var err error
		raw := tg.NewClient(client)

		res, _ := raw.ChannelsGetChannels(ctx, []tg.InputChannelClass{channel})
		channel_data := (*(*res.(*tg.MessagesChats)).Chats[0].(*tg.Channel))

		c.SetID(channel_data.ID)
		c.SetAccessHash(channel_data.AccessHash)
		c.SetTitle(channel_data.Title)
		c.SetUserName(channel_data.Username)
		c.SetParticipantsCount(channel_data.ParticipantsCount)

		return err
	}); err != nil {
		panic(err)
	}
}

func (c *Channel) ChannelSendMessage(ctx context.Context, client *telegram.Client, channel string, text string, photo string) {
	if err := client.Run(ctx, func(ctx context.Context) error {
		var err error
		raw := tg.NewClient(client)
		
		nm, _ := message.NewSender(raw).Resolve(channel).Upload(message.Upload(func(ctx context.Context, b message.Uploader) (tg.InputFileClass, error) {
			r, err := b.FromPath(ctx, photo)
			if err != nil {
				return nil, err
			}

			return r, nil
		})).Photo(ctx, html.String(nil, text))

		fmt.Println(nm)
		fmt.Println(err)

		return err
	}); err != nil {
		panic(err)
	}
	
}

func (a *Channel) CheckChannel(tdata_path string) bool {
	cfg_path := filepath.Join(tdata_path, "channel.ini")
	cfg_channel, err := ini.LoadFiles(cfg_path)
	if err != nil {
		return false
	}

	channel_id, _ := cfg_channel.Int("channel_id")
	channel_accesshash, _ := cfg_channel.Int("channel_access_hash")

	input_peer := &tg.InputChannel{ChannelID: int64(channel_id), AccessHash: int64(channel_accesshash)}
	a.SetChannel(input_peer)

	return true
}

func (a *Channel) Createchannel(client *telegram.Client, name string, about string, photo_path string, tdata_path string) bool {
	ctx := context.Background()
	if err := client.Run(ctx, func(ctx context.Context) error {
		raw := tg.NewClient(client)

		// Создаем канал
		req_create_channel := tg.ChannelsCreateChannelRequest{
			Title:     name,
			About:     about,
			Broadcast: true,
			Megagroup: false,
		}

		res_create_channel, err := raw.ChannelsCreateChannel(ctx, &req_create_channel)
		channel := ((*res_create_channel.(*tg.Updates)).Chats[0]).(*tg.Channel)

		ch_input := tg.InputChannel{ChannelID: channel.ID, AccessHash: channel.AccessHash}

		a.SetChannel(&ch_input)

		// Генерируем имя канала и назначаем его каналу
		var channel_username string
		for {
			tmp_channel_username := a.generateUsername()
			req_update_username := tg.ChannelsUpdateUsernameRequest{
				Channel:  &ch_input,
				Username: tmp_channel_username,
			}

			res_update_username, _ := raw.ChannelsUpdateUsername(ctx, &req_update_username)

			if res_update_username {
				channel_username = tmp_channel_username
				break
			}
		}

		// Устанавливаем аватарку каналу
		nm, _ := message.NewSender(raw).Resolve(channel_username).Upload(message.Upload(func(ctx context.Context, b message.Uploader) (tg.InputFileClass, error) {
			r, err := b.FromPath(ctx, photo_path)
			if err != nil {
				return nil, err
			}

			return r, nil
		})).Photo(ctx)

		photo := (*(*(*(*(*nm.(*tg.Updates)).Updates[2].(*tg.UpdateNewChannelMessage)).Message.(*tg.Message)).Media.(*tg.MessageMediaPhoto)).Photo.(*tg.Photo))

		// Ставим загруженное фото на аватарку
		in_photo := tg.InputPhoto{ID: photo.ID, AccessHash: photo.AccessHash, FileReference: photo.FileReference}
		chat_photo := tg.InputChatPhoto{ID: &in_photo}

		req_edit_photo := tg.ChannelsEditPhotoRequest{
			Channel: &ch_input,
			Photo:   &chat_photo,
		}

		raw.ChannelsEditPhoto(ctx, &req_edit_photo)

		a.SetID(channel.ID)
		a.SetAccessHash(channel.AccessHash)
		a.SetUserName(channel_username)
		a.SetTitle(name)

		a.createConfigChannel(tdata_path)
		return err
	}); err != nil {
		panic(err)
	}

	return true

}

func (a *Channel) generateUsername() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := strings.ReplaceAll(nameGenerator.Generate()+"_bratkov", "-", "_")

	return name
}

func (a *Channel) createConfigChannel(tdata_path string) {
	cfg_path := filepath.Join(tdata_path, "channel.ini")
	cfg_channel, _ := ini.LoadExists(cfg_path)
	input_channel := a.GetChannel()

	cfg_channel.SetInt("channel_id", int(input_channel.GetChannelID()))
	cfg_channel.SetInt("channel_access_hash", int(input_channel.GetAccessHash()))

	cfg_channel.WriteToFile(cfg_path)
}