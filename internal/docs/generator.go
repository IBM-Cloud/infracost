package docs

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/infracost/infracost/internal/providers/terraform"
)

func generateSupportedResourcesDocs(docsTemplatesPath string, outputPath string) error {
	tmpl, err := template.ParseFiles(docsTemplatesPath + "/supported_resources.md")
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Clean(outputPath + "/supported_resources.md"))
	if err != nil {
		return err
	}
	resourceRegistryMap := terraform.GetResourceRegistryMap()
	err = tmpl.Execute(f, resourceRegistryMap)
	if err != nil {
		return err
	}
	return nil
}

func GenerateDocs(docsTemplatesPath, outputPath string) error {
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}
	err = generateSupportedResourcesDocs(docsTemplatesPath, outputPath)
	if err != nil {
		return err
	}
	return nil
}
