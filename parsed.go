package main

import (
	"github.com/PuerkitoBio/goquery"
	mapset "github.com/deckarep/golang-set"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func parsedItemsIdsSet(path string) (mapset.Set){
	parsedItemsIdsSet := mapset.NewSet()

	if (filepath.Ext(path) == "") {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			parseFile(path + file.Name(), parsedItemsIdsSet)
		}
	} else {
		parseFile(path, parsedItemsIdsSet)
	}
	return parsedItemsIdsSet
}

func parseFile(filePath string, itemsIdsSet mapset.Set) mapset.Set {
	extension := filepath.Ext(filePath)
	if (extension == ".txt" || extension == ".html"){
		content, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}

		document, err := goquery.NewDocumentFromReader(content)
		if err != nil {
			log.Fatal("Error loading HTTP response body. ", err)
		}

		linkRegExp, _ := regexp.Compile(`(http(s)?:\/\/)?((w){3}.)?youtu(be|.be)?(\.com)?\/.+`)
		videoIdRegExp2, _ := regexp.Compile(`^.*(youtu\.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=)([^#\&\?]*).*`)
		document.Find("a").Each(func(index int, element *goquery.Selection) {
			link, exists := element.Attr("href")
			if exists {
				parserLink := linkRegExp.FindStringSubmatch(link)
				if (len(parserLink)>0) {
					//fmt.Println("link: "+parserLink[0])
					parserdId := videoIdRegExp2.FindStringSubmatch(parserLink[0])
					if (len(parserdId)>0) {
						//fmt.Println("id: "+parserdId[2])
						itemsIdsSet.Add(parserdId[2])
					}
				}
			}
		})
	}
	return itemsIdsSet
}

func sortParsedIds(service *youtube.Service, parsedIds mapset.Set) (mapset.Set, mapset.Set)  {
	validIds := mapset.NewSet()
	invalidIds := mapset.NewSet()
	it := parsedIds.Iterator()
	for id := range it.C {
		video := getVideo(service, id.(string))
		if (len(video.Items) > 0){
			validIds.Add(Video{id.(string), video.Items[0].Snippet.Title})
		} else {
			invalidIds.Add(Video{id.(string), ""})
		}
	}
	return validIds, invalidIds
}

func getVideo(service *youtube.Service, videoId string) *youtube.VideoListResponse {
	call := service.Videos.List("snippet")
	call = call.Id(videoId)
	response, err := call.Do()
	handleError(err, "")
	return response
}