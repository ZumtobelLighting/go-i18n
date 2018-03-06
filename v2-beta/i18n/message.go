package i18n

// Message is a message that can be localized.
type Message struct {
	// ID uniquely identifies the message.
	ID string

	// Description describes the message to give additional
	// context to translators that may be relevant for translation.
	Description string

	// Zero is the content of the message for the CLDR plural form "zero".
	Zero string

	// One is the content of the message for the CLDR plural form "one".
	One string

	// Two is the content of the message for the CLDR plural form "two".
	Two string

	// Few is the content of the message for the CLDR plural form "few".
	Few string

	// Many is the content of the message for the CLDR plural form "many".
	Many string

	// Otherko is the content of the message for the CLDR plural form "other"
	Other string
}

// Template returns a new MessageTemplate with the same data as the message.
func (m *Message) Template() (*MessageTemplate, error) {
	return NewMessageTemplate(m.ID, map[string]string{
		"description": m.Description,
		"zero":        m.Zero,
		"one":         m.One,
		"two":         m.Two,
		"few":         m.Few,
		"many":        m.Many,
		"other":       m.Other,
	})
}
