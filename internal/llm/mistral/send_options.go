package mistral

import (
	"github.com/bornholm/sidequest/internal/llm"
)

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

type Tool struct {
	Type     ToolType     `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Description string             `json:"description"`
	Name        string             `json:"name"`
	Parameters  FunctionParameters `json:"parameters"`
}

type FunctionParameters map[string]any

type ToolChoice string

const (
	ToolChoiceAuto ToolChoice = "auto"
	ToolChoiceNone ToolChoice = "none"
	ToolChoiceAny  ToolChoice = "any"
)

func WithTools(choice ToolChoice, tools ...Tool) llm.SendOptionFunc {
	return llm.WithAttrs(map[string]any{
		"tool_choice": string(choice),
		"tools":       tools,
	})
}

func WithJSONMode() llm.SendOptionFunc {
	return llm.WithAttr("response_format", struct {
		Type string `json:"type"`
	}{
		Type: "json_object",
	})
}
