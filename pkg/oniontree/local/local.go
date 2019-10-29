package local

import (
	"github.com/go-yaml/yaml"
	"github.com/onionltd/oniontree-tools/pkg/types/service"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

type Source struct {
	dir string
}

func (s Source) ListTag(name string) ([]string, error) {
	dirPath := path.Join(s.dir, "tagged", name)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	services := []string{}
	for _, f := range files {
		services = append(services, trimSuffix(f.Name()))
	}
	return services, nil
}

func (s Source) ListTags() ([]string, error) {
	dirPath := path.Join(s.dir, "tagged")
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	tags := []string{}
	for _, f := range files {
		tags = append(tags, f.Name())
	}
	return tags, nil
}

func (s Source) ListServiceFiles() ([]string, error) {
	dirPath := path.Join(s.dir, "unsorted")
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	ids := []string{}
	for _, f := range files {
		ids = append(ids, trimSuffix(f.Name()))
	}
	return ids, nil
}

func (s Source) GetServiceFile(id string) (service.Service, error) {
	filePath := path.Join(s.dir, "unsorted", id+".yaml")
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return service.Service{}, err
	}

	serv := service.Service{}
	if err := yaml.Unmarshal(content, &serv); err != nil {
		return service.Service{}, err
	}

	return serv, nil
}

func (s Source) ListServiceTags(id string) ([]string, error) {
	tags, err := s.ListTags()
	if err != nil {
		return []string{}, err
	}

	serviceTags := []string{}
	for _, tag := range tags {
		services, err := s.ListTag(tag)
		if err != nil {
			return []string{}, err
		}
		for _, service := range services {
			if service == id {
				serviceTags = append(serviceTags, tag)
			}
		}
	}
	return serviceTags, nil
}

func trimSuffix(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

func NewSource(dir string) Source {
	return Source{dir}
}
