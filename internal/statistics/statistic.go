package statistics

import (
	"context"
	"strings"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/telegram"

	channel "golang_tg/internal/channels"
)

type Stats struct {
	participants_count int
	posts_count int
	posts_views int
	offer_count int
	offer_views int
}

func (s *Stats) GetParticipantsCount() int { return s.participants_count }
func (s *Stats) GetPostsCount() int { return s.posts_count }
func (s *Stats) GetPostsViewsCount() int { return s.posts_views }
func (s *Stats) GetOfferCount() int { return s.offer_count }
func (s *Stats) GetOfferViewsCount() int { return s.offer_views }

func (s *Stats) SetParticipantsCount(count int) { s.participants_count = count }
func (s *Stats) setPostsCount(count int) { s.posts_count = count }
func (s *Stats) setPostViews(count int) { s.posts_views = count }
func (s *Stats) setOfferCount(count int) { s.offer_count = count }
func (s *Stats) setOfferViewsCount(count int) { s.offer_count = count }

func (s *Stats) GetStats(ctx context.Context, client *telegram.Client, tg_channel channel.Channel) {
	var err error

	client.Run(ctx, func(ctx context.Context) error {
		raw := tg.NewClient(client)

		ch := &tg.InputPeerChannel{ChannelID: tg_channel.GetID(), AccessHash: tg_channel.GetAccessHash()}
		var cnt_offer_views_count int
		var cnt_offer_count int
		var cnt_post_views_count int
		var cnt_posts_count int
		
		var offset int

		for {
			req_message_history := tg.MessagesGetHistoryRequest{
				Peer:      ch,
				Limit:     100,
				AddOffset: offset,
			}

			res_message_history, err := raw.MessagesGetHistory(ctx, &req_message_history)
			if err != nil {
				break
			}
			messages := (*res_message_history.(*tg.MessagesChannelMessages)).Messages

			for _, message := range messages {
			
				if message.TypeName() == "message" {
					if strings.Contains(message.((*tg.Message)).Message, "#myoffer") {
						cnt_offer_views_count += message.((*tg.Message)).Views
						s.setOfferViewsCount(cnt_offer_views_count)

						cnt_offer_count += 1
						s.setOfferCount(cnt_offer_count)
					}

					cnt_post_views_count += message.((*tg.Message)).Views
					s.setPostViews(cnt_post_views_count)

					cnt_posts_count += 1
					s.setPostsCount(cnt_posts_count)
				}
			}

			if len((*res_message_history.(*tg.MessagesChannelMessages)).Messages) == 100 {
				offset += 100
			} else {
				break
			}
		}

		return err
	})

}
