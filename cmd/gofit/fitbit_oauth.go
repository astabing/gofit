package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

var oauthConfig *oauth2.Config

var errExpiredToken = errors.New("Expired token")
var errInvalidToken = errors.New("Invalid token")
var errFitbitUndefEnv = errors.New("Fitbit API env variables undefined")

// From https://github.com/golang/oauth2/issues/84#issuecomment-332517319
// TokenNotifyFunc is a function that accepts an oauth2 Token upon refresh, and
// returns an error if it should not be used.
type tokenNotifyFunc func(string, *oauth2.Token) error

// NotifyRefreshTokenSource is essentially `oauth2.ResuseTokenSource` with `TokenNotifyFunc` added.
type notifyRefreshTokenSource struct {
	new oauth2.TokenSource
	mux sync.Mutex // guards tok
	tok *oauth2.Token
	tcf string          // token cache file
	fn  tokenNotifyFunc // called when token refreshed so new refresh token can be persisted
}

func storeTokenCache(tcf string, tok *oauth2.Token) error {
	tokData, err := json.Marshal(&tok)
	if err != nil {
		return fmt.Errorf("Could not encode token to json: %w", err)
	}
	// persist token
	return os.WriteFile(tcf, tokData, 0600)
}

func fetchTokenCache(tcf string) (*oauth2.Token, error) {
	tok := new(oauth2.Token)

	tokData, err := os.ReadFile(tcf)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(tokData, tok); err != nil {
		return nil, err
	}
	if !tok.Valid() && tok == nil {
		// Allow expired token, but not nil
		return nil, errInvalidToken
	}
	return tok, nil
}

// Token checks is *NotifyRefreshTokenSource token is valid and then returns it
// If token is not valid it will try to refresh it using provided oauth2.TokenSource
// Refreshed token is returned and written to token cache storage
func (src *notifyRefreshTokenSource) Token() (*oauth2.Token, error) {
	src.mux.Lock()
	defer src.mux.Unlock()
	if src.tok.Valid() {
		return src.tok, nil
	}
	// Refresh token
	tok, err := src.new.Token()
	if err != nil {
		return nil, err
	}
	src.tok = tok
	// Return refreshed token and write it to token cache storage
	return tok, src.fn(src.tcf, tok)
}

type transport struct{}

// Accept-Language header sets the measurement unit system to use for response values
// https://dev.fitbit.com/build/reference/web-api/developer-guide/application-design/#Localization
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	//req.Header.Add("Accept-Language", "en_US")

	// Debug dump
	//reqDump, err := httputil.DumpRequestOut(req, true)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Printf("REQUEST:\n%s", string(reqDump))

	return http.DefaultTransport.RoundTrip(req)
}

// GetClient initializes and returns Fitbit oauth2 *http.Client
func GetOauthClient() (*http.Client, error) {

	clientID := os.Getenv("FITBIT_CLIENTID")
	clientSecret := os.Getenv("FITBIT_SECRET")
	oauthRedirectURL := os.Getenv("FITBIT_REDIRURL")
	tokenCacheFile := os.Getenv("FITBIT_CACHEFILE")

	clientScopes := []string{"activity", "heartrate", "location", "profile", "nutrition", "sleep", "weight"}

	if clientID == "" || clientSecret == "" || oauthRedirectURL == "" || tokenCacheFile == "" {
		return nil, errFitbitUndefEnv
	}

	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       clientScopes,
		Endpoint:     fitbit.Endpoint,
		RedirectURL:  oauthRedirectURL,
	}

	hc := &http.Client{Transport: &transport{}}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, hc)

	tok, err := fetchTokenCache(tokenCacheFile)
	if err != nil {
		// Unable to retrieve cached token, try to get new
		log.Printf("Unable to retrieve cached token: %v", err)
		tok, err = oauthConfig.Exchange(ctx, authCode())
		if err != nil {
			return nil, err
		}

		if err = storeTokenCache(tokenCacheFile, tok); err != nil {
			log.Printf("Unable to persist token: %v", err)
		}
	}

	tokSrc := &notifyRefreshTokenSource{
		new: oauthConfig.TokenSource(ctx, tok),
		tok: tok,
		tcf: tokenCacheFile,
		fn:  storeTokenCache,
	}

	return oauth2.NewClient(ctx, tokSrc), nil
}

func authCode() string {
	fmt.Printf("Get auth code from:\n\n")
	fmt.Printf("%v\n\n", oauthConfig.AuthCodeURL("fit-auth-state", oauth2.AccessTypeOffline))
	return prompt("Enter auth code:")

}

// Prompt the user for an input line.  Return the given input.
func prompt(promptText string) string {
	fmt.Print(promptText)
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	return strings.TrimSpace(sc.Text())
}
