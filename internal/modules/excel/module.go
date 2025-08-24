package excel

import (
	"erp/internal/modules/excel/handler"
	"erp/internal/modules/excel/service"
)

// Module Excel模块
type Module struct {
	Handler    *handler.Handler
	Service    *service.Service
	GrpcClient *service.GrpcClient
}

// NewModule 创建Excel模块
func NewModule(grpcServerAddr string) (*Module, error) {
	// 创建gRPC客户端
	grpcClient, err := service.NewGrpcClient(grpcServerAddr)
	if err != nil {
		return nil, err
	}

	// 创建service
	svc := service.NewService(grpcClient)

	// 创建handler
	h := handler.NewHandler(svc)

	return &Module{
		Handler:    h,
		Service:    svc,
		GrpcClient: grpcClient,
	}, nil
}

// Close 关闭模块资源
func (m *Module) Close() error {
	if m.GrpcClient != nil {
		return m.GrpcClient.Close()
	}
	return nil
}
