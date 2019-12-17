package main

import (
	mapset "github.com/deckarep/golang-set"
	"google.golang.org/api/youtube/v3"
)

func playlistItemsList(service *youtube.Service, part string, playlistId string) *youtube.PlaylistItemListResponse {
	call := service.PlaylistItems.List(part)
	call = call.PlaylistId(playlistId)
	response, err := call.Do()
	handleError(err, "")
	return response
}

func playlistItemsIdsSet(service *youtube.Service, playlistId string) mapset.Set {
	playlistItemsIdsSet := mapset.NewSet()
	for _, playlistItem := range playlistItemsList(service, "snippet", playlistId).Items {
		playlistItem := Video{playlistItem.Snippet.ResourceId.VideoId,playlistItem.Snippet.Title}
		playlistItemsIdsSet.Add(playlistItem)
	}
	return playlistItemsIdsSet
}


