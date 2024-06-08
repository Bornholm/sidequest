package prompt

import (
	_ "embed"
)

//go:embed agent_prompt.txt
var Agent string

//go:embed character_prompt.txt
var Character string

//go:embed quest_prompt.txt
var Quest string
