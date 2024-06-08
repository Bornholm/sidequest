package mistral

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/bornholm/sidequest/internal/llm"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Chat implements llm.Service.
func (s *Service) Chat(ctx context.Context, messages ...llm.Message) (llm.ChatSession, error) {
	return &ChatSession{
		srv:      s,
		model:    s.model,
		messages: messages,
	}, nil
}

type ChatSession struct {
	srv      *Service
	model    string
	messages []llm.Message
}

type ChatRequest struct {
	Model       string        `mapstructure:"model"`
	Messages    []llm.Message `mapstructure:"messages"`
	Temperature float64       `mapstructure:"temperature,omitempty"`
}

//{"created":1708165814,"object":"chat.completion","id":"a64f9837-f819-4133-9172-bc907f89486e","model":"mistral-openorca","choices":[{"index":0,"finish_reason":"stop","message":{"role":"assistant","content":" I'm doing well, thank you for asking! How about you?\n"}}],"usage":{"prompt_tokens":0,"completion_tokens":0,"total_tokens":0}}

type ChatResponse struct {
	Created int    `json:"created"`
	Object  string `json:"object"`
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int      `json:"index"`
		FinishReason string   `json:"finish_reason"`
		Message      *Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type Message struct {
	Role      string     `json:"role" `
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

func (m *Message) Interface() llm.Message {
	return llm.NewMessage(m.Role, m.Content, m)
}

func (r *ChatResponse) Message() llm.Message {
	lastChoice := r.Choices[len(r.Choices)-1]
	return llm.NewMessage(lastChoice.Message.Role, lastChoice.Message.Content, lastChoice.Message)
}

type ToolCall struct {
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string  `json:"name"`
	Arguments string  `json:"arguments"`
	ID        *string `json:"id"`
}

// Send implements llm.ChatSession.
func (s *ChatSession) Send(ctx context.Context, message llm.Message, funcs ...llm.SendOptionFunc) (llm.ChatResponse, error) {
	opts := llm.NewSendOptions(funcs...)

	messages := make([]llm.Message, len(s.messages))
	copy(messages, s.messages)
	messages = append(messages, message)

	chatReq := ChatRequest{
		Model:       s.model,
		Temperature: opts.Temperature,
		Messages:    messages,
	}

	var payload map[string]any
	if err := mapstructure.Decode(chatReq, &payload); err != nil {
		return nil, errors.WithStack(err)
	}

	for key, value := range opts.Attrs {
		payload[key] = value
	}

	var body bytes.Buffer
	encoder := json.NewEncoder(&body)

	if err := encoder.Encode(payload); err != nil {
		return nil, errors.WithStack(err)
	}

	res, err := s.srv.apiPost("/v1/chat/completions", &body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer res.Body.Close()

	rawResBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var chatRes ChatResponse

	if err := json.Unmarshal(rawResBody, &chatRes); err != nil {
		return nil, errors.WithStack(err)
	}

	if chatRes.Error != nil {
		return nil, errors.Errorf("server error: %d - %s", chatRes.Error.Code, chatRes.Error.Message)
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("server error: %d - %s: %s", res.StatusCode, res.Status, rawResBody)
	}

	last := chatRes.Choices[0].Message

	s.messages = append(s.messages, last.Interface())

	return &chatRes, nil
}

var _ llm.ChatSession = &ChatSession{}
