package i18n_test

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2-beta/i18n"
)

func ExampleLocalizer_Localize_missingTranslation() {
	bundle := i18n.NewBundle()
	localizer := i18n.NewLocalizer(bundle, "es-es")
	localized, err := localizer.Localize("HelloWorld")
	fmt.Println(localized)
	fmt.Println(err)
	// Output:
	//
	// <nil>
}

func ExampleLocalizer_MustLocalize() {
	bundle := i18n.NewBundle()
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
HelloWorld = "Hola Mundo!"
`), "es.toml")
	localizer := i18n.NewLocalizer(bundle, "es-es")
	fmt.Println(localizer.MustLocalize("HelloWorld"))
	// Output:
	// Hola Mundo!
}

func ExampleDefaultLocalizer_MustLocalize() {
	bundle := i18n.NewBundle()
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
HelloWorld = "Hola Mundo!"
`), "es.toml")
	enTranslator := i18n.NewDefaultLocalizer(bundle, "en", "en")
	fmt.Println(enTranslator.MustLocalize("HelloWorld", "Hello World!"))
	esTranslator := i18n.NewDefaultLocalizer(bundle, "es-es", "en")
	fmt.Println(esTranslator.MustLocalize("HelloWorld", "Hello World!"))
	// Output:
	// Hello World!
	// Hola Mundo!
}

func ExampleLocalizer_MustLocalize_template() {
	bundle := i18n.NewBundle()
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
HelloPerson = "Hola {{.Person}}!"
`), "es-es.toml")
	localizer := i18n.NewLocalizer(bundle, "es-es")
	bobMap := map[string]interface{}{"Person": "Bob"}
	bobStruct := struct{ Person string }{Person: "Bob"}
	fmt.Println(localizer.MustLocalize("HelloPerson", bobMap))
	fmt.Println(localizer.MustLocalize("HelloPerson", bobStruct))
	// Output:
	// Hola Bob!
	// Hola Bob!
}

func ExampleDefaultLocalizer_MustLocalize_template() {
	bundle := i18n.NewBundle()
	localizer := i18n.NewDefaultLocalizer(bundle, "en", "en")
	bobMap := map[string]interface{}{"Person": "Bob"}
	bobStruct := struct{ Person string }{Person: "Bob"}
	fmt.Println(localizer.MustLocalize("HelloPerson", "Hello {{.Person}}!", bobMap))
	fmt.Println(localizer.MustLocalize("HelloPerson", "Hello {{.Person}}!", bobStruct))
	// Output:
	// Hello Bob!
	// Hello Bob!
}

func ExampleLocalizer_MustLocalize_plural() {
	bundle := i18n.NewBundle()
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
[YourHeightInMeters]
One = "You are {{.Count}} meter tall."
Other = "You are {{.Count}} meters tall."
`), "en.toml")
	localizer := i18n.NewLocalizer(bundle, "en")
	fmt.Println(localizer.MustLocalize("YourHeightInMeters", 0))
	fmt.Println(localizer.MustLocalize("YourHeightInMeters", 1))
	fmt.Println(localizer.MustLocalize("YourHeightInMeters", 2))
	fmt.Println(localizer.MustLocalize("YourHeightInMeters", "1.7"))
	// Output:
	// You are 0 meters tall.
	// You are 1 meter tall.
	// You are 2 meters tall.
	// You are 1.7 meters tall.
}

func ExampleDefaultLocalizer_MustLocalize_plural() {
	bundle := i18n.NewBundle()
	localizer := i18n.NewDefaultLocalizer(bundle, "en", "en")
	yourHeightInMeters := &i18n.Message{
		ID:          "YourHeightInMeters",
		Description: "A message that says how tall you are.",
		One:         "You are {{.Count}} meter tall.",
		Other:       "You are {{.Count}} meters tall.",
	}
	fmt.Println(localizer.MustLocalize(yourHeightInMeters, 0))
	fmt.Println(localizer.MustLocalize(yourHeightInMeters, 1))
	fmt.Println(localizer.MustLocalize(yourHeightInMeters, 2))
	fmt.Println(localizer.MustLocalize(yourHeightInMeters, "1.7"))
	// Output:
	// You are 0 meters tall.
	// You are 1 meter tall.
	// You are 2 meters tall.
	// You are 1.7 meters tall.
}

func ExampleLocalizer_MustLocalize_pluralTemplate() {
	bundle := i18n.NewBundle()
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
[PersonHeightInMeters]
One = "{{.Person}} is {{.Count}} meter tall."
Other = "{{.Person}} is {{.Count}} meters tall."
`), "en-en.toml")
	localizer := i18n.NewLocalizer(bundle, "en-en")
	bobMap := map[string]interface{}{"Person": "Bob"}
	bobStruct := struct{ Person string }{Person: "Bob"}
	fmt.Println(localizer.MustLocalize("PersonHeightInMeters", 0, bobMap))
	fmt.Println(localizer.MustLocalize("PersonHeightInMeters", 0, bobStruct))
	fmt.Println(localizer.MustLocalize("PersonHeightInMeters", 1, bobMap))
	fmt.Println(localizer.MustLocalize("PersonHeightInMeters", 1, bobStruct))
	fmt.Println(localizer.MustLocalize("PersonHeightInMeters", 2, bobMap))
	fmt.Println(localizer.MustLocalize("PersonHeightInMeters", 2, bobStruct))
	fmt.Println(localizer.MustLocalize("PersonHeightInMeters", "1.7", bobMap))
	fmt.Println(localizer.MustLocalize("PersonHeightInMeters", "1.7", bobStruct))
	// Output:
	// Bob is 0 meters tall.
	// Bob is 0 meters tall.
	// Bob is 1 meter tall.
	// Bob is 1 meter tall.
	// Bob is 2 meters tall.
	// Bob is 2 meters tall.
	// Bob is 1.7 meters tall.
	// Bob is 1.7 meters tall.
}

func ExampleDefaultLocalizer_MustLocalize_pluralTemplate() {
	bundle := i18n.NewBundle()
	localizer := i18n.NewDefaultLocalizer(bundle, "en", "en")
	personHeightInMeters := &i18n.Message{
		ID:          "PersonHeightInMeters",
		Description: "A message that says how tall a person is.",
		One:         "{{.Person}} is {{.Count}} meter tall.",
		Other:       "{{.Person}} is {{.Count}} meters tall.",
	}
	bobMap := map[string]interface{}{"Person": "Bob"}
	bobStruct := struct{ Person string }{Person: "Bob"}
	fmt.Println(localizer.MustLocalize(personHeightInMeters, 0, bobMap))
	fmt.Println(localizer.MustLocalize(personHeightInMeters, 0, bobStruct))
	fmt.Println(localizer.MustLocalize(personHeightInMeters, 1, bobMap))
	fmt.Println(localizer.MustLocalize(personHeightInMeters, 1, bobStruct))
	fmt.Println(localizer.MustLocalize(personHeightInMeters, 2, bobMap))
	fmt.Println(localizer.MustLocalize(personHeightInMeters, 2, bobStruct))
	fmt.Println(localizer.MustLocalize(personHeightInMeters, "1.7", bobMap))
	fmt.Println(localizer.MustLocalize(personHeightInMeters, "1.7", bobStruct))
	// Output:
	// Bob is 0 meters tall.
	// Bob is 0 meters tall.
	// Bob is 1 meter tall.
	// Bob is 1 meter tall.
	// Bob is 2 meters tall.
	// Bob is 2 meters tall.
	// Bob is 1.7 meters tall.
	// Bob is 1.7 meters tall.
}
