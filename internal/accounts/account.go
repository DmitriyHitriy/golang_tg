package account

import (
	"bufio"
	"context"
	"log"
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
	config "golang_tg/internal/configs"
	stats "golang_tg/internal/statistics"
)

type Account struct {
	first_name string
	last_name  string
	username   string
	phone      string
	tdata_path string
	last_use   time.Time
	next_use   time.Time
	client     *telegram.Client
	users      []*tg.User
	posts      []*tg.Message
	ctx        context.Context
	config     config.Configs
	Channel    channel.Channel
	Stats      stats.Stats
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
	err := a.client.Run(a.ctx, func(ctx context.Context) error {
		me, err := a.client.Self(ctx)

		a.SetFirstName(me.FirstName)
		a.SetLastName(me.LastName)
		a.SetUsername(me.Username)
		a.SetPhone(me.Phone)
		a.SetLastUse()

		log.Println(a.GetFullName() + "успешно авторизовались")
		return err
	})

	if err != nil {
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

func (a *Account) GetNextUse() time.Time {
	return a.next_use
}

func (a *Account) GetContext() *context.Context {
	return &a.ctx
}

func (a *Account) GetUsers() []*tg.User {
	return a.users
}

func (a *Account) GetPosts() []*tg.Message {
	return a.posts
}

func (a *Account) GetConfig() config.Configs {
	return a.config
}

func (a *Account) GetFullName() string {
	fullname_and_stat := a.GetFirstName() + " " + a.GetLastName() + " (очередь инвайтинга: " + strconv.Itoa(len(a.GetUsers())) + ", очередь постинга: " + strconv.Itoa(len(a.GetPosts())) + ") "
	return fullname_and_stat
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

		global_invited_users := a.readGlobalUsesDonor()

		if !a.isGlobalUsesDonor(global_invited_users, int(user.ID)) {
			a.addGlobalUsesDonor(int(user.ID))
			return user
		}
	}
}

func (a *Account) GetPostNext() *tg.Message {
	for {
		if len(a.GetPosts()) == 0 {
			return nil
		}

		index := rand.Intn(len(a.GetPosts()))
		post := a.GetPosts()[index]

		new_post_list := a.GetPosts()

		new_post_list = slices.Delete(new_post_list, index, index+1)

		a.SetPosts(new_post_list)

		// global_used_posts := a.readGlobalUsesDonor()
		// id_post, _ := strconv.Atoi(strconv.Itoa(post.Date) + strconv.Itoa(post.ID))
		// if !a.isGlobalUsesDonor(global_used_posts, id_post) {
		// 	a.addGlobalUsesDonor(id_post)
		// 	return post
		// }

		return post
	}

}

func (a *Account) PrintStats() string {
	a.Connect()
	a.Channel.GetPaticipantsCountFromChannel(*a.GetContext(), a.GetClient(), a.Channel.GetChannel())

	a.Connect()
	a.Stats.GetStats(*a.GetContext(), a.GetClient(), a.Channel)

	participants_count := a.Channel.GetParticipantsCount()
	offer_count := a.Stats.GetOfferCount()
	offer_views_count := a.Stats.GetOfferViewsCount()
	posts_count := a.Stats.GetPostsCount()
	posts_views_count := a.Stats.GetPostsViewsCount()

	return "Подписчиков: " + strconv.Itoa(participants_count) + " Офферов: " + strconv.Itoa(offer_count) + " Просмотров офферов: " + strconv.Itoa(offer_views_count) + " Постов: " + strconv.Itoa(posts_count) + " Просмотров постов: " + strconv.Itoa(posts_views_count) + " "
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

func (a *Account) SetNextUse(duration int) {
	a.next_use = time.Now().Add(time.Minute * 30)
}

func (a *Account) SetUsers(users []*tg.User) {
	a.users = users
}

func (a *Account) SetPosts(posts []*tg.Message) {
	a.posts = posts
}

func (a *Account) SetConfig(config config.Configs) {
	a.config = config
}

func (a *Account) readGlobalUsesDonor() []string {
	ids := make([]string, 0)

	f, e := os.Open("global_uses_donor")
	if e != nil {
		a.writeGlobalUsesDonor()
		return ids
	}

	defer f.Close()

	buf := bufio.NewScanner(f)

	for buf.Scan() {
		ids = append(ids, buf.Text())
	}

	return ids
}

func (a *Account) writeGlobalUsesDonor() {
	var data []byte
	os.WriteFile("global_uses_donor", data, 0700)
}

func (a *Account) addGlobalUsesDonor(id int) {
	f, _ := os.OpenFile("global_uses_donor", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	f.WriteString(strconv.Itoa(id) + "\n")
}

func (a *Account) isGlobalUsesDonor(haystack []string, id int) bool {
	for _, v := range haystack {
		if v == strconv.Itoa(id) {
			return true
		}
	}
	return false
}

func (a *Account) IsPossibleToUse() bool {
	current_date := time.Now()
	next_use := a.GetNextUse()

	if current_date.Unix() >= next_use.Unix() {
		return true
	} else {
		return false
	}
}