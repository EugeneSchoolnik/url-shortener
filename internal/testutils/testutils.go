package testutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func LoadTestEnv(t *testing.T) {
	root := projectRoot()
	fmt.Println("root", root)
	_ = godotenv.Load(filepath.Join(root, ".env.test"))
}

func projectRoot() string {
	wd, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		wd = filepath.Dir(wd)
	}
}
