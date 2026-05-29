package dnspod

import (
	"fmt"

	"my-project-golang/internal/config"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

// DescribeRecordList 获取DNS解析记录列表，返回第一条记录的 RecordId
func DescribeRecordList(cfg *config.TencentCloudConfig) (uint64, error) {
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, _ := sdk.NewClient(credential, "", cpf)

	request := sdk.NewDescribeRecordListRequest()
	request.Domain = common.StringPtr(cfg.Domain)
	request.Subdomain = common.StringPtr(cfg.Subdomain)

	response, err := client.DescribeRecordList(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return 0, fmt.Errorf("API错误: %s", err)
	}
	if err != nil {
		return 0, fmt.Errorf("请求失败: %w", err)
	}

	// 格式化输出完整响应
	fmt.Println("=== 记录列表 ===")
	config.PrintJSON(response.ToJsonString())

	// 提取第一条记录的 RecordId
	if len(response.Response.RecordList) > 0 {
		recordId := *response.Response.RecordList[0].RecordId
		fmt.Printf("\n获取到的 RecordId: %d\n", recordId)
		return recordId, nil
	}

	return 0, fmt.Errorf("未找到任何解析记录")
}
