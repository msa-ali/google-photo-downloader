package googlealbum

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const MediaItemsUrl = "https://photoslibrary.googleapis.com/v1/mediaItems:search"

type mediaItemBody struct {
	PageSize      string `json:"pageSize"`
	AlbumId       string `json:"albumId"`
	NextPageToken string `json:"pageToken"`
}

type VideoMetadata struct {
	CameraMake  string `json:"cameraMake"`
	CameraModel string `json:"cameraModel"`
	Fps         string `json:"fps"`
	Status      string `json:"status"`
}

type MediaMetadata struct {
	Width        string        `json:"width"`
	Height       string        `json:"height"`
	CreationTime string        `json:"creationTime"`
	Video        VideoMetadata `json:"video"`
}

type ContributorInfo struct {
	ProfilePictureBaseUrl string `json:"profilePictureBaseUrl"`
	DisplayName           string `json:"displayName"`
}

type MediaItem struct {
	Id              string          `json:"id"`
	Description     string          `json:"description"`
	ProductUrl      string          `json:"productUrl"`
	BaseUrl         string          `json:"baseUrl"`
	MimeType        string          `json:"mimeType"`
	Filename        string          `json:"filename"`
	MediaMetadata   MediaMetadata   `json:"mediaMetadata"`
	ContributorInfo ContributorInfo `json:"contributorInfo"`
}

func (album *Album) GetMediaItems(token string) ([]MediaItem, error) {
	var nextPageToken string
	var mediaItems []MediaItem
	for {
		bodyObj := mediaItemBody{
			PageSize: "100",
			AlbumId:  album.Id,
		}
		if nextPageToken != "" {
			bodyObj.NextPageToken = nextPageToken
		}
		bodyJson, err := json.Marshal(bodyObj)
		if err != nil {
			return nil, err
		}

		body := []byte(bodyJson)

		r, err := http.NewRequest("POST", MediaItemsUrl, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add("Authorization", "Bearer "+token)

		client := &http.Client{}
		res, err := client.Do(r)

		if res.StatusCode != http.StatusOK {
			return nil, errors.New(res.Status)
		}

		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		var data interface{}
		derr := json.NewDecoder(res.Body).Decode(&data)
		if derr != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			return nil, err
		}

		md, _ := data.(map[string]interface{})
		items, _ := md["mediaItems"].([]interface{})

		nextPageToken, _ = md["nextPageToken"].(string)
		fmt.Println(nextPageToken)

		for _, item := range items {
			media := item.(map[string]interface{})
			var mediaItem MediaItem
			mediaItem.Id = media["id"].(string)
			mediaItem.BaseUrl = media["baseUrl"].(string)
			mediaItem.Filename = media["filename"].(string)
			mediaItems = append(mediaItems, mediaItem)
		}
		// if next page token is available, it means media items are available in next page
		if nextPageToken == "" {
			break
		}
	}
	return mediaItems, nil
}

func (mediaItem *MediaItem) download(filePath string) (int64, error) {
	response, err := http.Get(mediaItem.BaseUrl + "=d")
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	file, err := os.Create(filePath + "/" + mediaItem.Filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	written, err := io.Copy(file, response.Body)
	if err != nil {
		return 0, err
	}
	return written, nil
}
