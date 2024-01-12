package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"main/pkg/utils"
	templatesList "main/templates"
)

type TemplateManager struct {
	Templates map[string]*template.Template
}

func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		Templates: make(map[string]*template.Template, 0),
	}
}

func (manager *TemplateManager) GetTemplate(name string) (*template.Template, error) {
	if t, ok := manager.Templates[name]; ok {
		return t, nil
	}

	filename := fmt.Sprintf("%s.html", name)

	t, err := template.New(filename).
		Funcs(template.FuncMap{
			"SerializeLink":  utils.SerializeLink,
			"FormatSize":     utils.FormatSize,
			"FormatBool":     utils.FormatBool,
			"FormatDuration": utils.FormatDuration,
		}).
		ParseFS(templatesList.Templates, filename)
	if err != nil {
		return nil, err
	}

	manager.Templates[name] = t
	return t, nil
}

func (manager *TemplateManager) Render(name string, data interface{}) (string, error) {
	t, err := manager.GetTemplate(name)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
