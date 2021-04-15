package gorsa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultVersion = "1.0"
)

//------------------------------------------------------------------------------
// Structure
//------------------------------------------------------------------------------

// Client is a client for QRadar REST API.
type Client struct {
	httpClient *http.Client

	BaseURL  string
	Username string
	Password string
	Version  string

	// Endpoints
	Incidents Incidents
}

//------------------------------------------------------------------------------
// Factory
//------------------------------------------------------------------------------

// NewClient returns a new QRadar API client.
func NewClient(httpClient *http.Client, baseURL, username, password string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// Create the client
	c := &Client{
		httpClient: httpClient,
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		Version:    defaultVersion,
	}

	// Add the endpoints
	c.Incidents = &Endpoint{client: c}

	return c
}

//------------------------------------------------------------------------------
// Factory
//------------------------------------------------------------------------------

// Authenticate to the REST API.
func (c *Client) Authenticate(ctx context.Context) (string, error) {
	// Prepare the URL
	var reqURL *url.URL
	reqURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("Error while parsing the URL : %s", err)
	}
	reqURL.Path += "/rest/api/auth/userpass"

	// Set the form values
	form := url.Values{
		"username": {c.Username},
		"password": {c.Password},
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL.String(), bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("Error while creating the request : %s", err)
	}

	// Set the headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	// Do the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error while doing the request : %s", err)
	}
	defer resp.Body.Close()

	// Read the respsonse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error while reading the request : %s", err)
	}

	// StatusCode is 500
	if resp.StatusCode == 500 {
		return "", fmt.Errorf("Status code is 500 and body is %s", string(body))
	}

	// Check the other status codes
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Status code is %d", resp.StatusCode)
	}

	// Prepare the response
	var tokenResponse *TokenResponse

	// Unmarshal the response
	err = json.Unmarshal([]byte(body), &tokenResponse)
	if err != nil {
		return "", fmt.Errorf("Error while unmarshalling the response : %s", err)
	}

	return tokenResponse.AccessToken, nil
}
