package balancer

import (
	"context"
	//"github.com/r-mol/balanser_highload_system/internal/proxy"

	//"google.golang.org/grpc"
	"testing"

	"github.com/golang/mock/gomock"
	Gmock_gen_proxy "github.com/r-mol/balanser_highload_system/internal/balancer/mock"
	mock_gen_client "github.com/r-mol/balanser_highload_system/internal/balancer/mock"
	data_transfer_api "github.com/r-mol/balanser_highload_system/protos"
	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProxy := Gmock_gen_proxy.NewMockProxy(ctrl)
	mockProxy.EXPECT().GetHost().Return("proxy1").AnyTimes()

	mockClient := mock_gen_client.NewMockKeyValueServiceClient(ctrl)
	//GetValueRequest := data_transfer_api.GetValueRequest{Key: "test"}
	mockClient.EXPECT().GetValue(data_transfer_api.GetValueRequest{Key: "test"}).Return(&data_transfer_api.GetValueResponse{}, nil)

	lb := &LoadBalancer{
		proxies: weightedProxiesBunch{&proxyWithWeight{Proxy: mockProxy, weight: 1}},
		Logger:  nil,
	}

	//origGrpcDialFn := grpcDialFn
	//grpcDialFn = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	//	return nil, nil
	//}
	//defer func() { grpcDialFn = origGrpcDialFn }()

	// Replace the data_transfer_api.NewKeyValueServiceClient function with a mock
	//origNewKeyValueServiceClientFn := newKeyValueServiceClientFn
	//newKeyValueServiceClientFn = func(conn *grpc.ClientConn) data_transfer_api.KeyValueServiceClient {
	//	return mockClient
	//}
	//defer func() { newKeyValueServiceClientFn = origNewKeyValueServiceClientFn }()

	// Test the GetValue method
	response, err := lb.GetValue(context.Background(), &data_transfer_api.GetValueRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestStoreValue(t *testing.T) {

}
