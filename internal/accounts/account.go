package account

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gookit/ini"
	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

type Account struct {
	first_name string
	last_name  string
	username   string
	phone      string
	tdata_path string
	last_use   time.Time
	client     *telegram.Client
	channel    *tg.InputPeerChannel
	ctx        context.Context
}

func (a *Account) Constructor(path string) {
	ctx := context.Background()
	a.ctx = ctx

	a.SetTDataPath(path)

	accounts, err := tdesktop.Read(a.GetTDataPath(), nil)
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

	a.client = telegram.NewClient(1, "s", telegram.Options{SessionStorage: storage})

}

func (a *Account) CheckAcc() bool {
	if err := a.client.Run(a.ctx, func(ctx context.Context) error {
		me, err := a.client.Self(ctx)
		fmt.Println("Успешно авторизовались: ", me.FirstName, me.LastName)
		a.SetLastUse()
		return err
	}); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (a *Account) GetClient() *telegram.Client {
	return a.client
}

func (a *Account) GetChannel() *tg.InputPeerChannel {
	return a.channel
}

func (a *Account) GetFirstName() string {
	return a.first_name
}

func (a *Account) GetLastName() string {
	return a.last_name
}

func (a *Account) GetUsername() string {
	return a.username
}

func (a *Account) GetPhone() string {
	return a.phone
}

func (a *Account) GetTDataPath() string {
	return a.tdata_path
}

func (a *Account) GetLastUse() time.Time {
	return a.last_use
}

func (a *Account) SetClient(client *telegram.Client) {
	a.client = client
}

func (a *Account) SetChannel(channel *tg.InputPeerChannel) {
	a.channel = channel
}

func (a *Account) SetFirstName(first_name string) {
	a.first_name = first_name
}

func (a *Account) SetLastName(last_name string) {
	a.last_name = last_name
}

func (a *Account) SetUsername(username string) {
	a.username = username
}

func (a *Account) SetPhone(phone string) {
	a.phone = phone
}

func (a *Account) SetTDataPath(tdata_path string) {
	a.tdata_path = tdata_path
}

func (a *Account) SetLastUse() {
	a.last_use = time.Now()
}

func NewAccount(first_name string, last_name string, username string, phone string, tdata_path string) *Account {
	return &Account{
		first_name: first_name,
		last_name:  last_name,
		username:   username,
		phone:      phone,
		tdata_path: tdata_path,
		last_use:   time.Now(),
	}
}

func (a *Account) CheckChannel() bool {
	cfg_path := filepath.Join(a.GetTDataPath(), "channel.ini")
	cfg_channel, err := ini.LoadFiles(cfg_path)
	if err != nil {
		return false
	}

	channel_id, _ := cfg_channel.Int("channel_id")
	channel_accesshash, _ := cfg_channel.Int("channel_access_hash")

	input_peer := &tg.InputPeerChannel{ChannelID: int64(channel_id), AccessHash: int64(channel_accesshash)}
	a.SetChannel(input_peer)

	return true
}
