package loaders

import (
	"errors"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/youtube"
)

func VideoLoader(url string) (*video.DataLoader, error) {
	if youtube.IsValidURL(url) {
		return youtube.New(url)
	}
	return nil, errors.New("no valid URL found")
}
