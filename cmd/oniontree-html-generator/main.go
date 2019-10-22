package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/onionltd/onionltd.github.io-generator/pkg/oniontree/local"
	"github.com/onionltd/oniontree-tools/pkg/types/service"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func loadTemplates(dir string) (*template.Template, error) {
	log.Printf("Load templates")
	return template.ParseGlob(dir + "/*.*")
}

func generateIndexHTML(t *template.Template) (string, error) {
	log.Printf("Generate index.html")
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	if err := t.ExecuteTemplate(&buffer, "index.html", nil); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func generateSearchHTML(t *template.Template) (string, error) {
	log.Printf("Generate search.html")
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	if err := t.ExecuteTemplate(&buffer, "search.html", nil); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func generateTagHTML(t *template.Template, name string, services []string) (string, error) {
	log.Printf("Generate tag: tags/%s.html", name)
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	data := struct {
		Name     string
		Services []string
	}{
		name,
		services,
	}
	if err := t.ExecuteTemplate(&buffer, "tag.html", data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func generateTagsHTML(t *template.Template, tags []string) (string, error) {
	log.Printf("Generate tags/index.html")
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	if err := t.ExecuteTemplate(&buffer, "tags.html", tags); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func generateServicesHTML(t *template.Template, services []string) (string, error) {
	log.Printf("Generate services/index.html")
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	if err := t.ExecuteTemplate(&buffer, "services.html", services); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func generateServiceHTML(t *template.Template, s service.Service) (string, error) {
	log.Printf("Generate service: %s", s.Name)
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	if err := t.ExecuteTemplate(&buffer, "service.html", s); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func main() {
	templates := flag.String("templates", "", "Directory with templates")
	data := flag.String("oniontree", "", "Oniontree directory")
	output := flag.String("output", ".", "Output directory")
	flag.Parse()

	if *templates == "" {
		panic(fmt.Errorf("templates not specified"))
	}
	if *data == "" {
		panic(fmt.Errorf("oniontree data not specified"))
	}

	oniontree := local.NewSource(*data)

	t, err := loadTemplates(*templates)
	if err != nil {
		panic(err)
	}

	html, err := generateIndexHTML(t)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(*output, "index.html"), []byte(html), 0644); err != nil {
		panic(err)
	}
	html, err = generateSearchHTML(t)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(*output, "search.html"), []byte(html), 0644); err != nil {
		panic(err)
	}
	services, err := oniontree.ListServiceFiles()
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(path.Join(*output, "services"), 0755); err != nil {
		panic(err)
	}
	html, err = generateServicesHTML(t, services)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(*output, "services", "index.html"), []byte(html), 0644); err != nil {
		panic(err)
	}
	for _, id := range services {
		s, err := oniontree.GetServiceFile(id)
		if err != nil {
			panic(err)
		}
		html, err := generateServiceHTML(t, s)
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(path.Join(*output, "services", id+".html"), []byte(html), 0644); err != nil {
			panic(err)
		}
	}

	// Generate tags
	tags, err := oniontree.ListTags()
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(path.Join(*output, "tags"), 0755); err != nil {
		panic(err)
	}
	html, err = generateTagsHTML(t, tags)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(*output, "tags", "index.html"), []byte(html), 0644); err != nil {
		panic(err)
	}
	// Generate a file for each tag
	for _, tag := range tags {
		services, err := oniontree.ListTag(tag)
		if err != nil {
			panic(err)
		}
		html, err := generateTagHTML(t, tag, services)
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(path.Join(*output, "tags", tag+".html"), []byte(html), 0644); err != nil {
			panic(err)
		}
	}
}
