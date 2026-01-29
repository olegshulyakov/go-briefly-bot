package loaders

import (
	"errors"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/provider"
)

func VideoLoader(url string) (*video.DataLoader, error) {
	switch {
	case provider.Youtube.IsValidURL(url):
		return provider.Youtube.BuildDataLoader(url)
	case provider.YoutubeShort.IsValidURL(url):
		return provider.YoutubeShort.BuildDataLoader(url)
	case provider.VkVideo.IsValidURL(url):
		return provider.VkVideo.BuildDataLoader(url)
	default:
		return nil, errors.New("no valid URL found")
	}
}

func ExtractURLs(text string) []string {
	switch {
	case provider.Youtube.IsValidURL(text):
		return provider.Youtube.ExtractURLs(text)
	case provider.YoutubeShort.IsValidURL(text):
		return provider.YoutubeShort.ExtractURLs(text)
	case provider.VkVideo.IsValidURL(text):
		return provider.VkVideo.ExtractURLs(text)
	default:
		return make([]string, 0)
	}
}
