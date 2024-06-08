package route

import (
	"math"

	"github.com/bornholm/sidequest/internal/env"
)

var (
	mistralBaseURL   string
	mistralAPIKey    string
	mistralChatModel string
	defaultMaxTokens int
)

const (
	pricePerToken        = 0.0000113 // In euro per token
	defaultUserFreeQuota = 0.25      // In euros
)

func init() {
	mistralBaseURL = env.String("SIDEQUEST_MISTRAL_BASE_URL", "https://api.mistral.ai")
	mistralAPIKey = env.String("SIDEQUEST_MISTRAL_API_KEY", "")
	mistralChatModel = env.String("SIDEQUEST_MISTRAL_CHAT_MODEL", "mistral-large-latest")
	defaultMaxTokens = env.Int("SIDEQUEST_DEFAULT_MAX_TOKENS", int(math.Ceil(defaultUserFreeQuota/pricePerToken)))
}
