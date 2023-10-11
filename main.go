package main

import (
	"fmt"
	//"os"
	"path/filepath"
	"time"
	"math/rand"

	accs "golang_tg/internal/accounts"
	cfg "golang_tg/internal/configs"
	donors "golang_tg/internal/donors"
)

func main() {
	// Читаем конфигурационный файл
	// Без него паникуем
	cfg := cfg.Configs{}
	cfg.New()


	// Собираем и чекаем аккаунты из папки accounts
	// внутри папки с аккаунтами должны лежать папки внутри которых лежат tdata
	// accounts - folder_x - tdata

	accounts_dir, _ := filepath.Glob("tdatas/*")
	var work_accounts accs.Accounts

	for _, acc_dir := range accounts_dir {
		tdata_folder_path := filepath.Join(acc_dir, "tdata")

		account := accs.Account{}
		account.Constructor(tdata_folder_path)

		// Проверяем, живой ли аккаунт
		if account.CheckAcc() {
			work_accounts.AddAccount(&account)

		}
		// Проверяем есть ли у аккаунта созданный рекламмный канал
		if !account.Channel.CheckChannel(account.GetTDataPath()) {
			fmt.Println("Канал не обнаружен. Создаем.")
			account.Connect()
			account.Channel.Createchannel(account.GetClient(), cfg.GetChannelName(), cfg.GetChannelDesc(), cfg.GetChannelPhoto(), account.GetTDataPath())
		}
	}

	// Собираем доноров со списка каналов input/channel_list
	for _, account := range work_accounts.Accounts {
		account.Connect()
		account.Channel.GetChannelInfo(*account.GetContext(), account.GetClient(),account.Channel.GetChannel())
		donor := donors.Donor{Account: account}
		donor.DonorGetUsers()
		donor.DonorGetPosts()
		p := donor.Account.GetPostNext()
		account.Connect()
		account.Channel.CreatePost(*account.GetContext(), account.GetClient(), account.Channel.GetUserName(), p)
		fmt.Println(p)
	}

	// Инвайтим юзеров или пишем пост с оффером в группу
	for {
		mode := rand.Intn(20-1) + 1

		for _, account := range work_accounts.Accounts {
			switch {
			case mode == 1:
				account.Connect()
				account.Channel.ChannelSendMessage(*account.GetContext(), account.GetClient(), account.Channel.GetUserName(), cfg.GetOfferText(), cfg.GetOfferPhoto())
				fmt.Println("Разместили рекламный оффер")
			case mode > 1 && mode <= 8:
				account.Connect()
				post := account.GetPostNext()
				account.Channel.CreatePost(*account.GetContext(), account.GetClient(), account.Channel.GetUserName(), post)
				fmt.Println("Разместили пост")
			default:
				us := account.GetUserNext()
				account.Connect()
				account.Channel.InviteToChannel(*account.GetContext(), account.GetClient(), account.Channel.GetChannel(), us)
				fmt.Println("Добавили человека в группу")
			}

			time.Sleep(5 * time.Second)
		}

		time.Sleep(300 * time.Second)
	}
	fmt.Println("Закончили работу")
}
