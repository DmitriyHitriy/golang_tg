package main

import (
	"fmt"
	//"os"
	"path/filepath"
	"time"

	accs "golang_tg/internal/accounts"
	cfg "golang_tg/internal/configs"
	donors "golang_tg/internal/donors"
)

func main() {
	// Читаем конфигурационный файл
	// Без него паникуем
	cfg := cfg.Configs{}
	cfg.New()
	fmt.Println(cfg)

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
		if !account.CheckChannel() {
			fmt.Println("Канал не обнаружен. Создаем его.")
			account.Constructor(tdata_folder_path)
			account.Createchannel(cfg.GetChannelName(), cfg.GetChannelDesc(), cfg.GetChannelPhoto())
		}
	}

	// Собираем доноров со списка каналов input/channel_list
	for _, account := range work_accounts.Accounts {
		account.Connect()
		account.Input_channel.GetChannelInfo(*account.GetContext(), account.GetClient(), account.GetChannel())
		account.Connect()
		account.Input_channel.ChannelSendMessage(*account.GetContext(), account.GetClient(), "morning_dew_bratkov")
		donor := donors.Donor{Account: account}
		donor.DonorGetUsers()
	}

	// Инвайтим юзеров
	for {
		for _, account := range work_accounts.Accounts {
			us := account.GetUserNext()
			account.Connect()
			account.Input_channel.InviteToChannel(*account.GetContext(), account.GetClient(), account.GetChannel(), us)
			time.Sleep(5 * time.Second)
		}

		time.Sleep(60 * time.Second)
	}
}
