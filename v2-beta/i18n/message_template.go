package i18n

// MessageTemplate is an executable template for a message.
type MessageTemplate struct {
	ID          string
	Description string
	PluralForms map[PluralForm]*Template
}

// NewMessageTemplate returns a new MessageTemplate parsed from data.
func NewMessageTemplate(id string, data map[string]string) (*MessageTemplate, error) {
	translation := &MessageTemplate{
		ID:          id,
		PluralForms: make(map[PluralForm]*Template),
	}
	for k, v := range data {
		switch k {
		case "description":
			translation.Description = v
		default:
			pluralForm, err := NewPluralForm(k)
			if err != nil {
				return nil, err
			}
			tmpl, err := NewTemplate(v)
			if err != nil {
				return nil, err
			}
			translation.PluralForms[pluralForm] = tmpl
		}
	}
	return translation, nil
}

// MustNewMessageTemplate is similar to NewMessageTemplate except it panics if an error happens.
func MustNewMessageTemplate(id string, data map[string]string) *MessageTemplate {
	t, err := NewMessageTemplate(id, data)
	if err != nil {
		panic(err)
	}
	return t
}

// Execute executes the template for the plural form
// and template data.
func (t *MessageTemplate) Execute(pluralForm PluralForm, data interface{}) string {
	tmpl := t.PluralForms[pluralForm]
	return tmpl.Execute(data)
}
