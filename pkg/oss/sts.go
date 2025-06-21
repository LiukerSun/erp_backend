package oss

import (
	"erp/config"
	"fmt"

	"github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// STSResponse STS临时凭证响应
type STSResponse struct {
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SecurityToken   string `json:"securityToken"`
	Expiration      string `json:"expiration"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
	Endpoint        string `json:"endpoint"`
}

// GetSTSCredentials 获取STS临时访问凭证 (用于前端直传)
func GetSTSCredentials() (*STSResponse, error) {
	// 检查必要的配置
	if config.AppConfig.OSSAccessKeyID == "" || config.AppConfig.OSSAccessKeySecret == "" {
		return nil, fmt.Errorf("OSS AccessKey配置缺失")
	}
	if config.AppConfig.OSSRoleARN == "" {
		return nil, fmt.Errorf("OSS_ROLE_ARN配置缺失，前端直传需要此配置")
	}

	// 创建STS客户端配置
	stsConfig := &client.Config{
		AccessKeyId:     tea.String(config.AppConfig.OSSAccessKeyID),
		AccessKeySecret: tea.String(config.AppConfig.OSSAccessKeySecret),
		Endpoint:        tea.String(getStsEndpoint(config.AppConfig.OSSRegion)),
	}

	// 创建STS客户端
	stsClient, err := sts20150401.NewClient(stsConfig)
	if err != nil {
		return nil, fmt.Errorf("创建STS客户端失败: %v", err)
	}

	// 准备AssumeRole请求
	assumeRoleRequest := &sts20150401.AssumeRoleRequest{
		RoleArn:         tea.String(config.AppConfig.OSSRoleARN),
		RoleSessionName: tea.String(config.AppConfig.OSSRoleSessionName),
		DurationSeconds: tea.Int64(3600), // 1小时有效期
	}

	// 发起AssumeRole请求
	runtime := &service.RuntimeOptions{}
	response, err := stsClient.AssumeRoleWithOptions(assumeRoleRequest, runtime)
	if err != nil {
		return nil, fmt.Errorf("获取STS临时凭证失败: %v", err)
	}

	// 构造返回结果
	credentials := response.Body.Credentials
	stsResponse := &STSResponse{
		AccessKeyID:     tea.StringValue(credentials.AccessKeyId),
		AccessKeySecret: tea.StringValue(credentials.AccessKeySecret),
		SecurityToken:   tea.StringValue(credentials.SecurityToken),
		Expiration:      tea.StringValue(credentials.Expiration),
		Region:          config.AppConfig.OSSRegion,
		Bucket:          config.AppConfig.OSSBucketName,
		Endpoint:        "https://" + config.AppConfig.GetOSSEndpoint(),
	}

	return stsResponse, nil
}

// getStsEndpoint 根据OSS区域获取对应的STS端点
func getStsEndpoint(region string) string {
	// 根据不同区域返回对应的STS端点
	switch region {
	case "cn-beijing":
		return "sts.cn-beijing.aliyuncs.com"
	case "cn-shanghai":
		return "sts.cn-shanghai.aliyuncs.com"
	case "cn-shenzhen":
		return "sts.cn-shenzhen.aliyuncs.com"
	case "cn-hangzhou":
		return "sts.cn-hangzhou.aliyuncs.com"
	case "cn-qingdao":
		return "sts.cn-qingdao.aliyuncs.com"
	case "cn-zhangjiakou":
		return "sts.cn-zhangjiakou.aliyuncs.com"
	case "cn-huhehaote":
		return "sts.cn-huhehaote.aliyuncs.com"
	case "cn-wulanchabu":
		return "sts.cn-wulanchabu.aliyuncs.com"
	case "cn-heyuan":
		return "sts.cn-heyuan.aliyuncs.com"
	case "cn-guangzhou":
		return "sts.cn-guangzhou.aliyuncs.com"
	case "cn-fuzhou":
		return "sts.cn-fuzhou.aliyuncs.com"
	case "cn-wuhan-lr":
		return "sts.cn-wuhan-lr.aliyuncs.com"
	case "cn-chengdu":
		return "sts.cn-chengdu.aliyuncs.com"
	case "cn-nanjing":
		return "sts.cn-nanjing.aliyuncs.com"
	// 海外区域
	case "ap-southeast-1":
		return "sts.ap-southeast-1.aliyuncs.com"
	case "ap-southeast-2":
		return "sts.ap-southeast-2.aliyuncs.com"
	case "ap-southeast-3":
		return "sts.ap-southeast-3.aliyuncs.com"
	case "ap-southeast-5":
		return "sts.ap-southeast-5.aliyuncs.com"
	case "ap-northeast-1":
		return "sts.ap-northeast-1.aliyuncs.com"
	case "ap-south-1":
		return "sts.ap-south-1.aliyuncs.com"
	case "us-east-1":
		return "sts.us-east-1.aliyuncs.com"
	case "us-west-1":
		return "sts.us-west-1.aliyuncs.com"
	case "eu-west-1":
		return "sts.eu-west-1.aliyuncs.com"
	case "eu-central-1":
		return "sts.eu-central-1.aliyuncs.com"
	case "me-east-1":
		return "sts.me-east-1.aliyuncs.com"
	default:
		// 默认使用杭州区域的STS端点
		return "sts.cn-hangzhou.aliyuncs.com"
	}
}
