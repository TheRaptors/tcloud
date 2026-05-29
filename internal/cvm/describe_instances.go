package cvm

import (
	"fmt"
	"time"

	"my-project-golang/internal/config"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

// DescribeInstances 查询实例信息，等待实例运行并返回公网IP
func DescribeInstances(cfg *config.TencentCloudConfig, instanceId string) (string, error) {
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, _ := sdk.NewClient(credential, cfg.Region, cpf)

	// 竞价实例创建后需要等待分配公网IP，轮询查询
	fmt.Println("等待实例分配公网IP...")
	maxRetries := 20
	for i := 0; i < maxRetries; i++ {
		request := sdk.NewDescribeInstancesRequest()
		request.InstanceIds = common.StringPtrs([]string{instanceId})

		response, err := client.DescribeInstances(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			return "", fmt.Errorf("API错误: %s", err)
		}
		if err != nil {
			return "", fmt.Errorf("请求失败: %w", err)
		}

		if len(response.Response.InstanceSet) > 0 {
			instance := response.Response.InstanceSet[0]

			// 检查实例状态是否为运行中
			if instance.InstanceState != nil && *instance.InstanceState == "RUNNING" {
				// 检查公网IP
				if len(instance.PublicIpAddresses) > 0 {
					publicIP := *instance.PublicIpAddresses[0]
					fmt.Printf("实例已运行，公网IP: %s\n", publicIP)
					return publicIP, nil
				}
			}

			state := "未知"
			if instance.InstanceState != nil {
				state = *instance.InstanceState
			}
			fmt.Printf("  第 %d/%d 次查询，实例状态: %s，等待中...\n", i+1, maxRetries, state)
		} else {
			fmt.Printf("  第 %d/%d 次查询，实例尚未就绪，等待中...\n", i+1, maxRetries)
		}

		time.Sleep(5 * time.Second)
	}

	return "", fmt.Errorf("等待超时：实例 %s 未能在规定时间内获取公网IP", instanceId)
}

// FindInstanceByPrivateIP 根据内网IP地址查找匹配的CVM实例，返回实例ID
func FindInstanceByPrivateIP(cfg *config.TencentCloudConfig) (string, error) {
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, _ := sdk.NewClient(credential, cfg.Region, cpf)

	request := sdk.NewDescribeInstancesRequest()

	response, err := client.DescribeInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", fmt.Errorf("API错误: %s", err)
	}
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}

	// 遍历所有实例，匹配内网IP
	for _, instance := range response.Response.InstanceSet {
		for _, ip := range instance.PrivateIpAddresses {
			if ip != nil && *ip == cfg.PrivateIP {
				instanceId := *instance.InstanceId
				state := "未知"
				if instance.InstanceState != nil {
					state = *instance.InstanceState
				}
				fmt.Printf("找到匹配内网IP %s 的实例: %s (状态: %s)\n", cfg.PrivateIP, instanceId, state)
				return instanceId, nil
			}
		}
	}

	return "", fmt.Errorf("未找到内网IP为 %s 的实例", cfg.PrivateIP)
}
