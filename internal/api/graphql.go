package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cemilsahin/arabamtaksit/internal/config"
	"github.com/gorilla/websocket"
)

type GraphQLService struct {
	url     string
	wsURL   string
	query   string
	model   any
	headers map[string]string
	vars    map[string]any
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

type GraphQLError struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type GraphQLResponse struct {
	Data   any            `json:"data"`
	Errors []GraphQLError `json:"errors"`
}

// New creates a new GraphQLService instance
func NewGql(model any) *GraphQLService {
	apiURL := os.Getenv("GQL_URL")
	wsURL := strings.Replace(apiURL, "http", "ws", 1)

	return &GraphQLService{
		url:     apiURL,
		wsURL:   wsURL,
		headers: make(map[string]string),
		vars:    make(map[string]any),
		model:   model,
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
		g.vars[k] = v
	}
	return g
}

// Subscribe creates a WebSocket connection and subscribes to a GraphQL subscription
func (g *GraphQLService) Subscribe(subscription string, callback func(model any)) error {
	if g.conn == nil {
		if g.headers["Authorization"] == "" {
			g.headers["Authorization"] = "Bearer " + config.App().Token
		}

		dialer := websocket.Dialer{
			EnableCompression: true,
			HandshakeTimeout:  time.Second * 10,
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
		Variables: g.vars,
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
				if err := json.Unmarshal(msg.Payload, &g.model); err != nil {
					continue
				}
				callback(g.model)
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
func (g *GraphQLService) Query(query string) error {
	g.query = query
	return g.execute()
}

// Mutation executes a GraphQL mutation
func (g *GraphQLService) Mutation(mutation string) error {
	g.query = mutation
	return g.execute()
}

func (g *GraphQLService) execute() error {
	reqBody := graphQLRequest{
		Query:     g.query,
		Variables: g.vars,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", err)
	}

	req, err := http.NewRequest("POST", g.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("create request error: %v", err)
	}

	if g.headers["Authorization"] == "" {
		g.headers["Authorization"] = "Bearer " + config.App().Token
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
		return fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response error: %v", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()

	var result map[string]any
	if err := decoder.Decode(&result); err != nil {
		return fmt.Errorf("json decode error: %v", err)
	}

	// Convert the decoded data back to JSON with proper time format
	formattedJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", err)
	}

	if err := json.Unmarshal(formattedJSON, &g.model); err != nil {
		return fmt.Errorf("json unmarshal error: %v", err)
	}

	g.reset()
	return nil
}

func (g *GraphQLService) reset() {
	g.query = ""
	g.model = nil
	g.vars = make(map[string]any)
	g.headers = make(map[string]string)
}

// QueryWithContext executes a GraphQL query with context
func (g *GraphQLService) QueryWithContext(ctx context.Context, query string) error {
	// context kullanımı
	return nil
}

// EXAMPLE
/*

	type FacilityResponse struct {
		Data struct {
			Facilities []model.Facility `json:"facilities"`
		} `json:"data"`
	}

	var facilityResponse FacilityResponse
	gql := api.NewGql(&facilityResponse)

	query := `
		query MyQuery {
			facilities {
				id
				facility_documents {
					id
					document {
						id
						name
						document
						created_at
						updated_at
					}
				}
				facility_translations(where: {language: {code: {_eq: "en"}}}) {
					id
					title
					description
				}
			}
		}
	`

	err = gql.Query(query)
	if err != nil {
		log.Printf("GraphQL Query Error: %v", err)
		return err
	}

	log.Println(facilityResponse.Data.Facilities[0].FacilityTranslations[0].Title)



*/
