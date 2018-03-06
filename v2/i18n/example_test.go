package i18n_test

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func ExampleTranslator_Translate_missingTranslation() {
	bundle := i18n.NewBundle()
	translator := i18n.NewTranslator(bundle, "es-es")
	translated, err := translator.Translate("HelloWorld")
	fmt.Println(translated)
	fmt.Println(err)
	// Output:
	//
	// <nil>
}

func ExampleTranslator_MustTranslate() {
	bundle := i18n.NewBundle()
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseTranslationFileBytes([]byte(`
HelloWorld = "Hola Mundo!"
`), "es.toml")
	translator := i18n.NewTranslator(bundle, "es-es")
	fmt.Println(translator.MustTranslate("HelloWorld"))
	// Output:
	// Hola Mundo!
}

func ExampleDefaultTranslator_MustTranslate() {
	bundle := i18n.NewBundle()
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseTranslationFileBytes([]byte(`
HelloWorld = "Hola Mundo!"
`), "es.toml")
	enTranslator := i18n.NewDefaultTranslator(bundle, "en", "en")
	fmt.Println(enTranslator.MustTranslate("HelloWorld", "Hello World!"))
	esTranslator := i18n.NewDefaultTranslator(bundle, "es-es", "en")
	fmt.Println(esTranslator.MustTranslate("HelloWorld", "Hello World!"))
	// Output:
	// Hello World!
	// Hola Mundo!
}

func ExampleTranslator_MustTranslate_template() {
	bundle := i18n.NewBundle()
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseTranslationFileBytes([]byte(`
HelloPerson = "Hola {{.Person}}!"
`), "es-es.toml")
	translator := i18n.NewTranslator(bundle, "es-es")
	bobMap := map[string]interface{}{"Person": "Bob"}
	bobStruct := struct{ Person string }{Person: "Bob"}
	fmt.Println(translator.MustTranslate("HelloPerson", bobMap))
	fmt.Println(translator.MustTranslate("HelloPerson", bobStruct))
	// Output:
	// Hola Bob!
	// Hola Bob!
}

func ExampleDefaultTranslator_MustTranslate_template() {
	bundle := i18n.NewBundle()
	translator := i18n.NewDefaultTranslator(bundle, "en", "en")
	bobMap := map[string]interface{}{"Person": "Bob"}
	bobStruct := struct{ Person string }{Person: "Bob"}
	fmt.Println(translator.MustTranslate("HelloPerson", "Hello {{.Person}}!", bobMap))
	fmt.Println(translator.MustTranslate("HelloPerson", "Hello {{.Person}}!", bobStruct))
	// Output:
	// Hello Bob!
	// Hello Bob!
}

func ExampleTranslator_MustTranslate_plural() {
	bundle := i18n.NewBundle()
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseTranslationFileBytes([]byte(`
[YourHeightInMeters]
One = "You are {{.Count}} meter tall."
Other = "You are {{.Count}} meters tall."
`), "en.toml")
	translator := i18n.NewTranslator(bundle, "en")
	fmt.Println(translator.MustTranslate("YourHeightInMeters", 0))
	fmt.Println(translator.MustTranslate("YourHeightInMeters", 1))
	fmt.Println(translator.MustTranslate("YourHeightInMeters", 2))
	fmt.Println(translator.MustTranslate("YourHeightInMeters", "1.7"))
	// Output:
	// You are 0 meters tall.
	// You are 1 meter tall.
	// You are 2 meters tall.
	// You are 1.7 meters tall.
}

func ExampleDefaultTranslator_MustTranslate_plural() {
	bundle := i18n.NewBundle()
	translator := i18n.NewDefaultTranslator(bundle, "en", "en")
	yourHeightInMeters := &i18n.PluralMessage{
		ID:          "YourHeightInMeters",
		Description: "A message that says how tall you are.",
		One:         "You are {{.Count}} meter tall.",
		Other:       "You are {{.Count}} meters tall.",
	}
	fmt.Println(translator.MustTranslate(yourHeightInMeters, 0))
	fmt.Println(translator.MustTranslate(yourHeightInMeters, 1))
	fmt.Println(translator.MustTranslate(yourHeightInMeters, 2))
	fmt.Println(translator.MustTranslate(yourHeightInMeters, "1.7"))
	// Output:
	// You are 0 meters tall.
	// You are 1 meter tall.
	// You are 2 meters tall.
	// You are 1.7 meters tall.
}

func ExampleTranslator_MustTranslate_pluralTemplate() {
	bundle := i18n.NewBundle()
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseTranslationFileBytes([]byte(`
[PersonHeightInMeters]
One = "{{.Person}} is {{.Count}} meter tall."
Other = "{{.Person}} is {{.Count}} meters tall."
`), "en-en.toml")
	translator := i18n.NewTranslator(bundle, "en-en")
	bobMap := map[string]interface{}{"Person": "Bob"}
	bobStruct := struct{ Person string }{Person: "Bob"}
	fmt.Println(translator.MustTranslate("PersonHeightInMeters", 0, bobMap))
	fmt.Println(translator.MustTranslate("PersonHeightInMeters", 0, bobStruct))
	fmt.Println(translator.MustTranslate("PersonHeightInMeters", 1, bobMap))
	fmt.Println(translator.MustTranslate("PersonHeightInMeters", 1, bobStruct))
	fmt.Println(translator.MustTranslate("PersonHeightInMeters", 2, bobMap))
	fmt.Println(translator.MustTranslate("PersonHeightInMeters", 2, bobStruct))
	fmt.Println(translator.MustTranslate("PersonHeightInMeters", "1.7", bobMap))
	fmt.Println(translator.MustTranslate("PersonHeightInMeters", "1.7", bobStruct))
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

func ExampleDefaultTranslator_MustTranslate_pluralTemplate() {
	bundle := i18n.NewBundle()
	translator := i18n.NewDefaultTranslator(bundle, "en", "en")
	personHeightInMeters := &i18n.PluralMessage{
		ID:          "PersonHeightInMeters",
		Description: "A message that says how tall a person is.",
		One:         "{{.Person}} is {{.Count}} meter tall.",
		Other:       "{{.Person}} is {{.Count}} meters tall.",
	}
	bobMap := map[string]interface{}{"Person": "Bob"}
	bobStruct := struct{ Person string }{Person: "Bob"}
	fmt.Println(translator.MustTranslate(personHeightInMeters, 0, bobMap))
	fmt.Println(translator.MustTranslate(personHeightInMeters, 0, bobStruct))
	fmt.Println(translator.MustTranslate(personHeightInMeters, 1, bobMap))
	fmt.Println(translator.MustTranslate(personHeightInMeters, 1, bobStruct))
	fmt.Println(translator.MustTranslate(personHeightInMeters, 2, bobMap))
	fmt.Println(translator.MustTranslate(personHeightInMeters, 2, bobStruct))
	fmt.Println(translator.MustTranslate(personHeightInMeters, "1.7", bobMap))
	fmt.Println(translator.MustTranslate(personHeightInMeters, "1.7", bobStruct))
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
