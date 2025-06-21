package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTExpireHours int
	ServerPort     string
	ServerMode     string
	// OSS配置
	OSSAccessKeyID     string
	OSSAccessKeySecret string
	OSSBucketName      string
	OSSRegion          string
	// STS配置
	OSSRoleARN         string
	OSSRoleSessionName string
}

var AppConfig *Config

func Init() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	AppConfig = &Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "password"),
		DBName:             getEnv("DB_NAME", "erp_db"),
		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		JWTExpireHours:     getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		ServerMode:         getEnv("SERVER_MODE", "debug"),
		OSSAccessKeyID:     getEnv("OSS_ACCESS_KEY_ID", ""),
		OSSAccessKeySecret: getEnv("OSS_ACCESS_KEY_SECRET", ""),
		OSSBucketName:      getEnv("OSS_BUCKET_NAME", ""),
		OSSRegion:          getEnv("OSS_REGION", "cn-beijing"),
		OSSRoleARN:         getEnv("OSS_ROLE_ARN", ""),
		OSSRoleSessionName: getEnv("OSS_ROLE_SESSION_NAME", "erp-frontend-upload"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetOSSEndpoint 根据区域获取OSS端点
func (c *Config) GetOSSEndpoint() string {
	return getOSSEndpointByRegion(c.OSSRegion)
}

// getOSSEndpointByRegion 根据区域获取OSS端点
func getOSSEndpointByRegion(region string) string {
	switch region {
	case "cn-beijing":
		return "oss-cn-beijing.aliyuncs.com"
	case "cn-shanghai":
		return "oss-cn-shanghai.aliyuncs.com"
	case "cn-hangzhou":
		return "oss-cn-hangzhou.aliyuncs.com"
	case "cn-shenzhen":
		return "oss-cn-shenzhen.aliyuncs.com"
	case "cn-qingdao":
		return "oss-cn-qingdao.aliyuncs.com"
	case "cn-zhangjiakou":
		return "oss-cn-zhangjiakou.aliyuncs.com"
	case "cn-huhehaote":
		return "oss-cn-huhehaote.aliyuncs.com"
	case "cn-wulanchabu":
		return "oss-cn-wulanchabu.aliyuncs.com"
	case "cn-heyuan":
		return "oss-cn-heyuan.aliyuncs.com"
	case "cn-guangzhou":
		return "oss-cn-guangzhou.aliyuncs.com"
	case "cn-fuzhou":
		return "oss-cn-fuzhou.aliyuncs.com"
	case "cn-wuhan-lr":
		return "oss-cn-wuhan-lr.aliyuncs.com"
	case "cn-chengdu":
		return "oss-cn-chengdu.aliyuncs.com"
	case "cn-nanjing":
		return "oss-cn-nanjing.aliyuncs.com"
	// 海外区域
	case "ap-southeast-1":
		return "oss-ap-southeast-1.aliyuncs.com"
	case "ap-southeast-2":
		return "oss-ap-southeast-2.aliyuncs.com"
	case "ap-southeast-3":
		return "oss-ap-southeast-3.aliyuncs.com"
	case "ap-southeast-5":
		return "oss-ap-southeast-5.aliyuncs.com"
	case "ap-northeast-1":
		return "oss-ap-northeast-1.aliyuncs.com"
	case "ap-south-1":
		return "oss-ap-south-1.aliyuncs.com"
	case "us-east-1":
		return "oss-us-east-1.aliyuncs.com"
	case "us-west-1":
		return "oss-us-west-1.aliyuncs.com"
	case "eu-west-1":
		return "oss-eu-west-1.aliyuncs.com"
	case "eu-central-1":
		return "oss-eu-central-1.aliyuncs.com"
	case "me-east-1":
		return "oss-me-east-1.aliyuncs.com"
	default:
		// 默认使用杭州区域
		return "oss-cn-hangzhou.aliyuncs.com"
	}
}
