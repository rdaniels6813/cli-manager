package github_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/rdaniels6813/cli-manager/pkg/github"
	"github.com/rdaniels6813/cli-manager/pkg/promptui"
)

func newServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "Default Message")
	}))
}
func newOauthServer(token string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, fmt.Sprintf(`{"token":"%s"}`, token))
	}))
}
func newOauthServerWithResponse(statusCode int, response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, response)
	}))
}

type LoginFixture struct {
	MockPrompter   *promptui.MockPrompter
	OAuthServer    *httptest.Server
	OAuthOTPServer *httptest.Server
	API            *github.API
}

func (f *LoginFixture) Close() {
	if f.OAuthOTPServer != nil {
		f.OAuthOTPServer.Close()
	}
	if f.OAuthServer != nil {
		f.OAuthServer.Close()
	}
}

func NewLoginFixture(t *testing.T,
	oAuthServer *httptest.Server,
	otpServer *httptest.Server) (*gomock.Controller, *LoginFixture) {
	ctrl := gomock.NewController(t)

	fixture := &LoginFixture{
		MockPrompter:   promptui.NewMockPrompter(ctrl),
		OAuthServer:    oAuthServer,
		OAuthOTPServer: otpServer,
	}
	api, err := github.CreateAPI(func(o *github.API) {
		if fixture.OAuthServer != nil {
			o.OAuthGetURL = fixture.OAuthServer.URL
		}
		if fixture.OAuthOTPServer != nil {
			o.OAuthOTPURL = fixture.OAuthOTPServer.URL
		}
		o.Prompter = fixture.MockPrompter
	})
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	fixture.API = api
	return ctrl, fixture
}

func TestLoginWithUsernamePasswordSucceeds(t *testing.T) {
	expectedToken := "my-awesome-token"
	oAuthServer := newOauthServer(expectedToken)
	otpServer := newServer()
	ctrl, fixture := NewLoginFixture(t, oAuthServer, otpServer)
	defer ctrl.Finish()
	defer fixture.Close()

	fixture.MockPrompter.EXPECT().PromptString(gomock.Any()).Return("username", nil)
	fixture.MockPrompter.EXPECT().PromptPassword(gomock.Any()).Return("password", nil)

	api := fixture.API

	token, err := api.CLILogin()

	assert.Nil(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestLoginWithUsernamePasswordFails(t *testing.T) {
	oauthServer := newOauthServerWithResponse(http.StatusUnauthorized, "Username or password incorrect")
	otpServer := newServer()
	ctrl, fixture := NewLoginFixture(t, oauthServer, otpServer)
	defer ctrl.Finish()
	defer fixture.Close()

	fixture.MockPrompter.EXPECT().PromptString(gomock.Any()).Return("username", nil)
	fixture.MockPrompter.EXPECT().PromptPassword(gomock.Any()).Return("password", nil)

	api := fixture.API

	token, err := api.CLILogin()

	assert.NotNil(t, err)
	assert.Equal(t, "", token)
}

func TestLoginOTPSuccess(t *testing.T) {

}
