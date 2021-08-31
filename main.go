package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/budhalantara/audiusdl/models"
)

var isTerminal bool

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("track_url is needed: audius <track_url>")
		return
	}

	fileInfo, _ := os.Stdout.Stat()
	isTerminal = (fileInfo.Mode() & os.ModeCharDevice) != 0

	discoveryNodesJSON, err := os.ReadFile("./storage/discovery_nodes.json")
	if err != nil {
		panic(err)
	}

	var discoveryNodes []string
	err = json.Unmarshal(discoveryNodesJSON, &discoveryNodes)
	if err != nil {
		panic(err)
	}

	trackUrl := os.Args[1]
	r := regexp.MustCompile(`audius\.co\/(\w.+)\/(\w.+)`)
	match := r.FindStringSubmatch(trackUrl)
	artistId := match[1]
	trackId := match[2]

	resp, err := getContent(discoveryNodes, fmt.Sprintf("/v1/full/tracks?handle=%s&slug=%s", artistId, trackId), "", "application/json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var trackDetail models.Track
	err = json.Unmarshal(resp.Body(), &trackDetail)
	if err != nil {
		fmt.Println("Error:", "Unable to unmarshal trackDetail ->", err)
	}

	fmt.Println()

	var wg sync.WaitGroup
	creatorNodes := strings.Split(trackDetail.Data.User.CreatorNodeEndpoint, ",")
	channel := make(chan models.TrackBuffer)
	for i, segment := range trackDetail.Data.TrackSegments {
		hash := segment.Multihash
		if isTerminal {
			fmt.Printf("[%d]started %s\n", i, hash)
		}
		go Download(&wg, channel, &creatorNodes, hash, i)
	}

	var trackBuffer []byte
	trackSegmentsLen := len(trackDetail.Data.TrackSegments)
	chunks := make([][]byte, trackSegmentsLen)
	for i := 0; i < trackSegmentsLen; i++ {
		buffer := <-channel
		chunks[buffer.ID] = buffer.Data
		if isTerminal {
			fmt.Println(fmt.Sprintf("[%d]Done", buffer.ID), buffer.Hash)
		}
	}

	wg.Wait()
	close(channel)

	for _, data := range chunks {
		trackBuffer = append(trackBuffer, data...)
	}

	if isTerminal {
		dest := fmt.Sprintf("./%s.mpeg", trackDetail.Data.Title)
		fmt.Println("Saved to:", dest)
		err := os.WriteFile(dest, trackBuffer, 0755)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Print(string(trackBuffer))
	}
}
