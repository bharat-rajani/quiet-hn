package main

import (
	"encoding/json"
	"fmt"
	io "io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {

	topStoriesEnpoint := "https://hacker-news.firebaseio.com/v0/topstories.json"
	resp, err := http.Get(topStoriesEnpoint)
	if err != nil {
		fmt.Errorf("%v\n", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("%v\n", err)
	}
	var data []int32
	json.Unmarshal(body, &data)
	// fmt.Println(data)

	stories := make([]interface{}, 500, 500)
	start := time.Now()

	var wg sync.WaitGroup
	for idx, storyId := range data {
		wg.Add(1)
		go func(stories []interface{}, idx int, storyId int32) {
			defer wg.Done()
			sresp, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyId))
			if err != nil {
				fmt.Errorf("Error for %d, %v", storyId, err.Error())
			}
			sbody, err := io.ReadAll(sresp.Body)
			if err != nil {
				fmt.Errorf("%v\n", err)
			}
			var sresponse interface{}
			json.Unmarshal(sbody, &sresponse)
			stories[idx] = sresponse
			// fmt.Println(sresponse)
		}(stories, idx, storyId)
	}
	wg.Wait()
	fmt.Println(time.Since(start))

	// fmt.Println(stories)
}
