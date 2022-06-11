package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type client struct {
	host       string
	httpclient httpClient
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
	Get(string) (*http.Response, error)
	Post(string, string, io.Reader) (*http.Response, error)
}

const path = "/v1/organisation/accounts/"

func New(httpclient httpClient, host string) (*client, error) {

	_, err := url.ParseRequestURI(host)
	if err != nil {
		return &client{}, err
	}

	return &client{
		host:       host,
		httpclient: httpclient,
	}, nil
}

func (client client) Create(account *Account) (*Account, error) {

	requesturl := client.host + path

	requestBody := createRequest{
		Data: *account,
	}

	postBody, err := json.Marshal(&requestBody)
	if err != nil {
		return &Account{}, err
	}

	response, err := client.httpclient.Post(requesturl, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return &Account{}, err
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	if response.StatusCode != 201 {
		var errorResponse errorResponse

		decoder.Decode(&errorResponse)
		if err != nil {
			return &Account{}, errors.New("deserialisation error for response")
		}

		return &Account{}, errors.New(errorResponse.ErrorMessage)
	}

	var createResponse createResponse

	decoder.Decode(&createResponse)
	if err != nil {
		return &Account{}, errors.New("deserialisation error for response")
	}

	return &createResponse.Data, err
}

func (client client) Fetch(accountID string) (*Account, error) {

	requesturl := client.host + path + accountID

	response, err := client.httpclient.Get(requesturl)
	if err != nil {
		return &Account{}, err
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	var fetchResponse fetchResponse

	decoder.Decode(&fetchResponse)
	if err != nil {
		return &Account{}, errors.New("deserialisation error for response")
	}

	return &fetchResponse.Data, err
}

func (client client) Delete(accountID string, version string) error {

	requesturl := fmt.Sprintf(client.host+path+"%v?version=%v", accountID, version)

	_, err := url.ParseRequestURI(requesturl)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodDelete, requesturl, nil)
	if err != nil {
		return errors.New("error creating delete request")
	}

	response, err := client.httpclient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 204 {
		return nil
	}

	decoder := json.NewDecoder(response.Body)
	var errorResponse errorResponse

	decoder.Decode(&errorResponse)
	if err != nil {
		return errors.New("deletion request failed, but response error message could not be read")
	}

	return errors.New(errorResponse.ErrorMessage)
}
