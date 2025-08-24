package service

import (
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"regexp"
	"time"

	"erp/internal/modules/excel/model"
	"erp/pkg/proto"
)

// Service Excel服务
type Service struct {
	grpcClient *GrpcClient
}

// NewService 创建Excel服务
func NewService(grpcClient *GrpcClient) *Service {
	return &Service{
		grpcClient: grpcClient,
	}
}

// ParseExcel 解析Excel文件
func (s *Service) ParseExcel(ctx context.Context, file *multipart.FileHeader, sheetName string) (*model.ExcelParseResponse, error) {
	// 验证文件
	if file == nil {
		return nil, errors.New("请选择要上传的文件")
	}

	// 验证文件类型
	if !isExcelFile(file.Filename) {
		return nil, errors.New("只支持Excel文件格式(.xlsx, .xls)")
	}

	// 读取文件内容
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("文件读取失败")
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		return nil, errors.New("文件读取失败")
	}

	// 调用gRPC服务解析Excel
	grpcResp, err := s.grpcClient.ParseExcel(ctx, fileData, sheetName)
	if err != nil {
		return nil, errors.New("Excel解析失败: " + err.Error())
	}

	// 转换gRPC响应为我们的模型
	products := convertGrpcResponseToProducts(grpcResp)

	// 对product中的product_spec进行处理，解析 颜色分类 尺码大小 发货时效   示例 "product_spec": "{颜色分类:蓝色， 尺码大小:2XL， 发货时效:10天内发货}"

	// 使用更灵活的正则表达式提取并处理数据
	// 支持多种格式：
	// 1. {颜色分类:蓝色， 尺码大小:2XL， 发货时效:10天内发货}
	// 2. {颜色分类:黑+白， 码数:均码}
	// 3. 颜色分类:蓝色,尺码大小:2XL,发货时效:10天内发货

	// 尝试匹配包含发货时效的格式
	re1 := regexp.MustCompile(`颜色分类:(.*?)[，,]\s*尺码大小:(.*?)[，,]\s*发货时效:(.*?)}`)
	// 尝试匹配只有颜色和码数的格式
	re2 := regexp.MustCompile(`颜色分类:(.*?)[，,]\s*码数:(.*?)}`)

	for i := range products {
		log.Printf("处理商品: %s, ProductSpec: %s", products[i].ProductName, products[i].ProductSpec)

		// 转换DeliveryTime编码为可读文本
		products[i].DeliveryTime = convertDeliveryTime(products[i].DeliveryTime)

		// 先尝试匹配包含发货时效的格式
		matches := re1.FindStringSubmatch(products[i].ProductSpec)
		if len(matches) > 0 {
			products[i].Color = matches[1]
			products[i].Size = matches[2]
			products[i].ShippingTime = matches[3]
		} else {
			// 尝试匹配只有颜色和码数的格式
			matches = re2.FindStringSubmatch(products[i].ProductSpec)
			if len(matches) > 0 {
				products[i].Color = matches[1]
				products[i].Size = matches[2]
				products[i].ShippingTime = "" // 没有发货时效信息
			} else {
				log.Printf("解析失败 - 未找到匹配的格式")
			}
		}
	}

	return &model.ExcelParseResponse{
		Success:  true,
		Message:  "Excel解析成功",
		Products: products,
		Total:    len(products),
		UploadAt: time.Now(),
	}, nil
}

// isExcelFile 检查是否为Excel文件
func isExcelFile(filename string) bool {
	// 简单的文件扩展名检查
	// 实际项目中可能需要更严格的MIME类型检查
	return len(filename) > 4 && (filename[len(filename)-5:] == ".xlsx" || filename[len(filename)-4:] == ".xls")
}

// convertGrpcResponseToProducts 转换gRPC响应为商品列表
func convertGrpcResponseToProducts(grpcResp *proto.ParseExcelResponse) []model.ProductInfo {
	var products []model.ProductInfo

	if grpcResp == nil || !grpcResp.Success {
		return products
	}

	for _, p := range grpcResp.Products {
		products = append(products, model.ProductInfo{
			ProductID:       p.ProductId,
			ProductName:     p.ProductName,
			CategoryLevel1:  p.CategoryLevel1,
			CategoryLevel2:  p.CategoryLevel2,
			CategoryLevel3:  p.CategoryLevel3,
			CategoryLevel4:  p.CategoryLevel4,
			ProductType:     p.ProductType,
			ProductGroup:    p.ProductGroup,
			MerchantCode:    p.MerchantCode,
			MerchantSkuCode: p.MerchantSkuCode,
			SpecID:          p.SpecId,
			ProductSpec:     p.ProductSpec,
			DeliveryTime:    p.DeliveryTime,
			Price:           p.Price,
			InStock:         p.InStock,
			PreSaleStock:    p.PreSaleStock,
			TieredStock:     p.TieredStock,
			SalesVolume:     p.SalesVolume,
			CommissionRate:  p.CommissionRate,
			AuditStatus:     p.AuditStatus,
			ProductLink:     p.ProductLink,
			ProductCode:     p.ProductCode,
		})
	}

	return products
}

// convertDeliveryTime 转换发货时间编码为可读文本
func convertDeliveryTime(deliveryTime string) string {
	switch deliveryTime {
	case "1":
		return "次日发"
	case "2":
		return "48小时"
	case "3":
		return "当日发"
	default:
		return deliveryTime // 如果不是已知编码，返回原值
	}
}

// testRegexParse 测试正则表达式解析（仅用于调试）
