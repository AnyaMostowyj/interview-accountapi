package client_test

import (
	"client"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

//const host = "http://accountapi:8080"
//const host = "http://localhost:8080"

var host = flag.String("host", "http://localhost:8080", "accountapi host address to execute tests against")

func TestCreateAccount(t *testing.T) {

	bytes := readTestData(t, "./test_data/createaccount.json")

	var account client.Account
	json.Unmarshal(bytes, &account)

	fmt.Println("Account type: " + account.Type)

	unitUnderTest, _ := client.New(&http.Client{}, *host)

	accountResponse, err := unitUnderTest.Create(&account)
	if err != nil {
		t.Error("Account creation error")
	}

	fmt.Println("Account response type: " + accountResponse.Type)

	tearDownTestData(t)
}

func TestCopCreateAccount(t *testing.T) {

	bytes := readTestData(t, "./test_data/cop_account_example.json")

	var account client.Account
	json.Unmarshal(bytes, &account)

	fmt.Println("Account type: " + account.Type)

	unitUnderTest, _ := client.New(&http.Client{}, *host)

	accountResponse, err := unitUnderTest.Create(&account)
	if err != nil {
		t.Error("Account creation error")
	}

	fmt.Println("Account response type: " + accountResponse.Type)

	tearDownTestData(t)
}

func TestCreateAccountFail(t *testing.T) {

	// Insert test account to accountapi
	setupTestData(t)

	bytes := readTestData(t, "./test_data/createaccount.json")

	var account client.Account
	json.Unmarshal(bytes, &account)

	fmt.Println("Account type: " + account.Type)

	unitUnderTest, _ := client.New(&http.Client{}, *host)

	// Attempt to create the same account again
	_, err := unitUnderTest.Create(&account)
	if err == nil {
		t.Error("test failed - expected error message for duplicate account creation")
	}

	tearDownTestData(t)
}

func TestGetAccount(t *testing.T) {

	// Insert test account to accountapi
	setupTestData(t)

	expectedAccountID := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

	unitUnderTest, _ := client.New(&http.Client{}, *host)

	account, err := unitUnderTest.Fetch(expectedAccountID)

	if err != nil {
		t.Errorf("Test request failed with error: '%v'", err.Error())
	}

	if account.Id != expectedAccountID {
		t.Errorf("Expected account ID: '%v' but got account ID: '%v'", expectedAccountID, account.Id)
	}

	tearDownTestData(t)
}

func TestDeleteAccount(t *testing.T) {

	setupTestData(t)

	expectedAccountID := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

	unitUnderTest, _ := client.New(&http.Client{}, *host)

	// Attempt to delete the account created in test setup
	err := unitUnderTest.Delete(expectedAccountID, "0")

	if err != nil {
		t.Errorf("Test request failed with error: '%v'", err.Error())
	}
}

func TestDeleteAccountFail(t *testing.T) {

	unitUnderTest, _ := client.New(&http.Client{}, *host)

	// Attempt to delete an account that doesn't exist
	err := unitUnderTest.Delete("invalidAccountId", "0")

	if err == nil {
		t.Error("test failed - expected error message for invalid deletion attempt")
	}
}

// These mocked test scenarios cover API responses that can't be triggered via normal API interactions
func TestCreateAccountWithMock(t *testing.T) {

	bytes := readTestData(t, "./test_data/createaccount.json")

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

	_, err := unitUnderTest.Create(&account)

	if err != nil {
		t.Errorf("Test request failed with error: '%v'", err.Error())
	}
}

func TestCreateCopAccountWithMock(t *testing.T) {

	bytes := readTestData(t, "./test_data/cop_request.json")

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

	_, err := unitUnderTest.Create(&account)

	if err != nil {
		t.Errorf("Test request failed with error: '%v'", err.Error())
	}
}

func TestCreateAccountFailWithMock(t *testing.T) {

	httpResponseStatus := []int{http.StatusInternalServerError, http.StatusGatewayTimeout, http.StatusBadGateway}

	bytes := readTestData(t, "./test_data/createaccount.json")

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

			_, err := unitUnderTest.Create(&account)

			if err == nil {
				t.Errorf("Test failed - Expected error response")
			}
		})
	}
}

func TestGetAccountWithMock(t *testing.T) {

	byteValue := readTestData(t, "./test_data/fetchresponse.json")

	mux := http.NewServeMux()

	testServer := httptest.NewServer(mux)

	mux.HandleFunc("/v1/organisation/accounts/41426819", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(byteValue)
	})

	defer testServer.Close()

	unitUnderTest, _ := client.New(&http.Client{}, testServer.URL)

	account, err := unitUnderTest.Fetch("41426819")

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

	err := unitUnderTest.Delete("41426819", expectedVersion)

	if err != nil {
		t.Errorf("Test request failed with error: '%v'", err.Error())
	}
}

func readTestData(t *testing.T, filename string) []byte {

	jsonFile, err := os.Open(filename)
	if err != nil {
		t.Fatalf("open of test setup data failed with error: '%v'", err.Error())
	}

	defer jsonFile.Close()

	fileBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("read of test setup data failed with error: '%v'", err.Error())
	}

	return fileBytes
}

func setupTestData(t *testing.T) {

	bytes := readTestData(t, "./test_data/createaccount.json")

	var account client.Account
	err := json.Unmarshal(bytes, &account)
	if err != nil {
		t.Fatalf("unmarshal of test setup data failed with error: '%v'", err.Error())
	}

	setupClient, _ := client.New(&http.Client{}, *host)

	_, err = setupClient.Create(&account)
	if err != nil {
		t.Fatalf("submission of test setup data failed with error: '%v'", err.Error())
	}
}

func tearDownTestData(t *testing.T) {

	testAccountID := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

	unitUnderTest, _ := client.New(&http.Client{}, *host)

	err := unitUnderTest.Delete(testAccountID, "0")
	if err != nil {
		t.Log("test tear down failed")
	}
}
