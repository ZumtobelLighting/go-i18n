package i18n

// TODO: rename to MessageTemplate?

// Translation contains the data for a single translation.
type Translation struct {
	ID          string
	Description string
	PluralForms map[PluralForm]*Template
}

// NewTranslation returns a new translation parsed from data.
// It returns an error if data contains invalid plural forms
// or invalid translation templates.
func NewTranslation(id string, data map[string]string) (*Translation, error) {
	translation := &Translation{
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

// MustNewTranslation is similar to NewTranslation except it panics if an error happens.
func MustNewTranslation(id string, data map[string]string) *Translation {
	t, err := NewTranslation(id, data)
	if err != nil {
		panic(err)
	}
	return t
}

// Translate returns the translated string for the plural form
// and template data.
func (t *Translation) Translate(pluralForm PluralForm, data interface{}) string {
	tmpl := t.PluralForms[pluralForm]
	return tmpl.Execute(data)
}
