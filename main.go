package main

import (
	"log"
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
		account.SetConfig(cfg)

		// Проверяем, живой ли аккаунт
		if account.CheckAcc() {
			work_accounts.AddAccount(&account)
		}
		
		// Проверяем есть ли у аккаунта созданный рекламмный канал
		if !account.Channel.CheckChannel(account.GetTDataPath()) {
			log.Print(account.GetFullName() + "создает канал")
			account.Connect()
			account.Channel.Createchannel(account.GetClient(), cfg.GetChannelName(), cfg.GetChannelDesc(), cfg.GetChannelPhoto(), account.GetTDataPath())
		}
	}

	// Собираем доноров со списка каналов input/channel_list
	for _, account := range work_accounts.Accounts {
		account.Connect()
		account.Channel.GetChannelInfo(*account.GetContext(), account.GetClient(),account.Channel.GetChannel())
		donor := donors.Donor{Account: account}
		log.Print(account.GetFullName() + "собирает аудиторию")
		donor.DonorGetUsers()
		log.Print(account.GetFullName() + "собирает контент")
		donor.DonorGetPosts()
	}

	// Инвайтим юзеров или пишем пост с оффером в группу
	for {
		mode := rand.Intn(20-1) + 1

		for _, account := range work_accounts.Accounts {
			switch {
			case mode == 1:
				account.Connect()
				account.Channel.ChannelSendMessage(*account.GetContext(), account.GetClient(), account.Channel.GetUserName(), cfg.GetOfferText(), cfg.GetOfferPhoto())
			case mode > 1 && mode <= 9:
				account.Connect()
				post := account.GetPostNext()
				if post != nil{
					account.Channel.CreatePost(*account.GetContext(), account.GetClient(), account.Channel.GetUserName(), post)
				}
			default:
				us := account.GetUserNext()
				if us != nil {
					account.Connect()
					_, err := account.Channel.InviteToChannel(*account.GetContext(), account.GetClient(), account.Channel.GetChannel(), us)
					log.Print(account.GetFullName() + " " + err.Error())
				}
			}
			log.Print(account.GetFullName() + account.PrintStats())

			time.Sleep(5 * time.Second)
		}

		log.Print("Спим 5 минут...\n")
		time.Sleep(600 * time.Second)
	}
}
