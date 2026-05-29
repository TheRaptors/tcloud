package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// LoadConfig 从配置文件加载腾讯云配置
func LoadConfig() (*TencentCloudConfig, error) {
	// 获取当前可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	configPath := filepath.Join(filepath.Dir(exePath), "config", "tencentcloud.json")

	// 如果可执行文件目录下没有配置文件，尝试使用源码目录
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join("config", "tencentcloud.json")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg TencentCloudConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if cfg.SecretID == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("配置文件中 secret_id 或 secret_key 为空")
	}

	return &cfg, nil
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
