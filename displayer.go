package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"crypto/rand"
)

type RadioButton struct {
	Name       string
	Value      string
	IsDisabled bool
	IsChecked  bool
	Text       string
}

type PageVariables struct {
	PageTitle        string
	PageRadioButtons []RadioButton
	Answer           string
}

func DisplayConnectGithub(w http.ResponseWriter, r *http.Request) {
	Title := "Connect to github"
	MyRadioButtons := []RadioButton{
		RadioButton{"ConnectGithub", "connect_github", false, false, "Connect Github"},
	}
	MyPageVariables := PageVariables{
		PageTitle:        Title,
		PageRadioButtons: MyRadioButtons,
	}
	t, err := template.ParseFiles("home.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, MyPageVariables)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func DisplayRepos(w http.ResponseWriter, r *http.Request) {

	Title := "Clone the repo"
	MyRadioButtons := []RadioButton{}
	for _, v := range repoDatas {
		fmt.Printf(v.Name)
		rb := RadioButton{
			Name:       "RepoName",
			Value:      v.Name,
			IsDisabled: false,
			IsChecked:  false,
			Text:       v.Name,
		}
		MyRadioButtons = append(MyRadioButtons, rb)
	}

	MyPageVariables := PageVariables{
		PageTitle:        Title,
		PageRadioButtons: MyRadioButtons,
	}

	t, err := template.ParseFiles("select.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, MyPageVariables)
	if err != nil {
		log.Print("template executing error: ", err)
	}

}

func UserSelected(w http.ResponseWriter, r *http.Request) {
	log.Println("starting UserSelected")
	r.ParseForm()

	selectedRepoName := r.Form.Get("RepoName")

	Title := "Selected Repo"
	MyPageVariables := PageVariables{
		PageTitle: Title,
		Answer:    selectedRepoName,
	}

	var cloneURL string
	for _, v := range repoDatas {
		if selectedRepoName == v.Name {
			cloneURL = v.CloneURL
		}
	}

	cmd := exec.Command("git", "clone", cloneURL)
	err := cmd.Run()
	if err != nil {
		log.Fatal("Error Cloning Repo:", err)
	}

	cmd = exec.Command("cd", selectedRepoName)
	err = cmd.Run()
	if err != nil {
		log.Fatal("Error in change dir:", err)
	}

	randint, _ := rand.Prime(rand.Reader, 32)
	newbranch := randint.String()
	cmd = exec.Command("git", "checkout", "-b", newbranch)
	err = cmd.Run()
	if err != nil {
		log.Fatal("Error in checkout new branch:", err)
	}

	cmd = exec.Command("touch", newbranch)
	err = cmd.Run()
	if err != nil {
		log.Fatal("Error in creating file:", err)
	}

	cmd = exec.Command("git", "add", newbranch)
	err = cmd.Run()
	if err != nil {
		log.Fatal("Error in adding file:", err)
	}

	cmd = exec.Command("git", "commit", "-a", "-m", newbranch)
	err = cmd.Run()
	if err != nil {
		log.Fatal("Error in commit:", err)
	}

	cmd = exec.Command("git", "push", "origin", "master")
	err = cmd.Run()
	if err != nil {
		log.Fatal("Error in push:", err)
	}

	t, err := template.ParseFiles("select.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, MyPageVariables)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}
