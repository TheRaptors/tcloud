package dnspod

import (
	"fmt"

	"my-project-golang/internal/config"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

// DescribeRecord 获取单条DNS解析记录详情
func DescribeRecord(cfg *config.TencentCloudConfig, recordId uint64) error {
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, _ := sdk.NewClient(credential, "", cpf)

	request := sdk.NewDescribeRecordRequest()
	request.Domain = common.StringPtr(cfg.Domain)
	request.RecordId = common.Uint64Ptr(recordId)

	response, err := client.DescribeRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("API错误: %s", err)
	}
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	fmt.Println("=== 记录详情 ===")
	config.PrintJSON(response.ToJsonString())
	return nil
}
