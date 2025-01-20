package goauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type installed struct {
	ClientID            string   `json:"client_id"`
	ProjectID           string   `json:"project_id"`
	AuthURI             string   `json:"auth_uri"`
	TokenURI            string   `json:"token_uri"`
	AuthProviderCertURL string   `json:"auth_provider_x509_cert_url"`
	ClientSecret        string   `json:"client_secret"`
	RedirectURIs        []string `json:"redirect_uris"`
}

type credentials struct {
	Installed installed `json:"installed"`
}

type OAuth struct {
	creds       credentials
	oauthConfig *oauth2.Config
	Client      *http.Client
}

// initializeCredentials method initializes the credentials for the OAuth client. It reads the GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables.
func (o *OAuth) initializeCredentials() {
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")

	if googleClientId == "" {
		log.Fatalln("GOOGLE_CLIENT_ID is not set")
	}

	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if googleClientSecret == "" {
		log.Fatalln("GOOGLE_CLIENT_SECRET is not set")
	}

	o.creds = credentials{
		Installed: installed{
			ProjectID:           "jiyeol-tech",
			AuthURI:             "https://accounts.google.com/o/oauth2/auth",
			TokenURI:            "https://oauth2.googleapis.com/token",
			AuthProviderCertURL: "https://www.googleapis.com/oauth2/v1/certs",
			RedirectURIs:        []string{"http://localhost:8000/callback"},
			ClientID:            googleClientId,
			ClientSecret:        googleClientSecret,
		},
	}
}

// setOAuthConfig method returns a new OAuth2 config.
func (o *OAuth) setOAuthConfig(scope ...string) error {
	o.initializeCredentials()

	b, err := json.Marshal(o.creds)
	if err != nil {
		return err
	}

	config, err := google.ConfigFromJSON(b, scope...)

	o.oauthConfig = config

	return err
}

// SetClient function retrieves a token, saves the token, then returns the generated client.
func (o *OAuth) SetClient(scope ...string) error {
	if o.Client != nil {
		return nil
	}

	err := o.setOAuthConfig(scope...)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	tokFile := home + "/token.json"
	tok, err := getTokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(o.oauthConfig)
		if err != nil {
			return err
		}
		err := saveToken(tokFile, tok)
		if err != nil {
			return err
		}
	}

	o.Client = o.oauthConfig.Client(context.Background(), tok)

	return nil
}

// getTokenFromWeb function retrieves a token from the web.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	err := exec.Command("open", authURL).Start()
	if err != nil {
		return nil, err
	}

	listenAddr := "localhost:8000"
	redirectURI := fmt.Sprintf("http://%s/callback", listenAddr)
	config.RedirectURL = redirectURI

	authCodeCh := make(chan string)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Authorization code not found", http.StatusBadRequest)
			return
		}

		authCodeCh <- code
		fmt.Fprintln(w, "Authorization code received. You can close this window.")
	})
	go func() {
		log.Fatal(http.ListenAndServe(listenAddr, nil))
	}()

	authCode := <-authCodeCh

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

// getTokenFromFile function retrieves a token from a local file.
func getTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken function saves a token to a local file.
func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return nil
}
