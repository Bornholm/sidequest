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
	generateQuestTool = mistral.Tool{
		Type: mistral.ToolTypeFunction,
		Function: mistral.ToolFunction{
			Name:        "createQuest",
			Description: "Create a new quest plot",
			Parameters: mistral.FunctionParameters{
				"type": "object",
				"properties": map[string]map[string]any{
					"title": {
						"type":        "string",
						"description": "The title of the quest",
					},
					"description": {
						"type":        "string",
						"description": "The description of the quest plot",
					},
					"characters": {
						"type":        "string",
						"description": "The description of the characters involved in the plot, with their roles in it",
					},
					"solution": {
						"type":        "string",
						"description": "The general description of the solution to the quest plot",
					},
					"clues": {
						"type":        "string",
						"description": "A list of the clues that the players can find to discover the solution to the plot. A complete description is provided for each clue, and the elements each brings to the plot.",
					},
				},
				"required": []string{"title", "description", "characters", "solution", "clues"},
			},
		},
	}
)

type Quest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Characters  string `json:"characters"`
	Solution    string `json:"solution"`
	Clues       string `json:"clues"`
}

type GenerateQuestPromptData struct {
	Language string          `json:"archetype"`
	Quest    QuestContext    `json:"quest"`
	Universe UniverseContext `json:"universe"`
}

type QuestContext struct {
	Quest `json:",inline"`
}

func GenerateQuest(c echo.Context) error {
	ctx := c.Request().Context()

	promptData := GenerateQuestPromptData{
		Language: "FR-fr",
		Quest: QuestContext{
			Quest: Quest{},
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

	questPrompt, err := generatePrompt(prompt.Quest, promptData)
	if err != nil {
		return errors.WithStack(err)
	}

	spew.Dump(questPrompt)

	response, err := sess.Send(
		ctx, llm.NewMessage("user", questPrompt, nil),
		llm.WithTemperature(0.8),
		mistral.WithTools(
			mistral.ToolChoiceAny,
			generateQuestTool,
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

	if toolCall.Function.Name != generateQuestTool.Function.Name {
		return errors.Errorf("unexpected tool call '%s'", toolCall.Function.Name)
	}

	var quest Quest

	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &quest); err != nil {
		return errors.Wrapf(err, "could not unmarshal executed tool call arguments")
	}

	return c.JSON(http.StatusOK, struct {
		Quest Quest `json:"quest"`
	}{
		Quest: quest,
	})
}
