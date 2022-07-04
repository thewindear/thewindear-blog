package tests

import (
    "github.com/thewindear/thewindear-blog/internal/app"
    "testing"
)

// TestApp 测试初始化应用
func TestApp(t *testing.T) {
    err := app.Run(TestAbsConfigFilePath)
    if err != nil {
        t.Fatal(err)
    }
    t.Logf("初始化应用成功: %v", app.Config)
}
