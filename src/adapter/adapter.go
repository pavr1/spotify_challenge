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
	client *http.Client
	config models.SpotifyConfigData
}

//NewAdapterImpl returns an instance of AdapterImpl
func NewAdapter(client *http.Client, config models.SpotifyConfigData) Adapter {
	return AdapterImpl{
		client: client,
		config: config,
	}
}

func (a AdapterImpl) Get(ISRC string) (*models.Metadata, error) {
	t, err := a.authorize()
	if err != nil {
		return nil, err
	}

	metadata, err := a.track(t.Token, ISRC)
	if err != nil {
		return nil, err
	}

	md := models.Metadata{}

	err = json.Unmarshal([]byte(metadata), &md)
	if err != nil {
		return nil, err
	}

	return &md, nil
}

func (a AdapterImpl) authorize() (*models.Token, error) {
	if a.client == nil {
		return nil, errors.New("client not provided")
	}

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, a.config.AuthAPI, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(a.config.ClientID+":"+a.config.Secret)))

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var t models.Token
	err = json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (a AdapterImpl) track(token string, ISRC string) (string, error) {
	if a.client == nil {
		return "", errors.New("client not provided")
	}

	url := fmt.Sprintf(a.config.TrackAPI, ISRC)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
