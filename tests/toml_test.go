package tests

import (
    "github.com/BurntSushi/toml"
    "github.com/thewindear/thewindear-blog/internal/app"
    "testing"
)

// TestUseTomlParseFile 测试解析toml配置文件
func TestUseTomlParseFile(t *testing.T) {
    configPath := "../config/config-dev.toml"
    var config app.Conf
    _, err := toml.DecodeFile(configPath, &config)
    if err != nil {
        t.Errorf("parse toml file error: %s", err.Error())
    }
    t.Log(config.Database)
}
