package mocks

import data_transfer_api "github.com/r-mol/balanser_highload_system/protos"

type KeyValueServiceClient interface {
	New() (*KeyValueServiceClient, error)
	GetValue(request data_transfer_api.GetValueRequest) (*data_transfer_api.GetValueResponse, error)
	StoreValue(request data_transfer_api.StoreValueRequest) (*data_transfer_api.StoreValueResponse, error)
}
