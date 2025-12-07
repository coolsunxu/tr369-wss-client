package model

import (
	"tr369-wss-client/pkg/api"
	tr181Model "tr369-wss-client/tr181/model"
)

// WSClient defines the interface for TR369 WebSocket client
type WSClient interface {
	// Connect establishes a WebSocket connection to the server
	Connect() error

	// Disconnect closes the WebSocket connection
	Disconnect()

	// StartMessageHandler starts the message handling goroutines
	StartMessageHandler()
}

// ClientRepository defines the interface for client data repository
type ClientRepository interface {

	// ConstructGetResp constructs a GET response from data map and paths
	ConstructGetResp(paths []string) api.Response_GetResp

	// HandleSetRequest handles a SET request by updating data map
	HandleSetRequest(path string, key string, value string)

	// IsExistPath checks if a path exists in the data map
	IsExistPath(path string) (isSuccess bool, nodePath string)

	// GetNewInstance gets a new instance path for a given path
	GetNewInstance(path string) (nodePath string)

	// HandleDeleteRequest handles a DELETE request by removing data
	HandleDeleteRequest(path string) (nodePath string, isFound bool)

	SaveData()

	StartClientRepository()

	AddListener(paramName string, listener tr181Model.Listener) error

	RemoveListener(paramName string) error
	NotifyListeners(paramName string, value interface{})
}

// ClientUseCase defines the interface for client use case
type ClientUseCase interface {
	// HandleMessage processes incoming USP messages
	HandleMessage(msg *api.Msg)
}
