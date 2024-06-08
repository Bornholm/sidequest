package route

import "github.com/bornholm/sidequest/internal/env"

var (
	mistralBaseURL   string
	mistralAPIKey    string
	mistralChatModel string
)

func init() {
	mistralBaseURL = env.String("SIDEQUEST_MISTRAL_BASE_URL", "https://api.mistral.ai")
	mistralAPIKey = env.String("SIDEQUEST_MISTRAL_API_KEY", "")
	mistralChatModel = env.String("SIDEQUEST_MISTRAL_CHAT_MODEL", "mistral-large-latest")
}
