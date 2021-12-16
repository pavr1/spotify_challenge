package adapter

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"spotify_challenge.com/src/models"
)

type Adapter interface {
	Get(ISRC string) (*models.Metadata, error)
}

type AdapterImpl struct {
	spoktifyURL string
	client      *http.Client
	config      models.SpotifyData
}

//NewAdapterImpl returns an instance of AdapterImpl
func NewAdapterImpl(url string, client *http.Client, config models.SpotifyData) Adapter {
	return AdapterImpl{
		spoktifyURL: url,
		client:      client,
		config:      config,
	}
}

func (a AdapterImpl) Get(ISRC string) (*models.Metadata, error) {
	t, err := a.authorize()
	if err != nil {
		return nil, err
	}
	metadata, err := a.track(t, ISRC)
	if err != nil {
		return nil, err
	}

	md := models.Metadata{}

	err = json.Unmarshal([]byte(metadata), &md)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (a AdapterImpl) authorize() (string, error) {
	if a.client == nil {
		return "", errors.New("client not provided")
	}

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", a.config.AuthAPI, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("Basic "+a.config.ClientID+":"+a.config.Secret)))

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (a AdapterImpl) track(token string, ISRC string) (string, error) {
	if a.client == nil {
		return "", errors.New("client not provided")
	}

	req, err := http.NewRequest("POST", a.config.TrackAPI, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
