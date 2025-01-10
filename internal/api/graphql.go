package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/mstgnz/starter-kit/internal/response"
)

type GraphQLService struct {
	url     string
	wsURL   string
	query   string
	headers map[string]string
	params  map[string]any
	conn    *websocket.Conn
}

type graphQLRequest struct {
	ID        string         `json:"id,omitempty"`
	Type      string         `json:"type,omitempty"`
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type wsMessage struct {
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// New creates a new GraphQLService instance
func NewGql() *GraphQLService {
	apiURL := os.Getenv("GQL_URL")
	wsURL := strings.Replace(apiURL, "http", "ws", 1)

	return &GraphQLService{
		url:     apiURL,
		wsURL:   wsURL,
		headers: make(map[string]string),
		params:  make(map[string]any),
	}
}

// WithHeader adds headers to the request
func (g *GraphQLService) WithHeader(headers map[string]string) *GraphQLService {
	for k, v := range headers {
		g.headers[k] = v
	}
	return g
}

// WithToken adds authorization token to headers
func (g *GraphQLService) WithToken(token string) *GraphQLService {
	g.headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	return g
}

// WithVariables adds variables to the GraphQL query
func (g *GraphQLService) WithVariables(vars map[string]any) *GraphQLService {
	for k, v := range vars {
		g.params[k] = v
	}
	return g
}

// Subscribe creates a WebSocket connection and subscribes to a GraphQL subscription
func (g *GraphQLService) Subscribe(subscription string, callback func(response *response.Response)) error {
	if g.conn == nil {
		dialer := websocket.Dialer{
			EnableCompression: true,
		}

		header := http.Header{}
		for k, v := range g.headers {
			header.Set(k, v)
		}

		conn, _, err := dialer.Dial(g.wsURL, header)
		if err != nil {
			return fmt.Errorf("websocket connection error: %v", err)
		}
		g.conn = conn

		// Send connection init message
		initMsg := wsMessage{
			Type: "connection_init",
		}
		if err := g.conn.WriteJSON(initMsg); err != nil {
			return fmt.Errorf("websocket init error: %v", err)
		}
	}

	// Send subscription request
	subscriptionMsg := graphQLRequest{
		ID:        "1", // You might want to generate unique IDs
		Type:      "start",
		Query:     subscription,
		Variables: g.params,
	}

	if err := g.conn.WriteJSON(subscriptionMsg); err != nil {
		return fmt.Errorf("subscription request error: %v", err)
	}

	// Start listening for messages
	go func() {
		for {
			var msg wsMessage
			err := g.conn.ReadJSON(&msg)
			if err != nil {
				// Handle connection errors
				if websocket.IsUnexpectedCloseError(err) {
					g.conn = nil
					return
				}
				continue
			}

			switch msg.Type {
			case "data":
				var resp response.Response
				if err := json.Unmarshal(msg.Payload, &resp); err != nil {
					continue
				}
				callback(&resp)
			case "complete":
				g.conn.Close()
				g.conn = nil
				return
			case "error":
				// Handle subscription errors
				continue
			}
		}
	}()

	return nil
}

// Unsubscribe closes the WebSocket connection
func (g *GraphQLService) Unsubscribe() error {
	if g.conn != nil {
		stopMsg := wsMessage{
			ID:   "1", // Should match the subscription ID
			Type: "stop",
		}
		if err := g.conn.WriteJSON(stopMsg); err != nil {
			return fmt.Errorf("unsubscribe error: %v", err)
		}

		err := g.conn.Close()
		g.conn = nil
		return err
	}
	return nil
}

// Query executes a GraphQL query
func (g *GraphQLService) Query(query string) (*response.Response, error) {
	g.query = query
	return g.execute()
}

// Mutation executes a GraphQL mutation
func (g *GraphQLService) Mutation(mutation string) (*response.Response, error) {
	g.query = mutation
	return g.execute()
}

func (g *GraphQLService) execute() (*response.Response, error) {
	reqBody := graphQLRequest{
		Query:     g.query,
		Variables: g.params,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %v", err)
	}

	req, err := http.NewRequest("POST", g.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request error: %v", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")

	// Set custom headers
	for k, v := range g.headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %v", err)
	}

	var result response.Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %v", err)
	}

	g.reset()
	return &result, nil
}

func (g *GraphQLService) reset() {
	g.query = ""
	g.params = make(map[string]any)
	g.headers = make(map[string]string)
}

// EXAMPLE
/*

type Town struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type District struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Towns []Town `json:"towns"`
}

type City struct {
    ID        string     `json:"id"`
    Name      string     `json:"name"`
    Districts []District `json:"districts"`
}

type Country struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
    Cities  []City `json:"cities"`
}

type LocationResponse struct {
    Countries []Country `json:"countries"`
}

graphqlService := api.NewGql()

query := `
    query MyQuery($_eq: smallint) {
        countries {
            id
            name
            cities(where: {id: {_eq: $_eq}}) {
                id
                name
                districts {
                    id
                    name
                    towns {
                        id
                        name
                    }
                }
            }
        }
    }

	variables := map[string]any{
        "_eq": 34
    }

	response, err := graphqlService.
        WithVariables(variables).
        Query(query)

    if err != nil {
        log.Fatalf("Query error: %v", err)
    }

	var locationData LocationResponse
    if err := json.Unmarshal(response.Data, &locationData); err != nil {
        log.Fatalf("JSON unmarshal error: %v", err)
    }

*/
