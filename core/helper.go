package core

import (
	"net/url"
)

func NewChannelInfo(channelType string, channelId string) *ChannelInfo {
	channel := new(ChannelInfo)
	channel.ChannelType = channelType
	channel.ChannelID = channelId

	return channel
}

func NewWatchRequest(channelInfo *ChannelInfo, rawURL string) (*WatchRequest, error) {
	request := new(WatchRequest)
	request.FeedbackChannelInfo = channelInfo
	URL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}
	request.URL = URL
	request.Interval = 3600

	return request, nil
}
