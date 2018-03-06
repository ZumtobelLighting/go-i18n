package i18n

import (
	"bytes"
	"encoding"
	"strings"
	gotemplate "text/template"
)

// Template stores the template for a string.
type Template struct {
	Template *gotemplate.Template
	Src      string
}

// NewTemplate returns a new template from src.
func NewTemplate(src string) (*Template, error) {
	var tmpl Template
	err := tmpl.parseTemplate(src)
	return &tmpl, err
}

func (t *Template) String() string {
	return t.Src
}

// Execute executes the template for the given data.
func (t *Template) Execute(data interface{}) string {
	if t.Template == nil {
		return t.Src
	}
	var buf bytes.Buffer
	if err := t.Template.Execute(&buf, data); err != nil {
		return err.Error()
	}
	return buf.String()
}

// MarshalText implements the TextMarshaler interface.
func (t *Template) MarshalText() ([]byte, error) {
	return []byte(t.Src), nil
}

// UnmarshalText implements the TextUnmarshaler interface.
func (t *Template) UnmarshalText(src []byte) error {
	return t.parseTemplate(string(src))
}

func (t *Template) parseTemplate(src string) (err error) {
	t.Src = src
	if strings.Contains(src, "{{") {
		t.Template, err = gotemplate.New(src).Parse(src)
	}
	return
}

var _ = encoding.TextMarshaler(&Template{})
var _ = encoding.TextUnmarshaler(&Template{})
