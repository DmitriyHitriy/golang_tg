package account

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"

	"golang.org/x/exp/slices"

	channel "golang_tg/internal/channels"
)

type Account struct {
	first_name string
	last_name  string
	username   string
	phone      string
	tdata_path string
	last_use   time.Time
	client     *telegram.Client
	users      []*tg.User
	ctx        context.Context
	Channel    channel.Channel
}

func (a *Account) Connect() {
	a.Constructor(a.GetTDataPath())
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

		a.SetFirstName(me.FirstName)
		a.SetLastName(me.LastName)
		a.SetUsername(me.Username)
		a.SetPhone(me.Phone)
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

func (a *Account) GetContext() *context.Context {
	return &a.ctx
}

func (a *Account) GetUsers() []*tg.User {
	return a.users
}

func (a *Account) GetUserNext() *tg.User {
	for {
		if len(a.GetUsers()) == 0 {
			return nil
		}

		index := rand.Intn(len(a.GetUsers()))
		user := a.GetUsers()[index]

		new_user_list := a.GetUsers()

		new_user_list = slices.Delete(new_user_list, index, index+1)

		a.SetUsers(new_user_list)

		global_invited_users := a.readGlobalInvitedUsers()
		if !a.isGlobalInvitedUsers(global_invited_users, int(user.ID)) {
			a.addGlobalInvitedUsers(int(user.ID))
			return user
		}
	}
}

func (a *Account) SetClient(client *telegram.Client) {
	a.client = client
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

func (a *Account) SetUsers(users []*tg.User) {
	a.users = users
}

func (a *Account) readGlobalInvitedUsers() []string {
	id_users := make([]string, 0)

	f, e := os.Open("global_invited_users")
	if e != nil {
		a.writeGlobalInvitedUsers()
		return id_users
	}

	defer f.Close()

	buf := bufio.NewScanner(f)

	for buf.Scan() {
		id_users = append(id_users, buf.Text())
	}

	return id_users
}

func (a *Account) writeGlobalInvitedUsers() {
	var data []byte
	os.WriteFile("global_invited_users", data, 0700)
}

func (a *Account) addGlobalInvitedUsers(id int) {
	f, _ := os.OpenFile("global_invited_users", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	f.WriteString(strconv.Itoa(id) + "\n")
}

func (a *Account) isGlobalInvitedUsers(haystack []string, id int) bool {
	for _, v := range haystack {
		if v == strconv.Itoa(id) {
			return true
		}
	}
	return false
}
