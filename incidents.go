package gorsa

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//------------------------------------------------------------------------------
// Structures
//------------------------------------------------------------------------------

// Incident is a RSA security incident
type Incident struct {
	ID                        string    `json:"id"`
	Title                     string    `json:"title"`
	Summary                   string    `json:"summary"`
	Priority                  string    `json:"priority"`
	RiskScore                 int       `json:"riskScore"`
	Status                    string    `json:"status"`
	AlertCount                int       `json:"alertCount"`
	AverageAlertRiskScore     int       `json:"averageAlertRiskScore"`
	Sealed                    bool      `json:"sealed"`
	TotalRemediationTaskCount int       `json:"totalRemediationTaskCount"`
	OpenRemediationTaskCount  int       `json:"openRemediationTaskCount"`
	Created                   time.Time `json:"created"`
	LastUpdated               time.Time `json:"lastUpdated"`
	LastUpdatedBy             string    `json:"lastUpdatedBy"`
	Assignee                  string    `json:"assignee"`
	Sources                   []string  `json:"sources"`
	RuleID                    string    `json:"ruleId"`
	FirstAlertTime            time.Time `json:"firstAlertTime"`
	Categories                []struct {
		ID     string `json:"id"`
		Parent string `json:"parent"`
		Name   string `json:"name"`
	} `json:"categories"`
	JournalEntries []struct {
		ID          string    `json:"id"`
		Author      string    `json:"author"`
		Notes       string    `json:"notes"`
		Created     time.Time `json:"created"`
		LastUpdated time.Time `json:"lastUpdated"`
		Milestone   string    `json:"milestone"`
	} `json:"journalEntries"`
	CreatedBy         string `json:"createdBy"`
	DeletedAlertCount int    `json:"deletedAlertCount"`
	EventCount        int    `json:"eventCount"`
	AlertMeta         struct {
		SourceIP      []string `json:"SourceIp"`
		DestinationIP []string `json:"DestinationIp"`
	} `json:"alertMeta"`
}

// IncidentsResponse ...
type IncidentsResponse struct {
	Items       []*Incident `json:"items"`
	PageNumber  int         `json:"pageNumber"`
	PageSize    int         `json:"pageSize"`
	TotalPages  int         `json:"totalPages"`
	TotalItems  int         `json:"totalItems"`
	HasNext     bool        `json:"hasNext"`
	HasPrevious bool        `json:"hasPrevious"`
}

// ErrorResponse ...
type ErrorResponse struct {
	Status    int       `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Errors    []struct {
		Message string `json:"message"`
		Field   string `json:"field,omitempty"`
	} `json:"errors"`
}

//------------------------------------------------------------------------------
// Functions
//------------------------------------------------------------------------------

// ListIncidents returns the incidents with given cirterias.
func (endpoint *Endpoint) ListIncidents(ctx context.Context, since, until string, perPage, pageNumber int) (*IncidentsResponse, error) {
	// Authenticate and get the token
	token, err := endpoint.client.Authenticate(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while authenticating : %s", err)
	}

	// Prepare the URL
	reqURL, err := url.Parse(endpoint.client.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing the URL : %s", err)
	}
	reqURL.Path += "/rest/api/incidents"
	parameters := url.Values{}
	parameters.Add("since", since)
	parameters.Add("pageSize", strconv.Itoa(perPage))
	parameters.Add("pageNumber", strconv.Itoa(pageNumber))
	reqURL.RawQuery = parameters.Encode()

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error while creating the request : %s", err)
	}

	// Set HTTP headers
	req.Header.Set("NetWitness-Token", token)
	req.Header.Set("NetWitness-Version", endpoint.client.Version)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	resp, err := endpoint.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while doing the request: %s", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading the request: %s", err)
	}

	// Unmarshal the response
	var incidentsResponse *IncidentsResponse
	err = json.Unmarshal([]byte(body), &incidentsResponse)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling the response : %s", err)
	}

	return incidentsResponse, nil
}
