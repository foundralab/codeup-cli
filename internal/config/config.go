// Package config 负责 codeup CLI 的配置加载与持久化。
//
// 配置优先级（高 -> 低）：命令行参数 > 环境变量 > 配置文件。
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// 环境变量名
const (
	EnvToken  = "CODEUP_TOKEN"
	EnvDomain = "CODEUP_DOMAIN"
	EnvOrgID  = "CODEUP_ORG_ID"
)

// Config 是 CLI 的全局配置。
type Config struct {
	// Domain 云效服务接入点，例如 openapi-rdc.aliyuncs.com（中心版）。
	Domain string `json:"domain,omitempty"`
	// Token 个人访问令牌（x-yunxiao-token）。
	Token string `json:"token,omitempty"`
	// OrganizationID 组织 ID。仅中心版需要；留空则按 Region 版调用。
	OrganizationID string `json:"organizationId,omitempty"`
}

// Path 返回配置文件路径 ~/.config/codeup/config.json。
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户主目录: %w", err)
	}
	return filepath.Join(home, ".config", "codeup", "config.json"), nil
}

// LoadFile 仅读取配置文件，不叠加环境变量。配置文件不存在时返回空配置。
func LoadFile() (*Config, error) {
	cfg := &Config{}

	path, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("解析配置文件 %s 失败: %w", path, err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
	}
	return cfg, nil
}

// Load 读取配置文件并叠加环境变量。配置文件不存在时不报错。
func Load() (*Config, error) {
	cfg, err := LoadFile()
	if err != nil {
		return nil, err
	}

	if v := os.Getenv(EnvDomain); v != "" {
		cfg.Domain = v
	}
	if v := os.Getenv(EnvToken); v != "" {
		cfg.Token = v
	}
	if v := os.Getenv(EnvOrgID); v != "" {
		cfg.OrganizationID = v
	}
	return cfg, nil
}

// Save 将配置写入配置文件（0600，因为包含令牌）。
func (c *Config) Save() error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	return nil
}
