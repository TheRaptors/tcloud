package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TencentCloudConfig 腾讯云配置结构体
type TencentCloudConfig struct {
	SecretID         string   `json:"secret_id"`
	SecretKey        string   `json:"secret_key"`
	Region           string   `json:"region"`
	Domain           string   `json:"domain"`
	Subdomain        string   `json:"subdomain"`
	PrivateIP        string   `json:"private_ip"`
	Zone             string   `json:"zone"`
	VpcId            string   `json:"vpc_id"`
	SubnetId         string   `json:"subnet_id"`
	SecurityGroupIds []string `json:"security_group_ids"`
	InstanceName     string   `json:"instance_name"`
	InstanceType     string   `json:"instance_type"`
	ImageId          string   `json:"image_id"`
	KeyId            string   `json:"key_id"`
	MaxPrice         string   `json:"max_price"`
}

// LoadConfig 从配置文件和环境变量加载腾讯云配置
// 优先级：环境变量 > JSON 配置文件
func LoadConfig() (*TencentCloudConfig, error) {
	var cfg TencentCloudConfig

	// 1. 尝试从配置文件加载（文件不存在时不报错，允许纯环境变量模式）
	configPath := findConfigPath()
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}

		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
	}

	// 2. 环境变量覆盖（环境变量优先级高于配置文件）
	applyEnvOverrides(&cfg)

	if cfg.SecretID == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("secret_id 或 secret_key 未配置（需通过配置文件或环境变量 TENCENTCLOUD_SECRET_ID / TENCENTCLOUD_SECRET_KEY 提供）")
	}

	return &cfg, nil
}

// findConfigPath 查找配置文件路径
func findConfigPath() string {
	// 获取当前可执行文件所在目录
	exePath, err := os.Executable()
	if err == nil {
		path := filepath.Join(filepath.Dir(exePath), "config", "tencentcloud.json")
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 尝试源码目录
	if _, err := os.Stat(filepath.Join("config", "tencentcloud.json")); err == nil {
		return filepath.Join("config", "tencentcloud.json")
	}

	return ""
}

// applyEnvOverrides 用环境变量覆盖配置（环境变量优先级更高）
func applyEnvOverrides(cfg *TencentCloudConfig) {
	// 字符串字段映射：环境变量名 → 配置字段指针
	stringFields := map[string]*string{
		"TENCENTCLOUD_SECRET_ID":     &cfg.SecretID,
		"TENCENTCLOUD_SECRET_KEY":    &cfg.SecretKey,
		"TENCENTCLOUD_REGION":        &cfg.Region,
		"TENCENTCLOUD_DOMAIN":        &cfg.Domain,
		"TENCENTCLOUD_SUBDOMAIN":     &cfg.Subdomain,
		"TENCENTCLOUD_PRIVATE_IP":    &cfg.PrivateIP,
		"TENCENTCLOUD_ZONE":          &cfg.Zone,
		"TENCENTCLOUD_VPC_ID":        &cfg.VpcId,
		"TENCENTCLOUD_SUBNET_ID":     &cfg.SubnetId,
		"TENCENTCLOUD_INSTANCE_NAME": &cfg.InstanceName,
		"TENCENTCLOUD_INSTANCE_TYPE": &cfg.InstanceType,
		"TENCENTCLOUD_IMAGE_ID":      &cfg.ImageId,
		"TENCENTCLOUD_KEY_ID":        &cfg.KeyId,
		"TENCENTCLOUD_MAX_PRICE":     &cfg.MaxPrice,
	}

	for envKey, field := range stringFields {
		if val := os.Getenv(envKey); val != "" {
			*field = val
		}
	}

	// 安全组列表：逗号分隔
	if val := os.Getenv("TENCENTCLOUD_SECURITY_GROUP_IDS"); val != "" {
		cfg.SecurityGroupIds = strings.Split(val, ",")
	}
}

// PrintJSON 格式化输出JSON
func PrintJSON(raw string) {
	var prettyJSON = new(bytes.Buffer)
	if err := json.Indent(prettyJSON, []byte(raw), "", "    "); err != nil {
		fmt.Printf("格式化JSON失败: %s\n", err)
		return
	}
	fmt.Printf("%s\n", prettyJSON.String())
}
