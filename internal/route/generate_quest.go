package route

import (
	"github.com/bornholm/sidequest/internal/llm/mistral"
	"github.com/bornholm/sidequest/internal/prompt"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
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
	Language string          `json:"language"`
	Quest    QuestContext    `json:"quest"`
	Universe UniverseContext `json:"universe"`
}

type QuestContext struct {
	Quest `json:",inline"`
}

func GenerateQuest(app *pocketbase.PocketBase) echo.HandlerFunc {
	return createGenerateHandler[GenerateQuestPromptData](
		app, prompt.Quest,
		generateQuestTool,
		func(quest Quest) any {
			return struct {
				Quest Quest `json:"quest"`
			}{
				Quest: quest,
			}
		},
	)
}
