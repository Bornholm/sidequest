package route

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bornholm/sidequest/internal/llm"
	"github.com/bornholm/sidequest/internal/llm/mistral"
	"github.com/bornholm/sidequest/internal/prompt"
	"github.com/labstack/echo/v5"
	"github.com/pkg/errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
)

type UniverseContext struct {
	Style string `json:"style"`
}

func createGenerateHandler[P any, E any](app *pocketbase.PocketBase, promptTemplate string, tool mistral.Tool, createResponse func(e E) any) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		user, ok := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		if !ok || user == nil {
			return errors.New("could not retrieve current user")
		}

		userStats, err := app.Dao().FindFirstRecordByData("user_stats", "user", user.Id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.WithStack(err)
		}

		if userStats == nil {
			collection, err := app.Dao().FindCollectionByNameOrId("user_stats")
			if err != nil {
				return errors.WithStack(err)
			}

			userStats = models.NewRecord(collection)
			userStats.Set("totalTokens", 0)
			userStats.Set("maxTokens", defaultMaxTokens)
			userStats.Set("user", user.Id)

			if err := app.Dao().SaveRecord(userStats); err != nil {
				return errors.WithStack(err)
			}
		}

		maxTokens := userStats.GetInt("maxTokens")
		if maxTokens == 0 {
			maxTokens = defaultMaxTokens
		}

		totalTokens := userStats.GetInt("totalTokens")
		if maxTokens != -1 && totalTokens >= maxTokens {
			return c.JSON(http.StatusPaymentRequired, struct {
				Message     string `json:"message"`
				TotalTokens int    `json:"totalTokens"`
				MaxTokens   int    `json:"maxTokens"`
			}{
				Message:     "tokens quota exceeded",
				TotalTokens: totalTokens,
				MaxTokens:   maxTokens,
			})
		}

		var promptData P

		if err := c.Bind(&promptData); err != nil {
			return errors.WithStack(err)
		}

		srv := mistral.NewService(mistralBaseURL, mistralAPIKey, mistralChatModel)

		sess, err := srv.Chat(ctx, llm.NewMessage("system", prompt.Agent, nil))
		if err != nil {
			return errors.WithStack(err)
		}

		prompt, err := generatePrompt(promptTemplate, promptData)
		if err != nil {
			return errors.WithStack(err)
		}

		response, err := sess.Send(
			ctx, llm.NewMessage("user", prompt, nil),
			llm.WithTemperature(0.8),
			mistral.WithTools(
				mistral.ToolChoiceAny,
				tool,
			),
		)
		if err != nil {
			return errors.WithStack(err)
		}

		mistralResponse, ok := response.(*mistral.ChatResponse)
		if !ok {
			return errors.Errorf("unexpected llm response type '%T'", response)
		}

		_, err = app.Dao().DB().
			NewQuery("UPDATE user_stats SET totalTokens = totalTokens + {:tokens} WHERE id = {:id} ").
			Bind(dbx.Params{
				"tokens": mistralResponse.Usage.TotalTokens,
				"id":     userStats.Id,
			}).Execute()
		if err != nil {
			return errors.WithStack(err)
		}

		if len(mistralResponse.Choices) == 0 {
			return errors.New("no choice available")
		}

		choice := mistralResponse.Choices[0]

		if len(choice.Message.ToolCalls) == 0 {
			return errors.New("no tool call available")
		}

		toolCall := choice.Message.ToolCalls[0]

		if toolCall.Function.Name != tool.Function.Name {
			return errors.Errorf("unexpected tool call '%s'", toolCall.Function.Name)
		}

		var entity E

		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &entity); err != nil {
			return errors.Wrapf(err, "could not unmarshal executed tool call arguments")
		}

		return c.JSON(http.StatusOK, createResponse(entity))
	}
}
