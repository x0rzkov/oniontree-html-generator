package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/onionltd/onionltd.github.io-generator/pkg/oniontree/local"
	"github.com/onionltd/oniontree-tools/pkg/types/service"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func loadTemplates(dir string) (*template.Template, error) {
	log.Printf("Load templates")
	return template.ParseGlob(dir + "/*.*")
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

func generateTagHTML(t *template.Template, name string, ids []string, services []service.Service) (string, error) {
	log.Printf("Generate tag: tags/%s.html", name)
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	data := struct {
		Name string
		Data map[string][]struct {
			ID      string
			Service service.Service
		}
	}{
		name,
		make(map[string][]struct {
			ID      string
			Service service.Service
		}),
	}
	for idx, s := range services {
		letter := strings.ToUpper(string(s.Name[0]))
		data.Data[letter] = append(data.Data[letter], struct {
			ID      string
			Service service.Service
		}{
			ID:      ids[idx],
			Service: s,
		})
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
	data := make(map[string][]string)
	for _, tag := range tags {
		letter := strings.ToUpper(string(tag[0]))
		data[letter] = append(data[letter], tag)
	}
	if err := t.ExecuteTemplate(&buffer, "tags.html", data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func generateServicesHTML(t *template.Template, ids []string, services []service.Service) (string, error) {
	log.Printf("Generate services/index.html")
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	data := make(map[string][]struct {
		ID      string
		Service service.Service
	})
	for idx, s := range services {
		letter := strings.ToUpper(string(s.Name[0]))
		data[letter] = append(data[letter], struct {
			ID      string
			Service service.Service
		}{
			ID:      ids[idx],
			Service: s,
		})
	}
	if err := t.ExecuteTemplate(&buffer, "services.html", data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func generateServiceHTML(t *template.Template, id string, tags []string, s service.Service) (string, error) {
	log.Printf("Generate service: %s", s.Name)
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	data := struct {
		ID      string
		Tags    []string
		Service service.Service
	}{
		id,
		tags,
		s,
	}
	if err := t.ExecuteTemplate(&buffer, "service.html", data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func generateServiceJSON(s service.Service) (string, error) {
	log.Printf("Generate service JSON: %s", s.Name)
	buff, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(buff), nil
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

	// Generate search
	html, err := generateSearchHTML(t)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(*output, "search.html"), []byte(html), 0644); err != nil {
		panic(err)
	}

	// Generate services
	serviceIDs, err := oniontree.ListServiceFiles()
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(path.Join(*output, "services"), 0755); err != nil {
		panic(err)
	}
	services := []service.Service{}
	for _, id := range serviceIDs {
		s, err := oniontree.GetServiceFile(id)
		if err != nil {
			panic(err)
		}
		tags, err := oniontree.ListServiceTags(id)
		if err != nil {
			panic(err)
		}
		html, err := generateServiceHTML(t, id, tags, s)
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(path.Join(*output, "services", id+".html"), []byte(html), 0644); err != nil {
			panic(err)
		}
		json, err := generateServiceJSON(s)
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(path.Join(*output, "services", id+".json"), []byte(json), 0644); err != nil {
			panic(err)
		}
		services = append(services, s)
	}
	html, err = generateServicesHTML(t, serviceIDs, services)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(*output, "index.html"), []byte(html), 0644); err != nil {
		panic(err)
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
	for _, tag := range tags {
		serviceIDs, err := oniontree.ListTag(tag)
		if err != nil {
			panic(err)
		}
		services := []service.Service{}
		for _, id := range serviceIDs {
			s, err := oniontree.GetServiceFile(id)
			if err != nil {
				panic(err)
			}
			services = append(services, s)
		}
		html, err := generateTagHTML(t, tag, serviceIDs, services)
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(path.Join(*output, "tags", tag+".html"), []byte(html), 0644); err != nil {
			panic(err)
		}
	}
}
