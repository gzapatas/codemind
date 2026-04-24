package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Ping checks availability by requesting a model list.
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/models", nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// Embeddings generates embeddings for the given text. (stub)
func (c *Client) Embeddings(ctx context.Context, text string) ([]float64, error) {
	return nil, fmt.Errorf("Embeddings not implemented yet")
}

// ListModels returns a list of available model names from the Ollama server.
func (c *Client) ListModels(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/models", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var raw any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	out := make([]string, 0)

	switch v := raw.(type) {
	case []any:
		for _, it := range v {
			if m, ok := it.(map[string]any); ok {
				if name := extractName(m); name != "" {
					out = append(out, name)
					continue
				}
				if b, err := json.Marshal(m); err == nil {
					out = append(out, string(b))
					continue
				}
			}
		}
	case map[string]any:
		// common shapes:
		// { "models": [...] }
		// { "data": [...] }
		// { "object": "list", "data": [...] }
		if modelsRaw, ok := v["models"]; ok {
			if arr, ok := modelsRaw.([]any); ok {
				for _, it := range arr {
					if m, ok := it.(map[string]any); ok {
						if name := extractName(m); name != "" {
							out = append(out, name)
							continue
						}
						if b, err := json.Marshal(m); err == nil {
							out = append(out, string(b))
							continue
						}
					}
				}
			}
		} else if dataRaw, ok := v["data"]; ok {
			if arr, ok := dataRaw.([]any); ok {
				for _, it := range arr {
					if m, ok := it.(map[string]any); ok {
						if name := extractName(m); name != "" {
							out = append(out, name)
							continue
						}
						if b, err := json.Marshal(m); err == nil {
							out = append(out, string(b))
							continue
						}
					}
				}
			} else if m, ok := dataRaw.(map[string]any); ok {
				// data is a map, maybe contains models inside
				if inner, ok := m["models"]; ok {
					if arr, ok := inner.([]any); ok {
						for _, it := range arr {
							if mm, ok := it.(map[string]any); ok {
								if name := extractName(mm); name != "" {
									out = append(out, name)
									continue
								}
								if b, err := json.Marshal(mm); err == nil {
									out = append(out, string(b))
									continue
								}
							}
						}
					}
				} else {
					// try to extract name from the map itself
					if name := extractName(m); name != "" {
						out = append(out, name)
					} else {
						for k := range m {
							out = append(out, k)
						}
					}
				}
			}
		} else {
			// fallback: maybe the map keys are model names
			for k := range v {
				out = append(out, k)
			}
		}
	default:
		// fallback: stringify
		if b, err := json.Marshal(raw); err == nil {
			out = append(out, string(b))
		}
	}
	return out, nil
}

func extractName(m map[string]any) string {
	if s, ok := m["name"].(string); ok && s != "" {
		return s
	}
	if s, ok := m["model"].(string); ok && s != "" {
		return s
	}
	if s, ok := m["id"].(string); ok && s != "" {
		return s
	}
	return ""
}
