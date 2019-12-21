package main

import (
	"fmt"
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

func (t TF) tfFormatPGPFingerprint(s string) string {
	return fmt.Sprintf("%s %s %s %s %s  %s %s %s %s %s",
		s[0:4], s[4:8], s[8:12], s[12:16], s[16:20],
		s[20:24], s[24:28], s[28:32], s[32:36], s[36:40])
}
