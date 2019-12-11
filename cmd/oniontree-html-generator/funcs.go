package main

import (
	"os"
	"path"
	"strings"
)

type TF struct {
	dataPath string
	// OnionTree Bookmarks version
	otbVersion string
	target     string
}

func (t TF) tfTarget() string {
	return t.target
}

func (t TF) tfServiceFileLastUpdate(id string) string {
	info, err := os.Stat(path.Join(t.dataPath, "unsorted", id+".yaml"))
	if err != nil {
		return ""
	}
	return info.ModTime().Format("2006-01-02 15:04:05 MST")
}

func (t TF) tfToUpper(s string) string {
	return strings.ToUpper(s)
}

func (t TF) tfOnionTreeBookmarksVersion() string {
	return t.otbVersion
}

func (t TF) tfNL2Space(s string) string {
	return strings.Replace(s, "\n", " ", -1)
}
