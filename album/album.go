package google_album

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const AlbumURL = "https://photoslibrary.googleapis.com/v1/albums"

type Album struct {
	Id                    string `json:"id"`
	Title                 string `json:"title"`
	ProductUrl            string `json:"productUrl"`
	MediaItemsCount       string `json:"mediaItemsCount"`
	CoverPhotoBaseUrl     string `json:"coverPhotoBaseUrl"`
	CoverPhotoMediaItemId string `json:"CoverPhotoMediaItemId"`
}

type AlbumsResponse struct {
	Albums        []Album `json:"albums"`
	NextPageToken string `json:"nextPageToken"`
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
