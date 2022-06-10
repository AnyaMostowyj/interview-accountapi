package client_test

import (
	"client"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const host = "http://accountapi:8080"

//const host = "http://localhost:8080"

func TestCreateAccount(t *testing.T) {

	jsonFile, err := os.Open("createaccount.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	var account client.Account
	json.Unmarshal(bytes, &account)

	fmt.Println("Account type: " + account.Type)

	unitUnderTest, _ := client.New(&http.Client{}, host)

	accountResponse, err := unitUnderTest.CreateAccount(&account)
	if err != nil {
		t.Error("Account creation error")
	}

	fmt.Println("Account response type: " + accountResponse.Type)
}

func TestGetAccount(t *testing.T) {

	expectedAccountID := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

	unitUnderTest, _ := client.New(&http.Client{}, host)

	account, err := unitUnderTest.GetAccount(expectedAccountID)

	if err != nil {
		t.Error("It didn't work")
	}

	if account.Id != expectedAccountID {
		t.Errorf("Expected account ID: '%v' but got account ID: '%v'", expectedAccountID, account.Id)
	}
}

func TestDeleteAccount(t *testing.T) {

	expectedAccountID := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

	unitUnderTest, _ := client.New(&http.Client{}, host)

	err := unitUnderTest.DeleteAccount(expectedAccountID, "0")

	if err != nil {
		t.Error("It didn't work")
	}
}

// These mocked test scenarios cover failure responses difficult to trigger via integration tests
func TestCreateAccountWithMock(t *testing.T) {

	jsonFile, err := os.Open("createaccount.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	var account client.Account
	json.Unmarshal(bytes, &account)

	mux := http.NewServeMux()

	testServer := httptest.NewServer(mux)

	mux.HandleFunc("/v1/organisation/accounts/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(bytes)
	})

	defer testServer.Close()

	unitUnderTest, _ := client.New(&http.Client{}, testServer.URL)

	_, err = unitUnderTest.CreateAccount(&account)

	if err != nil {
		t.Errorf("Test request failed with error: '%v'", err.Error())
	}
}

func TestCreateCopAccountWithMock(t *testing.T) {

	jsonFile, err := os.Open("./test_data/cop_request.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	var account client.Account
	json.Unmarshal(bytes, &account)

	mux := http.NewServeMux()

	testServer := httptest.NewServer(mux)

	mux.HandleFunc("/v1/organisation/accounts/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(bytes)
	})

	defer testServer.Close()

	unitUnderTest, _ := client.New(&http.Client{}, testServer.URL)

	_, err = unitUnderTest.CreateAccount(&account)

	if err != nil {
		t.Errorf("Test request failed with error: '%v'", err.Error())
	}
}

func TestCreateAccountFailWithMock(t *testing.T) {

	httpResponseStatus := []int{http.StatusInternalServerError, http.StatusGatewayTimeout, http.StatusBadGateway}

	jsonFile, err := os.Open("createaccount.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	var account client.Account
	json.Unmarshal(bytes, &account)

	for _, h := range httpResponseStatus {

		t.Run(http.StatusText(h), func(t *testing.T) {

			mux := http.NewServeMux()

			testServer := httptest.NewServer(mux)

			mux.HandleFunc("/v1/organisation/accounts/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(h)
			})

			defer testServer.Close()

			unitUnderTest, _ := client.New(&http.Client{}, testServer.URL)

			_, err = unitUnderTest.CreateAccount(&account)

			if err == nil {
				t.Errorf("Test failed - Expected error response")
			}
		})
	}
}

func TestGetAccountWithMock(t *testing.T) {

	byteValue, err := readTestData("fetchresponse.json")
	if err != nil {
		t.Errorf("Test request failed with error: '%v'", err.Error())
	}

	mux := http.NewServeMux()

	testServer := httptest.NewServer(mux)

	mux.HandleFunc("/v1/organisation/accounts/41426819", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(byteValue)
	})

	defer testServer.Close()

	unitUnderTest, _ := client.New(&http.Client{}, testServer.URL)

	account, err := unitUnderTest.GetAccount("41426819")

	if err != nil {
		t.Errorf("Test request failed with error: '%v'", err.Error())
	}

	if account.Attributes.AccountNumber != "41426819" {
		t.Error("AccountNumber was not mapped to response.")
	}
}

func TestDeleteAccountWithMock(t *testing.T) {

	expectedVersion := "0"

	mux := http.NewServeMux()

	testServer := httptest.NewServer(mux)

	mux.HandleFunc("/v1/organisation/accounts/41426819", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		versionQuery := r.URL.Query().Get("version")
		if versionQuery != expectedVersion {
			t.Errorf("Request made with incorrect version query parameter. Expected '%v' but got '%v'", expectedVersion, versionQuery)
		}
	})

	defer testServer.Close()

	unitUnderTest, _ := client.New(&http.Client{}, testServer.URL)

	err := unitUnderTest.DeleteAccount("41426819", expectedVersion)

	if err != nil {
		t.Errorf("Test request failed with error: '%v'", err.Error())
	}
}

func readTestData(filename string) ([]byte, error) {

	jsonFile, err := os.Open(filename)
	if err != nil {
		return []byte{}, err
	}

	defer jsonFile.Close()

	return ioutil.ReadAll(jsonFile)
}
