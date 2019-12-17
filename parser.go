package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/api/youtube/v3"
)

const pathIndex = 1

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

type Video struct {
	Id string
	Title string
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message + ": %v", err.Error())
	}
}

// Retrieve resource for the authenticated user's channel
func playlistsListMine(service *youtube.Service, part string) *youtube.PlaylistListResponse {
	call := service.Playlists.List(part)
	call = call.Id("PL1eQAFFC9DaMS2zQBiLxW8EPcoYCkzZug")
	response, err := call.Do()
	handleError(err, "")
	return response
}

func removeIdsFromList(set mapset.Set, input string){
	slice := strings.Split(input, ",")
	for _, id := range slice {
		set.Remove(strings.TrimSpace(id))
	}
}

func addIdsToList(set mapset.Set, input string){
	slice := strings.Split(input, ",")
	for _, id := range slice {
		set.Add(strings.TrimSpace(id))
	}
}

func pushVideosToPlaylist(service *youtube.Service, set mapset.Set, playlistId string){
	it := set.Iterator()
	for video := range it.C {
		playlistItemInsert(service, video.(Video).Id, playlistId)
		println(video.(Video).Id + "   " + video.(Video).Title)
	}
}

func playlistItemInsert(service *youtube.Service, videoId string, playlistId string) {
	resourceId := youtube.ResourceId{
		Kind:    "youtube#video",
		VideoId: videoId,
	}

	playlistItemSnippet := youtube.PlaylistItemSnippet{
		PlaylistId: playlistId,
		ResourceId: &resourceId,
	}

	playlistItem :=  youtube.PlaylistItem{
		Snippet: &playlistItemSnippet,
	}

	call := service.PlaylistItems.Insert("id", &playlistItem)
	_, err := call.Do()
	handleError(err, "")
}

func printList(set mapset.Set, headerMessage string)  {
	println(headerMessage)
	it := set.Iterator()
	for video := range it.C {
		println(video.(Video).Id + "   " + video.(Video).Title)
	}
}

func main() {
	fmt.Println("Hey Dude! I know you have a lot of youtube links in your logs. Lets push them to one playlist.")
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := GetClient(ctx, config)

	service, err := youtube.New(client)

	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	var playlistId string
	playlistId = `PL1eQAFFC9DaMS2zQBiLxW8EPcoYCkzZug`
	path := os.Args[pathIndex]

	//PL1eQAFFC9DaMS2zQBiLxW8EPcoYCkzZug
	var command string

	newPlaylistSet := parsedItemsIdsSet(path)

	validSet, invalidIdsSet := sortParsedIds(service, newPlaylistSet)
	currentPlaylistSet := playlistItemsIdsSet(service, playlistId)
	setToPush := validSet.Difference(currentPlaylistSet)
	addedSet := validSet.Intersect(currentPlaylistSet)

	printList(setToPush, "We are going to push this video: ")
	printList(addedSet, "This video is already exists in your playlist: ")
	printList(invalidIdsSet, "This video has broken links: ")

	println(`If list is ready to push just print "apply" or "exit" to exit from the app`)

	fmt.Scanf("%s", &command)
	pushVideosToPlaylist(service, setToPush, playlistId)
	println(`Chanson débile. Bye-bye!`)
}