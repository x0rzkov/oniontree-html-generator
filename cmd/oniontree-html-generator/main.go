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

func generateSearchHTML(output string, t *template.Template) error {
	log.Printf("Generate search.html")
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	if err := t.ExecuteTemplate(&buffer, "search.html", nil); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(output, "search.html"), buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func generateAPIHTML(output string, t *template.Template) error {
	log.Printf("Generate api.html")
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	if err := t.ExecuteTemplate(&buffer, "api.html", nil); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(output, "api.html"), buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func generateTagHTML(output string, t *template.Template, name string, ids []string, services []service.Service) error {
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
		return err
	}
	if err := ioutil.WriteFile(path.Join(output, "tags", name+".html"), buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func generateTagsHTML(output string, t *template.Template, tags []string) error {
	log.Printf("Generate tags/index.html")
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	data := make(map[string][]string)
	for _, tag := range tags {
		letter := strings.ToUpper(string(tag[0]))
		data[letter] = append(data[letter], tag)
	}
	if err := t.ExecuteTemplate(&buffer, "tags.html", data); err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(output, "tags"), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(output, "tags", "index.html"), buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func generateServicesHTML(output string, t *template.Template, ids []string, services []service.Service) error {
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
		return err
	}
	if err := ioutil.WriteFile(path.Join(output, "index.html"), buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func generateServiceHTML(output string, t *template.Template, id string, s service.Service, tags []string) error {
	log.Printf("Generate service: services/%s.html", id)
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
		return err
	}
	if err := os.MkdirAll(path.Join(output, "services"), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(output, "services", id+".html"), buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func generateServiceJSON(output string, id string, s service.Service) error {
	log.Printf("Generate service: services/%s.json", id)
	buff, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(output, "services", id+".json"), buff, 0644); err != nil {
		return err
	}
	return nil
}

func generatePGPTXT(output string, k service.PublicKey) error {
	log.Printf("Generate service: keys/%s.txt", k.ID)
	if err := os.MkdirAll(path.Join(output, "keys"), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(output, "keys", k.ID+".txt"), []byte(k.Value), 0644); err != nil {
		return err
	}
	return nil
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
	if err = generateSearchHTML(*output, t); err != nil {
		panic(err)
	}

	// Generate api docs
	if err = generateAPIHTML(*output, t); err != nil {
		panic(err)
	}

	// Generate services
	serviceIDs, err := oniontree.ListServiceFiles()
	if err != nil {
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
		if err := generateServiceHTML(*output, t, id, s, tags); err != nil {
			panic(err)
		}
		if err := generateServiceJSON(*output, id, s); err != nil {
			panic(err)
		}
		for _, k := range s.PublicKeys {
			if err := generatePGPTXT(*output, k); err != nil {
				panic(err)
			}
		}
		services = append(services, s)
	}
	if err = generateServicesHTML(*output, t, serviceIDs, services); err != nil {
		panic(err)
	}

	// Generate tags
	tags, err := oniontree.ListTags()
	if err != nil {
		panic(err)
	}
	if err = generateTagsHTML(*output, t, tags); err != nil {
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
		if err := generateTagHTML(*output, t, tag, serviceIDs, services); err != nil {
			panic(err)
		}
	}
}
