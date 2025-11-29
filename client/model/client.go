package model

import (
	"sync"
	"tr369-wss-client/pkg/api"
)

// ParameterChangeListener defines a callback for parameter changes
type ParameterChangeListener func(name string, oldValue, newValue interface{})

type TR181DataModel struct {
	Parameters map[string]interface{}
	Lock       sync.RWMutex
	Listeners  map[string][]ParameterChangeListener
}

// WSClient defines the interface for TR369 WebSocket client
type WSClient interface {
	// Connect establishes a WebSocket connection to the server
	Connect() error

	// Disconnect closes the WebSocket connection
	Disconnect()

	// StartMessageHandler starts the message handling goroutines
	StartMessageHandler()

	// HandleGetRequest handles incoming GET requests
	HandleGetRequest(inComingMsg *api.Msg)

	// HandleSetRequest handles incoming SET requests
	HandleSetRequest(inComingMsg *api.Msg)

	// HandleAddRequest handles incoming ADD requests
	HandleAddRequest(inComingMsg *api.Msg)

	// HandleDeleteRequest handles incoming DELETE requests
	HandleDeleteRequest(inComingMsg *api.Msg)

	// HandleOperateRequest handles incoming OPERATE requests
	HandleOperateRequest(inComingMsg *api.Msg)

	// SendOperateCompleteNotify sends an operation complete notification
	SendOperateCompleteNotify(objPath string, commandName string, commandKey string, outputArgs map[string]string)

	// HandleMTPMsgTransmit handles MTP message transmission
	HandleMTPMsgTransmit(msg *api.Msg) error
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
}
