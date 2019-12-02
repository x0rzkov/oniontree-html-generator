package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/onionltd/oniontree-tools/pkg/oniontree"
	"github.com/onionltd/oniontree-tools/pkg/types/service"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

const (
	targetClearnet string = "clearnet"
	targetOnion    string = "onion"
)

var targets = []string{targetClearnet, targetOnion}

func loadTemplates(dir string, funcs *TF) (*template.Template, error) {
	log.Printf("Load templates")
	return template.New("").Funcs(template.FuncMap{
		"getTarget":                 funcs.tfTarget,
		"toUpper":                   funcs.tfToUpper,
		"serviceLastUpdate":         funcs.tfServiceFileLastUpdate,
		"onionTreeBookmarksVersion": funcs.tfOnionTreeBookmarksVersion,
	}).ParseGlob(dir + "/*.*")
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

func generateDownloadHTML(output string, t *template.Template) error {
	log.Printf("Generate download.html")
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	if err := t.ExecuteTemplate(&buffer, "download.html", nil); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(output, "download.html"), buffer.Bytes(), 0644); err != nil {
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

func generateServiceHTML(output string, t *template.Template, id string, s service.Service, tags []string, alerts []map[string]string) error {
	log.Printf("Generate service: services/%s.html", id)
	buffer := bytes.Buffer{}
	bufio.NewWriter(&buffer)
	data := struct {
		ID      string
		Tags    []string
		Service service.Service
		Alerts  []map[string]string
	}{
		id,
		tags,
		s,
		alerts,
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
	omitTags := flag.String("frontpage-omit-tags", "", "Services tagged with these tags won't show on the frontpage")
	alertFile := flag.String("alerts", "", "YAML file containing service alerts")
	otbVersion := flag.String("otb-version", "", "OnionTree Bookmarks version")
	target := flag.String("target", targetClearnet, fmt.Sprintf("Generate for target [%s]", strings.Join(targets, ",")))
	flag.Parse()

	if *templates == "" {
		panic(fmt.Errorf("templates not specified"))
	}
	if *data == "" {
		panic(fmt.Errorf("oniontree data not specified"))
	}
	if *otbVersion == "" {
		panic(fmt.Errorf("oniontree bookmarks version not specified"))
	}

	alerts := map[string][]map[string]string{}
	if *alertFile != "" {
		b, err := ioutil.ReadFile(*alertFile)
		if err != nil {
			panic(err)
		}
		if err := yaml.Unmarshal(b, &alerts); err != nil {
			panic(err)
		}
	}

	omittedTags := []string{}
	if *omitTags != "" {
		for _, t := range strings.Split(*omitTags, ",") {
			omittedTags = append(omittedTags, strings.TrimSpace(t))
		}
	}

	onionTree, err := oniontree.Open(*data)
	if err != nil {
		panic(err)
	}

	t, err := loadTemplates(*templates, &TF{
		dataPath:   *data,
		otbVersion: *otbVersion,
		target:     *target,
	})
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

	// Generate download
	if err = generateDownloadHTML(*output, t); err != nil {
		panic(err)
	}

	// Generate services
	serviceIDs, err := onionTree.List()
	if err != nil {
		panic(err)
	}
	ids := []string{}
	services := []service.Service{}
	for _, id := range serviceIDs {
		s, err := onionTree.Get(id)
		if err != nil {
			panic(err)
		}
		tags, err := listServiceTags(onionTree, id)
		if err != nil {
			panic(err)
		}
		serviceAlerts, ok := alerts[id]
		if !ok {
			serviceAlerts = nil
		}
		if err := generateServiceHTML(*output, t, id, s, tags, serviceAlerts); err != nil {
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
		// Is service tagged with tag that should be omitted from the frontpage? If it is, jump to next service.
		if containsIgnoredTag(tags, omittedTags...) {
			log.Printf("Omitting service: services/%s.html (tagged as one of: %s)", id, omittedTags)
			continue
		}
		ids = append(ids, id)
		services = append(services, s)
	}
	if err = generateServicesHTML(*output, t, ids, services); err != nil {
		panic(err)
	}

	// Generate tags
	tags, err := onionTree.ListTags()
	if err != nil {
		panic(err)
	}
	if err = generateTagsHTML(*output, t, tags); err != nil {
		panic(err)
	}
	for idx, _ := range tags {
		tag, err := onionTree.GetTag(tags[idx])
		if err != nil {
			panic(err)
		}
		services := []service.Service{}
		for _, id := range tag.Services {
			s, err := onionTree.Get(id)
			if err != nil {
				panic(err)
			}
			services = append(services, s)
		}
		if err := generateTagHTML(*output, t, tag.ID, tag.Services, services); err != nil {
			panic(err)
		}
	}
}

func containsIgnoredTag(slice []string, tags ...string) bool {
	for _, item := range slice {
		for _, t := range tags {
			if t == item {
				return true
			}
		}
	}
	return false
}

func listServiceTags(onionTree *oniontree.OnionTree, id string) ([]string, error) {
	res := []string{}
	tags, err := onionTree.ListTags()
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		t, err := onionTree.GetTag(tag)
		if err != nil {
			return nil, err
		}
		for _, service := range t.Services {
			if service == id {
				res = append(res, tag)
				break
			}
		}
	}
	sort.Strings(res)
	return res, nil
}
