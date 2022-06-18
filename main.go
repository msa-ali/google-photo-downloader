package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	// GET https://photoslibrary.googleapis.com/v1/albums
	// https://www.googleapis.com/oauth2/v2/userinfo
	getData("https://photoslibrary.googleapis.com/v1/albums", w, r, token)

}

func getData(url string, w http.ResponseWriter, r *http.Request, token *oauth2.Token) {
	// res, err := http.Get(url + "?access_token=" + token.AccessToken)
	// if err != nil {
	// 	fmt.Println("Could make get request...")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// defer r.Body.Close()

	// if res.StatusCode == http.StatusOK {
	// 	content, _ := ioutil.ReadAll(res.Body)
	// 	bodyString := string(content)

	// 	fmt.Println(bodyString)
	// 	fmt.Fprintf(w, "Response: %s", content)
	// }
	googlealbum.GetAllAlbums(token.AccessToken)
}
