package route

import (
	"bytes"
	"html/template"

	"github.com/pkg/errors"
)

func generatePrompt(rawTemplate string, data any) (string, error) {
	promptTmpl, err := template.New("").Parse(rawTemplate)
	if err != nil {
		return "", errors.WithStack(err)
	}

	var prompt bytes.Buffer

	if err := promptTmpl.Execute(&prompt, data); err != nil {
		return "", errors.WithStack(err)
	}

	return prompt.String(), nil
}
