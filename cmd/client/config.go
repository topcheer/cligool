package main

/*
 配置文件支持

配置文件查找优先级：
1. ./cligool.json (当前目录)
2. ~/.cligool.json (用户主目录)

如果都不存在，创建 ~/.cligool.json 并设置默认值
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 配置结构
type Config struct {
	Server string `json:"server"`
	Cols   int    `json:"cols"`
	Rows   int    `json:"rows"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Server: "https://cligool.zty8.cn",
		Cols:   0, // 0 表示自动检测
		Rows:   0, // 0 表示自动检测
	}
}

// LoadConfig 加载配置文件
// 查找优先级：./cligool.json -> ~/.cligool.json
// 如果都不存在，创建 ~/.cligool.json 并返回默认配置
func LoadConfig() (Config, string, error) {
	// 1. 尝试当前目录的配置文件
	localConfigPath := "cligool.json"
	if config, err := loadConfigFile(localConfigPath); err == nil {
		return config, localConfigPath, nil
	} else if !os.IsNotExist(err) {
		return Config{}, "", fmt.Errorf("读取 %s 失败: %w", localConfigPath, err)
	}

	// 2. 尝试用户主目录的配置文件
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, "", fmt.Errorf("获取用户主目录失败: %w", err)
	}

	homeConfigPath := filepath.Join(homeDir, ".cligool.json")
	if config, err := loadConfigFile(homeConfigPath); err == nil {
		return config, homeConfigPath, nil
	} else if !os.IsNotExist(err) {
		return Config{}, "", fmt.Errorf("读取 %s 失败: %w", homeConfigPath, err)
	}

	// 3. 配置文件不存在，创建默认配置
	defaultConfig := DefaultConfig()
	if err := saveConfigFile(homeConfigPath, defaultConfig); err != nil {
		return Config{}, "", fmt.Errorf("创建默认配置文件失败: %w", err)
	}

	return defaultConfig, homeConfigPath, nil
}

// loadConfigFile 从指定路径加载配置文件
func loadConfigFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return config, nil
}

// saveConfigFile 保存配置到指定路径
func saveConfigFile(path string, config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// MergeWithFlags 合并配置文件和命令行参数
// 命令行参数优先级高于配置文件
func MergeWithFlags(config Config, server *string, cols *int, rows *int) Config {
	result := config

	if server != nil && *server != "" {
		result.Server = *server
	}
	if cols != nil && *cols != 0 {
		result.Cols = *cols
	}
	if rows != nil && *rows != 0 {
		result.Rows = *rows
	}

	return result
}
