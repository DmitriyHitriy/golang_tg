package main

import (
	"fmt"
	"time"

	functions "golang_tg/cmd"

	"github.com/gotd/td/tg"
)

func main() {
	//TGToolsGenerateChannel("Histórias incríveis de vitórias 💎", "🔥 Histórias coletadas do Brasil sobre vitórias incríveis de pessoas. Tente repetir suas histórias de sucesso. 💲", "casino.jpg")

	//os.Exit(1)
	var channels_donor []*tg.Channel
	var users_donor []*tg.User

	channel_list, _ := functions.Get_rows_in_file("input/channel_list")

	for _, tg_channel := range channel_list {
		channels_result := search_contact(tg_channel, 10)

		for _, channel := range channels_result {
			participants_count, _ := channel.GetParticipantsCount()

			if (channel.Megagroup == true) && (participants_count > 100) {
				fmt.Println(channel.Title, channel.ParticipantsCount)
				channels_donor = append(channels_donor, channel)
			}
		}

	}

	for _, tg_channel := range channels_donor {
		users_in_channel := get_users_in_messages_from_channel(tg_channel.ID, tg_channel.AccessHash)
		for _, user := range users_in_channel {
			if elementExists(users_donor, user.ID) == false {
				users_donor = append(users_donor, user)
			}
		}
	}

	for i, user := range users_donor {
		if i <= 8 {
			ch_input := &tg.InputChannel{ChannelID: 1941890406, AccessHash: 5678437962803547359}
			us_unput := &tg.InputUser{UserID: user.ID, AccessHash: user.AccessHash}
			res_invite := invite_to_channel(ch_input, us_unput)

			if res_invite == true {
				time.Sleep(10 * time.Second)
			}
		} else {
			fmt.Println("Исчерпан лимит инвайта")
			break
		}
	}

}
