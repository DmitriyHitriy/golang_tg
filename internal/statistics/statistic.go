package statistics

import (
	"context"
	"github.com/gotd/td/tg"
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

func (s *Stats) setParticipantsCount(count int) { s.participants_count = count }
func (s *Stats) setPostsCount(count int) { s.posts_count = count }
func (s *Stats) setPostViews(count int) { s.posts_views = count }
func (s *Stats) setOfferCount(count int) { s.offer_count = count }
func (s *Stats) setOfferViewsCount(count int) { s.offer_count = count }

func (s *Stats) GetStats(ctx context.Context, client *tg.Client, ) {

}