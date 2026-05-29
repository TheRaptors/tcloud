package cvm

import (
	"fmt"

	"my-project-golang/internal/config"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

// RunInstances 创建CVM竞价实例，返回实例ID
func RunInstances(cfg *config.TencentCloudConfig) (string, error) {
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, _ := sdk.NewClient(credential, cfg.Region, cpf)

	request := sdk.NewRunInstancesRequest()

	request.InstanceChargeType = common.StringPtr("SPOTPAID")
	request.Placement = &sdk.Placement{
		Zone:      common.StringPtr(cfg.Zone),
		ProjectId: common.Int64Ptr(0),
	}
	request.InstanceType = common.StringPtr(cfg.InstanceType)
	request.ImageId = common.StringPtr(cfg.ImageId)
	request.SystemDisk = &sdk.SystemDisk{
		DiskType: common.StringPtr("CLOUD_BSSD"),
		DiskSize: common.Int64Ptr(50),
	}
	request.VirtualPrivateCloud = &sdk.VirtualPrivateCloud{
		VpcId:              common.StringPtr(cfg.VpcId),
		SubnetId:           common.StringPtr(cfg.SubnetId),
		AsVpcGateway:       common.BoolPtr(false),
		PrivateIpAddresses: common.StringPtrs([]string{cfg.PrivateIP}),
		Ipv6AddressCount:   common.Uint64Ptr(0),
	}
	request.InternetAccessible = &sdk.InternetAccessible{
		InternetChargeType:      common.StringPtr("TRAFFIC_POSTPAID_BY_HOUR"),
		InternetMaxBandwidthOut: common.Int64Ptr(200),
		PublicIpAssigned:        common.BoolPtr(true),
		InternetServiceProvider: common.StringPtr("BGP"),
	}
	request.InstanceCount = common.Int64Ptr(1)
	request.InstanceName = common.StringPtr(cfg.InstanceName)
	request.LoginSettings = &sdk.LoginSettings{
		KeyIds: common.StringPtrs([]string{cfg.KeyId}),
	}
	request.SecurityGroupIds = common.StringPtrs(cfg.SecurityGroupIds)
	request.EnhancedService = &sdk.EnhancedService{
		SecurityService: &sdk.RunSecurityServiceEnabled{
			Enabled: common.BoolPtr(true),
		},
		MonitorService: &sdk.RunMonitorServiceEnabled{
			Enabled: common.BoolPtr(true),
		},
		AutomationService: &sdk.RunAutomationServiceEnabled{
			Enabled: common.BoolPtr(false),
		},
	}
	request.InstanceMarketOptions = &sdk.InstanceMarketOptionsRequest{
		SpotOptions: &sdk.SpotMarketOptions{
			MaxPrice: common.StringPtr(cfg.MaxPrice),
		},
	}
	request.DisableApiTermination = common.BoolPtr(false)

	response, err := client.RunInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", fmt.Errorf("API错误: %s", err)
	}
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}

	fmt.Println("=== 创建实例结果 ===")
	config.PrintJSON(response.ToJsonString())

	// 返回第一个实例ID
	if len(response.Response.InstanceIdSet) > 0 {
		instanceId := *response.Response.InstanceIdSet[0]
		fmt.Printf("\n实例ID: %s\n", instanceId)
		return instanceId, nil
	}

	return "", fmt.Errorf("未获取到实例ID")
}
