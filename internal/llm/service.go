package llm

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

type Service interface {
	Embeddings(ctx context.Context, text string) ([]float64, error)
	Chat(ctx context.Context, messages ...Message) (ChatSession, error)
}

type ChatSession interface {
	Send(ctx context.Context, message Message, funcs ...SendOptionFunc) (ChatResponse, error)
}

type ChatResponse interface {
	Message() Message
}

type SendOptions struct {
	Temperature float64
	MaxTokens   int
	Attrs       map[string]any
}

func NewSendOptions(funcs ...SendOptionFunc) *SendOptions {
	opts := &SendOptions{
		Temperature: 0.4,
		Attrs:       make(map[string]any),
	}
	for _, fn := range funcs {
		fn(opts)
	}
	return opts
}

type SendOptionFunc func(opts *SendOptions)

func WithTemperature(temperature float64) SendOptionFunc {
	return func(opts *SendOptions) {
		opts.Temperature = temperature
	}
}

func WithAttr(name string, value any) SendOptionFunc {
	return func(opts *SendOptions) {
		opts.Attrs[name] = value
	}
}

func WithAttrs(attrs map[string]any) SendOptionFunc {
	return func(opts *SendOptions) {
		for key, value := range attrs {
			opts.Attrs[key] = value
		}
	}
}

type Message interface {
	Role() string
	Content() string
	Source() any
}

type GenericMessage struct {
	role    string
	content string
	source  any
}

func (m *GenericMessage) Role() string {
	return m.role
}

func (m *GenericMessage) Content() string {
	return m.content
}

func (m *GenericMessage) Source() any {
	return m.source
}

func (m *GenericMessage) MarshalJSON() ([]byte, error) {
	var source any = m
	if src := m.Source(); src != nil {
		source = src
	} else {
		source = struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			Role:    m.role,
			Content: m.content,
		}
	}

	data, err := json.Marshal(source)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

func NewMessage(role string, content string, source any) *GenericMessage {
	return &GenericMessage{
		role:    role,
		content: content,
		source:  source,
	}
}
