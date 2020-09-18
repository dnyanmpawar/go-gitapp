package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"github.com/gosimple/oauth2"
	"io/ioutil"
	"net/http"
)

var service *oauth2.OAuth2Service
var repoDatas = []repoInfo{}

type configDefinition struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackURL  string `json:"callback_url"`
}

var config configDefinition

type userInfo struct {
	Name     string `json:"Name"`
	Login    string `json:"Login"`
	ReposURL string `json:"repos_url"`
}

type repoInfo struct {
	CloneURL string `json:"clone_url"`
	Name     string `json:"name"`
}

func github_callback(w http.ResponseWriter, r *http.Request) {

	codes := r.URL.Query()["code"]
	code := codes[0]
	if code == "" {
		log.Println("Invalid Code")
		return
	}
	// Get access token.
	token, err := service.GetAccessToken(code)
	if err != nil {
		fmt.Println("Get access token error: ", err)
		return
	}

	// Prepare resource request
	apiBaseURL := "https://api.github.com/"

	github := oauth2.Request(apiBaseURL, token.AccessToken)
	github.AccessTokenInHeader = true
	github.AccessTokenInHeaderScheme = "token"

	// Get User Info
	apiEndPoint := "user"
	githubUserData, err := github.Get(apiEndPoint)
	if err != nil {
		log.Println("Get: ", err)
		return
	}
	defer githubUserData.Body.Close()

	body, err := ioutil.ReadAll(githubUserData.Body)
	if err != nil {
		log.Println("Error reading response")
		return
	}

	response := userInfo{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Failed to unmarshal user info")
		return
	}
	endpoint := "users/" + response.Login + "/repos"

	repoData, err := github.Get(endpoint)
	if err != nil {
		fmt.Println("Error reading response")
		return
	}
	defer repoData.Body.Close()

	body, err = ioutil.ReadAll(repoData.Body)
	if err != nil {
		fmt.Println("Error reading response")
		return
	}

	err = json.Unmarshal(body, &repoDatas)
	if err != nil{
		log.Println("Failed to unmarshal Repo Info")
	}

	http.Redirect(w, r, "/display_repos", 301)

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go app integration with github")
}

func loadLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Loading github login page")

	//service.RedirectURL = "http://127.0.0.1:80/github_callback"
	service.RedirectURL = config.CallbackURL

	// Get authorization url.
	authUrl := service.GetAuthorizeURL("")
	http.Redirect(w, r, authUrl, 301)
}


func main() {
	fmt.Println("Starting Server")
	http.HandleFunc("/", DisplayConnectGithub)
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/connect_github", loadLogin)
	http.HandleFunc("/github_callback", github_callback)
	http.HandleFunc("/display_repos", DisplayRepos)
	http.HandleFunc("/selected", UserSelected)

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("file config.json not found ")
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("invalid format of config file")
	}

	service = oauth2.Service(
		config.ClientID,
		config.ClientSecret,
		"https://github.com/login/oauth/authorize",
		"https://github.com/login/oauth/access_token",
	)
	errs := make(chan error)
	go func() {
		errs <- http.ListenAndServe(":80", nil)
	}()
	fmt.Println("Exiting with Status ", <-errs)
}
