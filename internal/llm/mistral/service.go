package mistral

import (
	"io"
	"net/http"
	"strings"

	"github.com/bornholm/sidequest/internal/llm"
	"github.com/pkg/errors"
)

type Service struct {
	baseURL string
	apiKey  string
	model   string
}

func NewService(baseURL string, apiKey string, model string) *Service {
	return &Service{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		apiKey:  apiKey,
		model:   model,
	}
}

func (s *Service) apiPost(path string, body io.Reader) (*http.Response, error) {
	return s.apiDo("POST", path, body)
}

func (s *Service) apiDo(method string, path string, body io.Reader) (*http.Response, error) {
	url := s.baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json; charset=utf-8")

	if s.apiKey != "" {
		req.Header.Add("Authorization", "Bearer "+s.apiKey)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

var _ llm.Service = &Service{}
