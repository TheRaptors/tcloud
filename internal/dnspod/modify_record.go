package dnspod

import (
	"fmt"

	"my-project-golang/internal/config"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

// ModifyRecord 修改DNS解析记录
func ModifyRecord(cfg *config.TencentCloudConfig, recordId uint64, value string) error {
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, _ := sdk.NewClient(credential, "", cpf)

	request := sdk.NewModifyRecordRequest()
	request.Domain = common.StringPtr(cfg.Domain)
	request.RecordType = common.StringPtr("A")
	request.RecordLine = common.StringPtr("默认")
	request.Value = common.StringPtr(value)
	request.RecordId = common.Uint64Ptr(recordId)
	request.SubDomain = common.StringPtr(cfg.Subdomain)

	response, err := client.ModifyRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("API错误: %s", err)
	}
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	fmt.Println("=== 修改记录结果 ===")
	config.PrintJSON(response.ToJsonString())
	return nil
}
