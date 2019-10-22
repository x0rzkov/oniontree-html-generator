package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/onionltd/oniontree-tools/pkg/types/service"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func ListTag(name string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/onionltd/oniontree/contents/tagged/%s", name)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid API response code %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := make([]map[string]interface{}, 128)
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}

	services := []string{}
	for _, t := range data {
		services = append(services, trimSuffix(t["name"].(string)))
	}
	return services, nil
}

func ListTags() ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/onionltd/oniontree/contents/tagged")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid API response code %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := make([]map[string]interface{}, 128)
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}

	tags := []string{}
	for _, t := range data {
		tags = append(tags, t["name"].(string))
	}
	return tags, nil
}

func ListServiceFiles() ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/onionltd/oniontree/contents/unsorted")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid API response code %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := make([]map[string]interface{}, 128)
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}

	ids := []string{}
	for _, s := range data {
		ids = append(ids, trimSuffix(s["name"].(string)))
	}
	return ids, nil
}

func GetServiceFile(id string) (service.Service, error) {
	url := fmt.Sprintf("https://api.github.com/repos/onionltd/oniontree/contents/unsorted/%s.yaml", id)
	resp, err := http.Get(url)
	if err != nil {
		return service.Service{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return service.Service{}, fmt.Errorf("invalid API response code %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return service.Service{}, err
	}

	data := make(map[string]interface{})
	if err := json.Unmarshal(respBody, &data); err != nil {
		return service.Service{}, err
	}

	if _, ok := data["content"]; !ok {
		return service.Service{}, fmt.Errorf("content field not defined")
	}

	if _, ok := data["content"].(string); !ok {
		return service.Service{}, fmt.Errorf("content field not defined")
	}

	content := strings.Replace(data["content"].(string), "\\n", "\n", -1)

	contentDecoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return service.Service{}, err
	}

	s := service.Service{}
	if err := yaml.Unmarshal(contentDecoded, &s); err != nil {
		return service.Service{}, err
	}

	return s, nil
}

func trimSuffix(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}
