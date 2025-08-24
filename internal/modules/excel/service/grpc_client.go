package service

import (
	"context"
	"erp/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClient gRPC客户端
type GrpcClient struct {
	conn   *grpc.ClientConn
	client proto.ExcelParseServiceClient
}

// NewGrpcClient 创建gRPC客户端
func NewGrpcClient(serverAddr string) (*GrpcClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// 创建gRPC客户端
	client := proto.NewExcelParseServiceClient(conn)

	return &GrpcClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close 关闭连接
func (gc *GrpcClient) Close() error {
	return gc.conn.Close()
}

// ParseExcel 解析Excel文件
func (gc *GrpcClient) ParseExcel(ctx context.Context, excelData []byte, sheetName string) (*proto.ParseExcelResponse, error) {
	// 创建请求
	req := &proto.ParseExcelRequest{
		ExcelData: excelData,
		SheetName: sheetName,
	}

	// 调用gRPC服务
	resp, err := gc.client.ParseExcel(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
