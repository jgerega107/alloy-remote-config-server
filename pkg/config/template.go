package config

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var templates = make(map[string]*template.Template)

func LoadTemplates(path string) error {
	files, err := filepath.Glob(filepath.Join(path, "*.conf.tmpl"))
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		fullName := filepath.Base(file)
		tmpl, err := template.New(fullName).Parse(string(content))
		if err != nil {
			return err
		}
		trimmedName := strings.TrimSuffix(fullName, ".conf.tmpl")
		if _, exists := templates[trimmedName]; exists {
			log.Printf("Warning: Template '%s' already loaded, overwriting with %s", trimmedName, file)
		}
		templates[trimmedName] = tmpl
	}

	return nil
}
