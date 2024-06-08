package route

import (
	"encoding/json"
	"net/http"

	"github.com/bornholm/sidequest/internal/llm"
	"github.com/bornholm/sidequest/internal/llm/mistral"
	"github.com/bornholm/sidequest/internal/prompt"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v5"
	"github.com/pkg/errors"
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
						"description": "The life objectives of the character. Should not be empty",
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
	Language  string           `json:"archetype"`
	Character CharacterContext `json:"character"`
	Universe  UniverseContext  `json:"universe"`
}

type CharacterContext struct {
	Character `json:",inline"`
	Alignment string `json:"alignment"`
	Archetype string `json:"archetype"`
}

func GenerateCharacter(c echo.Context) error {
	ctx := c.Request().Context()

	promptData := GenerateCharacterPromptData{
		Language: "FR-fr",
		Character: CharacterContext{
			Character: Character{},
		},
	}

	if err := c.Bind(&promptData); err != nil {
		return errors.WithStack(err)
	}

	srv := mistral.NewService(mistralBaseURL, mistralAPIKey, mistralChatModel)

	sess, err := srv.Chat(ctx, llm.NewMessage("system", prompt.Agent, nil))
	if err != nil {
		return errors.WithStack(err)
	}

	characterPrompt, err := generatePrompt(prompt.Character, promptData)
	if err != nil {
		return errors.WithStack(err)
	}

	spew.Dump(characterPrompt)

	response, err := sess.Send(
		ctx, llm.NewMessage("user", characterPrompt, nil),
		llm.WithTemperature(0.8),
		mistral.WithTools(
			mistral.ToolChoiceAny,
			generateCharacterTool,
		),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	spew.Dump(response)

	mistralResponse, ok := response.(*mistral.ChatResponse)
	if !ok {
		return errors.Errorf("unexpected llm response type '%T'", response)
	}

	if len(mistralResponse.Choices) == 0 {
		return errors.New("no choice available")
	}

	choice := mistralResponse.Choices[0]

	if len(choice.Message.ToolCalls) == 0 {
		return errors.New("no tool call available")
	}

	toolCall := choice.Message.ToolCalls[0]

	if toolCall.Function.Name != generateCharacterTool.Function.Name {
		return errors.Errorf("unexpected tool call '%s'", toolCall.Function.Name)
	}

	var character Character

	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &character); err != nil {
		return errors.Wrapf(err, "could not unmarshal executed tool call arguments")
	}

	return c.JSON(http.StatusOK, struct {
		Character Character `json:"character"`
	}{
		Character: character,
	})
}
