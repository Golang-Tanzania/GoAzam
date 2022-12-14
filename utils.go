package GoAzam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type APICONTEXT struct {
	AppName      string
	ClientID     string
	ClientSecret string
	Token        string
	BaseURL      string
	Bearer       string
}

type kEYS struct {
	AppName      string
	ClientId     string
	ClientSecret string
	Token        string
}

func (api *APICONTEXT) generateSessionID() string {
	api.BaseURL = "https://sandbox.azampay.co.tz"
	endpoint := "/AppRegistration/GenerateToken"
	url := fmt.Sprintf("https://authenticator-sandbox.azampay.co.tz%v", endpoint)
	parameters := fmt.Sprintf(`{"appName":"%v", "clientId": "%v", "clientSecret": "%v"}`, api.AppName, api.ClientID, api.ClientSecret)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(parameters)))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err.Error()
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err.Error()
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err.Error()
	}

	type Result struct {
		Data       map[string]string
		Message    string
		Success    bool
		StatusCode int
	}
	var result Result
	json.Unmarshal([]byte(body), &result)

	if result.Data["accessToken"] != "" {
		api.Bearer = result.Data["accessToken"]
		return api.Bearer
	} else {
		return string(body)
	}
}

func (api *APICONTEXT) LoadKeys(file string) *APICONTEXT {

	keys, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println(err)
	}

	var Keys kEYS

	err = json.Unmarshal(keys, &Keys)

	if err != nil {
		fmt.Println(err)
	}

	api.AppName = Keys.AppName
	api.ClientID = Keys.ClientId
	api.ClientSecret = Keys.ClientSecret
	api.Token = Keys.Token

	return api
}

func (api *APICONTEXT) sendRequest(endpoint string, query map[string]string) string {
	jsonParameters, err := json.Marshal(query)
	if err != nil {
		fmt.Println(err)
	}
	url := api.BaseURL + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonParameters)))
	req.Header.Set("Content-Type", "application/json")
	bearer := fmt.Sprintf("Bearer %v", api.Bearer)
	req.Header.Set("Authorization", bearer)
	req.Header.Set("X-API-KEY", api.Token)
	if err != nil {
		return err.Error()
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err.Error()
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err.Error()
	}

	return string(body)
}
