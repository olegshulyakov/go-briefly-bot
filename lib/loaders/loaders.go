package loaders

import (
	"errors"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/youtube"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/youtube_short"
)

func VideoLoader(url string) (*video.DataLoader, error) {
	if youtube.IsValidURL(url) {
		return youtube.New(url)
	} else if youtube_short.IsValidURL(url) {
		return youtube_short.New(url)
	}
	return nil, errors.New("no valid URL found")
}
