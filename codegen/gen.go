package codegen

import (
	"bufio"
	"bytes"
	"embed"
	"go-faker/db"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

type M struct {
	StructName string
	TableName  string
	Fields     []string
	SM         map[string]string
}

func GenTableFile(name string) error {
	tmpl := template.Must(template.New("table").ParseFS(templates, "templates/*.tmpl"))
	m := M{
		StructName: ToCamelCase(name),
		TableName:  name,
		Fields:     make([]string, 0),
	}

	schema, err := db.DB.ParseSchema(name)
	if err != nil {
		return err
	}

	m.Fields = schema

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err = tmpl.ExecuteTemplate(w, "tableFile", &m); err != nil {
		return err
	}
	if err = w.Flush(); err != nil {
		return err
	}

	code, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}

	if !strings.HasSuffix(name, ".go") {
		name += ".go"
	}

	if err = os.MkdirAll(filepath.Dir("model"), 0755); err != nil {
		return err
	}
	if _, err = os.Stat("model/" + name); os.IsNotExist(err) {
		if err = os.WriteFile("model/"+name, code, 0644); err != nil {
			return err
		}
	}

	if err = GenModelFile(); err != nil {
		return err
	}

	return nil
}

func GenModelFile() error {
	tmpl := template.Must(template.New("table").ParseFS(templates, "templates/*.tmpl"))
	m := M{
		SM: map[string]string{},
	}

	files, err := ioutil.ReadDir("model/")
	if err != nil {
		return err
	}

	for _, file := range files {
		name := strings.TrimSuffix(file.Name(), ".go")
		if name != "model" {
			m.SM[name] = ToCamelCase(name)
		}
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err = tmpl.ExecuteTemplate(w, "modelFile", &m); err != nil {
		return err
	}
	if err = w.Flush(); err != nil {
		return err
	}

	code, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir("model"), 0755); err != nil {
		return err
	}
	if err = os.WriteFile("model/model.go", code, 0644); err != nil {
		return err
	}

	return nil
}

func ToCamelCase(inputUnderScoreStr string) (camelCase string) {
	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}
