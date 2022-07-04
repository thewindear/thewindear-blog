package tests

import (
    "github.com/thewindear/thewindear-blog/internal/app"
    "log"
    "path/filepath"
    "runtime"
)

const (
    configPath = "config/config-dev.toml"
)

var TestAbsRootDir string
var TestAbsConfigFilePath string

func init() {
    _, fullFilename, _, _ := runtime.Caller(0)
    TestAbsRootDir = filepath.Dir(filepath.Dir(fullFilename))
    TestAbsConfigFilePath = filepath.Join(TestAbsRootDir, configPath)
}

func InitApp() {
    err := app.Init(TestAbsConfigFilePath)
    if err != nil {
        log.Fatalln(err)
    }
}
