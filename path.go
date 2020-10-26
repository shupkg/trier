package trier

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	exeName       string
	fetchNameOnce sync.Once
)

func ExeName() string {
	fetchNameOnce.Do(func() {
		name, _ := os.Executable()
		name = filepath.Base(name)
		ext := filepath.Ext(name)
		if ext != "" {
			name = strings.TrimSuffix(name, ext)
		}
		exeName = name
	})

	return exeName
}
