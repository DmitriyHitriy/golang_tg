package donors

import (
	"bufio"
	"context"
	"os"

	"github.com/gotd/td/tg"

	"golang_tg/internal/accounts"
)

type Donor struct {
	Account     *account.Account
	Users []*tg.User
}

func (d *Donor) DonorSetUsers(users []*tg.User) {
	d.Users = users
}

func (d *Donor) DonorGetUsers() []*tg.User {
	var users []*tg.User

	channel_list, err := d.donorGetChannelList()

	if err != nil {
		return users
	}

	var channels []*tg.Channel

	for _, channel_string := range channel_list {
		d.Account.Connect()
		tg_channels_result := d.donorSearchChannelFromQueryString(channel_string, 10)

		for _, tg_channel := range tg_channels_result {
			participants_count, _ := tg_channel.GetParticipantsCount()

			if tg_channel.Megagroup && participants_count > 100 {
				channels = append(channels, tg_channel)
			}
		} 
	}

	for _, channel := range channels {
		d.Account.Connect()
		users_tmp := d.donorSearchUsersFromMessagesChannel(channel)
		for _, user := range users_tmp {
			if !d.DonorIsUniqUser(users, user.ID) {
				users = append(users, user)
			}
		}
	}

	d.DonorSetUsers(users)

	return users

}

func (d *Donor) donorSearchUsersFromMessagesChannel(tg_channel *tg.Channel) []*tg.User {
	var users []*tg.User
	var err error

	if err := d.Account.GetClient().Run(*d.Account.GetContext(), func(ctx context.Context) error {
		raw := tg.NewClient(d.Account.GetClient())

		ch := &tg.InputPeerChannel{ChannelID: tg_channel.AsInput().ChannelID, AccessHash: tg_channel.AsInput().AccessHash}
		var offset int

		for {
			if len(users) <= 50 {
				req := tg.MessagesGetHistoryRequest{
					Peer:      ch,
					Limit:     100,
					AddOffset: offset,
				}

				res, err := raw.MessagesGetHistory(ctx, &req)
				if err != nil {
					break
				}
				res_user_list := (*res.(*tg.MessagesChannelMessages)).Users

				for _, user := range res_user_list {
					if user.((*tg.User)).Bot {
						continue
					}
					users = append(users,  user.((*tg.User)))
					
					// if elementExists(users, user.((*tg.User)).ID) == false {
					// 	users = append(users, user.((*tg.User)))
					// }
				}

				if len((*res.(*tg.MessagesChannelMessages)).Messages) == 100 {
					offset += 100
				} else {
					break
				}

			} else {
				break
			}
		}

		return err
	}); err != nil {
		panic(err)
	}

	return users

}

func (d *Donor) donorSearchChannelFromQueryString(query string, limit int) []*tg.Channel {
	var chats_results []*tg.Channel

	if err := d.Account.GetClient().Run(*d.Account.GetContext(), func(ctx context.Context) error {
		var err error
		raw := tg.NewClient(d.Account.GetClient())

		req_contact_search := tg.ContactsSearchRequest{
			Q:     query,
			Limit: limit,
		}

		res_contact_search, _ := raw.ContactsSearch(ctx, &req_contact_search)

		for _, chat := range (*res_contact_search).Chats {
			chats_results = append(chats_results, chat.(*tg.Channel))
		}

		return err
	}); err != nil {
		panic(err)
	}

	return chats_results

}

func (d *Donor) donorGetChannelList() ([]string, error) {
	file, err := os.Open("input/channel_list")
	if err != nil {
		return nil, err
	}

	buf := bufio.NewScanner(file)
	rows := make([]string, 0)

	for buf.Scan() {
		rows = append(rows, buf.Text())
	}

	return rows, nil

}

func (d *Donor) DonorIsUniqUser(haystack []*tg.User, user_id int64) bool {
	for _, v := range haystack {
		if v.ID == user_id {
			return true
		}
	}
	return false

}