package oss

import (
	"erp/config"
	"fmt"
)

// STSCredentials STS临时凭证响应
type STSCredentials struct {
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SecurityToken   string `json:"securityToken"`
	Expiration      string `json:"expiration"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
	Endpoint        string `json:"endpoint"`
}

// InitOSS 初始化OSS服务（仅用于STS凭证服务）
func InitOSS() error {
	// 检查必要的配置
	if config.AppConfig.OSSAccessKeyID == "" || config.AppConfig.OSSAccessKeySecret == "" {
		return fmt.Errorf("OSS AccessKey配置缺失")
	}
	if config.AppConfig.OSSRoleARN == "" {
		return fmt.Errorf("OSS_ROLE_ARN配置缺失，前端直传需要此配置")
	}
	if config.AppConfig.OSSBucketName == "" {
		return fmt.Errorf("OSS Bucket名称配置缺失")
	}

	return nil
}
