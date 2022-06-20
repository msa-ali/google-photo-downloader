package googlealbum

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

const AlbumURL = "https://photoslibrary.googleapis.com/v1/albums"

type Album struct {
	Id                    string `json:"id"`
	Title                 string `json:"title"`
	ProductUrl            string `json:"productUrl"`
	MediaItemsCount       string `json:"mediaItemsCount"`
	CoverPhotoBaseUrl     string `json:"coverPhotoBaseUrl"`
	CoverPhotoMediaItemId string `json:"coverPhotoMediaItemId"`
}

type AlbumsResponse struct {
	Albums        []Album `json:"albums"`
	NextPageToken string  `json:"nextPageToken"`
}

func GetAllAlbums(token string) (albums AlbumsResponse, err error) {
	res, err := http.Get(AlbumURL + "?access_token=" + token)

	if err != nil {
		return AlbumsResponse{}, err
	}

	if res.StatusCode == http.StatusOK {
		content, _ := ioutil.ReadAll(res.Body)
		var data AlbumsResponse
		err = json.Unmarshal(content, &data)
		if err != nil {
			return AlbumsResponse{}, err
		}
		return data, nil
	}

	return AlbumsResponse{}, errors.New("invalid Status code")
}

func (album *Album) DownloadAllMediaItems(token string, downloadPath string, maxRetry int8) error {
	mediaItems, err := album.GetMediaItems(token)
	if err != nil {
		return err
	}
	failedMediaItems := downloadMediaItems(mediaItems, downloadPath)
	var retryCounter int8 = 0
	for {
		if len(failedMediaItems) == 0 || retryCounter < maxRetry {
			break
		}
		failedMediaItems = downloadMediaItems(failedMediaItems, downloadPath)
		retryCounter++
	}

	if len(failedMediaItems) > 0 {
		fmt.Println("Failed Media Items: ", failedMediaItems)
		return errors.New("failed to download few media items")
	}

	return nil
}

func downloadMediaItems(mediaItems []MediaItem, downloadPath string) (failedMediaItems []MediaItem) {
	if len(mediaItems) == 0 {
		return nil
	}
	var wg sync.WaitGroup
	for _, mediaItem := range mediaItems {
		wg.Add(1)
		go func(mediaItem MediaItem) {
			defer wg.Done()
			_, err := mediaItem.download(downloadPath)
			if err != nil {
				failedMediaItems = append(failedMediaItems, mediaItem)
			}
		}(mediaItem)
	}
	wg.Wait()
	if len(failedMediaItems) > 0 {
		return failedMediaItems
	}
	return nil
}
