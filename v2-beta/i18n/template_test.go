package i18n

import (
	"bytes"
	"fmt"
	"testing"
	gotemplate "text/template"
)

func TestNilTemplate(t *testing.T) {
	expected := "hello"
	tmpl := &Template{
		Template: nil,
		Src:      expected,
	}
	if actual := tmpl.Execute(nil); actual != expected {
		t.Errorf("Execute(nil) returned %s; expected %s", actual, expected)
	}
}

func TestMarshalText(t *testing.T) {
	tmpl := &Template{
		Template: gotemplate.Must(gotemplate.New("id").Parse("this is a {{.foo}} template")),
		Src:      "boom",
	}
	expectedBuf := []byte(tmpl.Src)
	if buf, err := tmpl.MarshalText(); !bytes.Equal(buf, expectedBuf) || err != nil {
		t.Errorf("MarshalText() returned %#v, %#v; expected %#v, nil", buf, err, expectedBuf)
	}
}

func TestUnmarshalText(t *testing.T) {
	tmpl := &Template{}
	tmpl.UnmarshalText([]byte("hello {{.World}}"))
	result := tmpl.Execute(map[string]string{
		"World": "world!",
	})
	expected := "hello world!"
	if result != expected {
		t.Errorf("expected %#v; got %#v", expected, result)
	}
}

var benchmarkResult string

func BenchmarkExecuteNilTemplate(b *testing.B) {
	template := &Template{Src: "hello world"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkResult = template.Execute(nil)
	}
}

func BenchmarkExecuteHelloWorldTemplate(b *testing.B) {
	template, err := NewTemplate("hello world")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkResult = template.Execute(nil)
	}
}

// Executing a simple template like this is ~6x slower than Sprintf
// but it is still only a few microseconds which should be sufficiently fast.
// The benefit is that we have nice semantic tags in the translation.
func BenchmarkExecuteHelloNameTemplate(b *testing.B) {
	template, err := NewTemplate("hello {{.Name}}")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkResult = template.Execute(map[string]string{
			"Name": "Nick",
		})
	}
}

func BenchmarkSprintf(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkResult = fmt.Sprintf("hello %s", "nick")
	}
}
