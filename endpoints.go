package gorsa

import "context"

//------------------------------------------------------------------------------
// Structure
//------------------------------------------------------------------------------

// Endpoint is an API endpoint.
type Endpoint struct {
	client *Client
}

//------------------------------------------------------------------------------
// Interfaces
//------------------------------------------------------------------------------

// Incidents endpoint.
type Incidents interface {
	ListIncidents(context.Context, string, string, int, int) (*IncidentsResponse, error)
}
