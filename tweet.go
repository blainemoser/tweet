package tweet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	go_oauth1 "github.com/blainemoser/go-oauth1"
)

const v2URL = "https://api.twitter.com/2/tweets"

var expects = []string{
	"ConsumerKey",
	"ConsumerSecret",
	"AccessToken",
	"AccessSecret",
	"SignatureMethod",
}

type Tweet struct {
	header  string
	content map[string]interface{}
}

func NewTweet(authCreds map[string]string, properties map[string]interface{}) (*Tweet, error) {
	err := validateAuthCreds(authCreds)
	if err != nil {
		return nil, err
	}
	t := &Tweet{
		header:  getHeader(authCreds),
		content: properties,
	}
	return t, nil
}

func (t *Tweet) Send() (int, error) {
	data, err := t.getData()
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(http.MethodPost, v2URL, data)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", t.header)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	return res.StatusCode, nil
}

func (t *Tweet) getData() (io.Reader, error) {
	tweetJSON, err := json.Marshal(t.content)
	if err != nil {
		return strings.NewReader(""), err
	}
	return strings.NewReader(string(tweetJSON)), nil
}

func validateAuthCreds(authCreds map[string]string) error {
	errs := []string{}
	for _, expected := range expects {
		if authCreds[expected] == "" {
			errs = append(errs, expected)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("missing the following credentials:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func getHeader(authCreds map[string]string) string {
	auth := go_oauth1.OAuth1{
		ConsumerKey:     authCreds["ConsumerKey"],
		ConsumerSecret:  authCreds["ConsumerSecret"],
		AccessToken:     authCreds["AccessToken"],
		AccessSecret:    authCreds["AccessSecret"],
		SignatureMethod: "HMAC-SHA1",
	}
	return auth.BuildOAuth1Header(http.MethodPost, v2URL, map[string]string{})
}
