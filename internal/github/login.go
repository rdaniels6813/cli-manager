package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rdaniels6813/cli-manager/pkg/promptui"
)

// API api struct for logging into github and retrieving a new personal access token
type API struct {
	OAuthGetURL string
	OAuthOTPURL string
	Prompter    promptui.Prompter
}

// CreateAPI create an api struct for logging into github and retrieving a new personal access token
func CreateAPI(options ...func(*API)) (*API, error) {
	api := &API{
		OAuthGetURL: "https://api.github.com/user",
		OAuthOTPURL: "https://api.github.com/authorizations",
	}
	for _, option := range options {
		option(api)
	}
	return api, nil
}

// CLILogin prompt the user via the command line for credentials to log into github
func (g *API) CLILogin(scopes ...string) (string, error) {
	username, err := g.getUserName()
	if err != nil {
		return "", err
	}
	password, err := g.getPassword()
	if err != nil {
		return "", err
	}
	return g.tryLogin(username, password)
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (g *API) tryLogin(username string, password string) (string, error) {
	req, err := http.NewRequest("GET", g.OAuthGetURL, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, password)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode == 401 {
		otpHeader := res.Header.Get("x-github-otp")
		if otpHeader == "required: app;" || otpHeader == "required: sms;" {
			return g.loginWithOTP(username, password, otpHeader)
		}
		return "", fmt.Errorf("Request failed with status: %s", res.Status)
	}
	defer res.Body.Close()
	var token tokenResponse
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func (g *API) loginWithOTP(username string, password string, otpHeader string) (string, error) {
	req, err := http.NewRequest("POST", g.OAuthOTPURL, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, password)
	otp, err := g.getOTP(otpHeader)
	if err != nil {
		return "", err
	}
	req.Header.Set("x-github-otp", otp)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("Request failed with status: %s", res.Status)
	}
	defer res.Body.Close()
	return "", nil
}

func (g *API) getOTP(header string) (string, error) {
	otpType := "Authenticator App"
	if header == "required: sms;" {
		otpType = "Sent via SMS"
	}
	return g.Prompter.PromptString(fmt.Sprintf("GitHub One-Time Passcode (%s)", otpType))
}
func (g *API) getUserName() (string, error) {
	return g.Prompter.PromptString("GitHub Username")
}
func (g *API) getPassword() (string, error) {
	return g.Prompter.PromptPassword("GitHub Password")
}
