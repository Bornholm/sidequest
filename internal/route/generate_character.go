package route

import (
	"github.com/bornholm/sidequest/internal/llm/mistral"
	"github.com/bornholm/sidequest/internal/prompt"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
)

var (
	generateCharacterTool = mistral.Tool{
		Type: mistral.ToolTypeFunction,
		Function: mistral.ToolFunction{
			Name:        "createCharacter",
			Description: "Create a new character for the role playing game",
			Parameters: mistral.FunctionParameters{
				"type": "object",
				"properties": map[string]map[string]any{
					"name": {
						"type":        "string",
						"description": "The full name of the character",
					},
					"story": {
						"type":        "string",
						"description": "The origin story of the character. Must be at least 2 paragraphs.",
					},
					"race": {
						"type":        "string",
						"description": "The race of the character",
					},
					"sex": {
						"type":        "string",
						"description": "The sex of the character",
					},
					"age": {
						"type":        "number",
						"description": "The age of the character, in years",
					},
					"objectives": {
						"type":        "string",
						"description": "The life objectives of the character. Non empty.",
					},
				},
				"required": []string{"name", "story", "sex", "age", "race", "objectives"},
			},
		},
	}
)

type Character struct {
	Name       string `json:"name"`
	Story      string `json:"story"`
	Sex        string `json:"sex"`
	Race       string `json:"race"`
	Objectives string `json:"objectives"`
	Age        int    `json:"age"`
}

type GenerateCharacterPromptData struct {
	Language  string           `json:"language"`
	Character CharacterContext `json:"character"`
	Universe  UniverseContext  `json:"universe"`
}

type CharacterContext struct {
	Character `json:",inline"`
	Alignment string `json:"alignment"`
	Archetype string `json:"archetype"`
}

func GenerateCharacter(app *pocketbase.PocketBase) echo.HandlerFunc {
	return createGenerateHandler[GenerateCharacterPromptData](
		app, prompt.Character,
		generateCharacterTool,
		func(character Character) any {
			return struct {
				Character Character `json:"character"`
			}{
				Character: character,
			}
		},
	)
}
