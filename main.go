package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Altamashattari/google-photo-downloader/googlealbum"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOAuthConfig *oauth2.Config
	randomState       = "random"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env variables ...")
	}
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/photoslibrary.readonly",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.ListenAndServe(":8080", nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	var html = `
		<html>
			<body>
				<a href="/login">Login Using Google</a>
			</body>
		</html>
	`
	fmt.Fprint(w, html)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != randomState {
		fmt.Println("State is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	token, err := googleOAuthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))

	if err != nil {
		fmt.Println("Error while retrieving token..")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	getData(w, r, token)

}

func getData(w http.ResponseWriter, r *http.Request, token *oauth2.Token) {
	albums, err := googlealbum.GetAllAlbums(token.AccessToken)

	if err != nil {
		fmt.Println("Could make get request...")
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	firstAlbum := albums.Albums[0]
	mediaItems, err := firstAlbum.GetMediaItems(token.AccessToken)

	if err != nil {
		fmt.Println("Couldn't get media items")
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var wg sync.WaitGroup

	for _, mediaItem := range mediaItems {
		wg.Add(1)
		go func(mediaItem googlealbum.MediaItem) {
			defer wg.Done()
			googlealbum.DownloadMediaItem(
				os.Getenv("DOWNLOAD_PATH"),
				&mediaItem,
			)
		}(mediaItem)
	}
	wg.Wait()

	stringifiedAlbums, _ := json.Marshal(albums)

	fmt.Fprintf(w, "Response: %s", stringifiedAlbums)
}
