package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/budhalantara/audiusdl/models"

	"github.com/go-resty/resty/v2"
)

func getContent(endpoints []string, path string, logPrefix string, contentType string) (*resty.Response, error) {
	for _, endp := range endpoints {
		if isTerminal {
			fmt.Println(fmt.Sprintf("%sTrying endpoint:", logPrefix), endp)
		}
		resp, err := resty.New().R().
			SetHeader("Accept-Encoding", "gzip, deflate, br").
			Get(fmt.Sprintf("%s%s", endp, path))
		if err != nil {
			return resp, err
		}

		if contentType != "" {
			// fmt.Println(resp.Header().Get("Content-Type"))
			if !strings.Contains(resp.Header().Get("Content-Type"), contentType) {
				continue
			}
		}

		if resp.StatusCode() == 200 {
			if isTerminal {
				fmt.Println(fmt.Sprintf("%sSelected endpoint:", logPrefix), endp)
			}
			return resp, nil
		}
	}
	return &resty.Response{}, errors.New("unable to use provided endpoints")
}

func Download(wg *sync.WaitGroup, buffer chan<- models.TrackBuffer, endpoints *[]string, multihash string, id int) {
	defer wg.Done()

	wg.Add(1)

	resp, err := getContent(*endpoints, fmt.Sprintf("/ipfs/%s", multihash), fmt.Sprintf("[%d]", id), "")
	if err != nil {
		panic(err)
	}

	data := models.TrackBuffer{ID: id, Data: resp.Body(), Hash: multihash}

	buffer <- data
}
