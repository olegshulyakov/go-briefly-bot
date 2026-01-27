package loaders

import (
	"errors"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/provider"
)

func VideoLoader(url string) (*video.DataLoader, error) {
	if provider.Youtube.IsValidURL(url) {
		return provider.Youtube.BuildDataLoader(url)
	} else if provider.YoutubeShort.IsValidURL(url) {
		return provider.YoutubeShort.BuildDataLoader(url)
	} else if provider.VkVideo.IsValidURL(url) {
		return provider.VkVideo.BuildDataLoader(url)
	}
	return nil, errors.New("no valid URL found")
}

func ExtractURLs(text string) []string {
	if provider.Youtube.IsValidURL(text) {
		return provider.Youtube.ExtractURLs(text)
	} else if provider.YoutubeShort.IsValidURL(text) {
		return provider.YoutubeShort.ExtractURLs(text)
	} else if provider.VkVideo.IsValidURL(text) {
		return provider.VkVideo.ExtractURLs(text)
	}
	return make([]string, 0)
}
