package cvm

import (
	"fmt"

	"my-project-golang/internal/config"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

// TerminateInstances 销毁CVM竞价实例
func TerminateInstances(cfg *config.TencentCloudConfig, instanceId string) error {
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, _ := sdk.NewClient(credential, cfg.Region, cpf)

	request := sdk.NewTerminateInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{instanceId})

	response, err := client.TerminateInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("API错误: %s", err)
	}
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	fmt.Println("=== 销毁实例结果 ===")
	config.PrintJSON(response.ToJsonString())
	return nil
}
