package adapter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Adapter interface {
	Get(ISRC string) ([]byte, error)
}

type AdapterImpl struct {
	spoktifyURL string
	client      *http.Client
}

//NewAdapterImpl returns an instance of AdapterImpl
func NewAdapterImpl(url string, client *http.Client) Adapter {
	return AdapterImpl{
		spoktifyURL: url,
		client:      client,
	}
}

func (a AdapterImpl) Get(ISRC string) ([]byte, error) {
	if a.client == nil {
		return nil, errors.New("client not provided")
	}

	resp, err := a.client.Get(fmt.Sprintf(a.spoktifyURL, ISRC))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
