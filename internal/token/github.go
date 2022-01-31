package token

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/cli/cli/api"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/cli/oauth"
	"github.com/rdaniels6813/cli-manager/internal/store"
)

type TokenManager interface {
	GetNewOrSavedToken(scopes []string) (string, error)
}

type OSTokenManager struct {
	store store.Store
}

// NewConfigTokenManager uses a store interface to persist tokens
func NewConfigTokenManager(store store.Store) *OSTokenManager {
	return &OSTokenManager{store: store}
}

const serviceName = "cli-manager:github"

func scopesToAccount(scopes []string) string {
	sort.Strings(scopes)
	return strings.Join(scopes, ";")
}

func (t *OSTokenManager) GetNewOrSavedToken(scopes []string) (string, error) {
	token, err := t.store.Get(serviceName, scopesToAccount(scopes))
	if token == "" || err != nil {
		newToken, err := authFlow("github.com", &iostreams.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}, "", scopes)
		if err != nil {
			return "", err
		}
		err = t.store.Set(serviceName, scopesToAccount(scopes), newToken)
		if err != nil {
			fmt.Printf("failed saving token to keyring: %s", err)
		}
		return newToken, nil
	}
	return token, nil
}

func authFlow(oauthHost string, inputOutput *iostreams.IOStreams, notice string, additionalScopes []string) (string, error) {
	w := inputOutput.ErrOut
	cs := inputOutput.ColorScheme()

	httpClient := http.DefaultClient
	if envDebug := os.Getenv("DEBUG"); envDebug != "" {
		logTraffic := strings.Contains(envDebug, "api") || strings.Contains(envDebug, "oauth")
		httpClient.Transport = api.VerboseLog(inputOutput.ErrOut, logTraffic, inputOutput.ColorEnabled())(httpClient.Transport)
	}

	scopes := []string{"repo", "read:org", "gist"}
	scopes = append(scopes, additionalScopes...)

	callbackURI := "http://127.0.0.1/callback"
	host := oauth.GitHubHost("https://github.com")
	flow := &oauth.Flow{
		Host:        host,
		ClientID:    "c6f78436de6ccad2fb30",
		CallbackURI: callbackURI,
		Scopes:      scopes,
		DisplayCode: func(code, verificationURL string) error {
			fmt.Fprintf(w, "%s First copy your one-time code: %s\n", cs.Yellow("!"), cs.Bold(code))
			return nil
		},
		BrowseURL: func(url string) error {
			fmt.Fprintf(w, "- %s to open %s in your browser... ", cs.Bold("Press Enter"), oauthHost)
			browser := cmdutil.NewBrowser(os.Getenv("BROWSER"), inputOutput.Out, inputOutput.ErrOut)
			if err := browser.Browse(url); err != nil {
				fmt.Fprintf(w, "%s Failed opening a web browser at %s\n", cs.Red("!"), url)
				fmt.Fprintf(w, "  %s\n", err)
				fmt.Fprint(w, "  Please try entering the URL in your browser manually\n")
			}
			return nil
		},
		WriteSuccessHTML: func(w io.Writer) {
			fmt.Fprint(w, oauthSuccessPage)
		},
		HTTPClient: httpClient,
		Stdin:      inputOutput.In,
		Stdout:     w,
	}

	fmt.Fprintln(w, notice)

	token, err := flow.DetectFlow()
	if err != nil {
		return "", err
	}

	return token.Token, nil
}

const oauthSuccessPage = `
<!doctype html>
<meta charset="utf-8">
<title>Success: cli-manager</title>
<style type="text/css">
body {
  color: #1B1F23;
  background: #F6F8FA;
  font-size: 14px;
  font-family: -apple-system, "Segoe UI", Helvetica, Arial, sans-serif;
  line-height: 1.5;
  max-width: 620px;
  margin: 28px auto;
  text-align: center;
}
h1 {
  font-size: 24px;
  margin-bottom: 0;
}
p {
  margin-top: 0;
}
.box {
  border: 1px solid #E1E4E8;
  background: white;
  padding: 24px;
  margin: 28px;
}
</style>
<body>
  <svg height="52" class="octicon octicon-mark-github" viewBox="0 0 16 16" version="1.1" width="52" aria-hidden="true"><path fill-rule="evenodd" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"></path></svg>
  <div class="box">
    <h1>Successfully authenticated cli-manager</h1>
    <p>You may now close this tab and return to the terminal.</p>
  </div>
</body>
`
