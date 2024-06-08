package mistral

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type embeddingsRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embeddingsResponse struct {
	ID      string `json:"id"`
	Created int    `json:"created"`
	Object  string `json:"object"`
	Model   string `json:"model"`
	Data    []struct {
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// Embeddings implements llm.Service.
func (s *Service) Embeddings(ctx context.Context, text string) ([]float64, error) {
	req := embeddingsRequest{
		Model: s.model,
		Input: text,
	}

	var body bytes.Buffer
	encoder := json.NewEncoder(&body)

	if err := encoder.Encode(req); err != nil {
		return nil, errors.WithStack(err)
	}

	res, err := s.apiPost("/v1/embeddings", &body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer res.Body.Close()

	rawResBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var result embeddingsResponse

	if err := json.Unmarshal(rawResBody, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	if result.Error != nil {
		return nil, errors.Errorf("server error: %d - %s", result.Error.Code, result.Error.Message)
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("server error: %d - %s: %s", res.StatusCode, res.Status, rawResBody)
	}

	if len(result.Data) < 1 {
		return nil, errors.New("unexpected number of results")
	}

	return result.Data[0].Embedding, nil
}
