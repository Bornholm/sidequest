package prompt

import (
	_ "embed"
)

//go:embed agent_prompt.txt.gotmpl
var Agent string

//go:embed character_prompt.txt.gotmpl
var Character string

//go:embed quest_prompt.txt.gotmpl
var Quest string
