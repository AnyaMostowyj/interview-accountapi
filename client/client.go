package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func GetClient(httpclient httpClient, host string) *client {

	return &client{
		host:       host,
		httpclient: httpclient,
	}
}

func (client client) GetList() (*ListResponse, error) {

	url := client.host + path

	response, err := client.httpclient.Get(url)
	if err != nil {
		return &ListResponse{}, err
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	var listResponse ListResponse

	decoder.Decode(&listResponse)
	if err != nil {
		return &ListResponse{}, errors.New("deserialisation error for response")
	}

	return &listResponse, err
}

func (client client) CreateAccount(account *Account) (*Account, error) {

	url := client.host + path

	requestBody := CreateRequest{
		Data: *account,
	}

	postBody, err := json.Marshal(&requestBody)
	if err != nil {
		return &Account{}, err
	}

	response, err := client.httpclient.Post(url, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return &Account{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != 201 {
		return &Account{}, errors.New("unexpected response code")
	}

	newdecoder := json.NewDecoder(response.Body)

	var createResponse CreateResponse

	newdecoder.Decode(&createResponse)
	if err != nil {
		return &Account{}, errors.New("deserialisation error for response")
	}

	return &createResponse.Data, err
}

func (client client) GetAccount(accountID string) (*DeprecatedAccount, error) {

	url := client.host + path + accountID

	response, err := client.httpclient.Get(url)
	if err != nil {
		return &DeprecatedAccount{}, err
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	var fetchResponse FetchResponse

	decoder.Decode(&fetchResponse)
	if err != nil {
		return &DeprecatedAccount{}, errors.New("deserialisation error for response")
	}

	return &fetchResponse.Data, err
}

func (client client) DeleteAccount(accountID string, version string) error {

	url := fmt.Sprintf(client.host+path+"%v?version=%v", accountID, version)

	request, err := http.NewRequest(http.MethodDelete, url, nil)
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

	return errors.New("something went wrong")
}
