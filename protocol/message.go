package protocol

import (
	"encoding/json"
	"fmt"
	"time"
)

// MessageType defines the type of TR369 message
const (
	MessageTypeGet    = "get"
	MessageTypeSet    = "set"
	MessageTypeInform = "inform"
	MessageTypeAdd    = "add"
	MessageTypeDelete = "delete"
	MessageTypeReplace = "replace"
	MessageTypeOperation = "operation"
)

// ErrorCode defines the error codes for TR369 responses
const (
	ErrorCodeNoError              = 0
	ErrorCodeInvalidParameterName = 9001
	ErrorCodeInvalidParameterType = 9002
	ErrorCodeInvalidParameterValue = 9003
	ErrorCodeInternalError        = 9010
	ErrorCodeResourceBusy         = 9011
)

// Message represents the base TR369 message structure
type Message struct {
	MessageType string      `json:"msg_type"`
	MsgID       string      `json:"msg_id"`
	Timestamp   int64       `json:"timestamp,omitempty"`
	From        string      `json:"from,omitempty"`
	To          string      `json:"to,omitempty"`
	Body        interface{} `json:"body"`
}

// GetRequest represents a GET request message body
type GetRequest struct {
	Parameters []ParameterRef `json:"parameters"`
	Alias      string         `json:"alias,omitempty"`
}

// SetRequest represents a SET request message body
type SetRequest struct {
	Parameters []ParameterValue `json:"parameters"`
	Alias      string           `json:"alias,omitempty"`
}

// InformRequest represents an INFORM request message body
type InformRequest struct {
	Event      string           `json:"event"`
	Parameters []ParameterValue `json:"parameters,omitempty"`
	Alias      string           `json:"alias,omitempty"`
}

// ParameterRef represents a parameter reference (name only)
type ParameterRef struct {
	Name string `json:"name"`
}

// ParameterValue represents a parameter with its value
type ParameterValue struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// Error represents an error in a response
type Error struct {
	Parameter string `json:"parameter,omitempty"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
}

// GetResponse represents a GET response message body
type GetResponse struct {
	Parameters []ParameterValue `json:"parameters"`
	Errors     []Error          `json:"errors,omitempty"`
	Alias      string           `json:"alias,omitempty"`
}

// SetResponse represents a SET response message body
type SetResponse struct {
	Parameters []ParameterValue `json:"parameters"`
	Errors     []Error          `json:"errors,omitempty"`
	Alias      string           `json:"alias,omitempty"`
}

// InformResponse represents an INFORM response message body
type InformResponse struct {
	Status string `json:"status"`
	Alias  string `json:"alias,omitempty"`
}

// NewGetRequest creates a new GET request message
func NewGetRequest(parameters []string) *Message {
	paramRefs := make([]ParameterRef, len(parameters))
	for i, name := range parameters {
		paramRefs[i] = ParameterRef{Name: name}
	}

	return &Message{
		MessageType: MessageTypeGet,
		MsgID:       generateMsgID(),
		Timestamp:   time.Now().Unix(),
		Body: GetRequest{
			Parameters: paramRefs,
		},
	}
}

// NewSetRequest creates a new SET request message
func NewSetRequest(parameters []ParameterValue) *Message {
	return &Message{
		MessageType: MessageTypeSet,
		MsgID:       generateMsgID(),
		Timestamp:   time.Now().Unix(),
		Body: SetRequest{
			Parameters: parameters,
		},
	}
}

// NewInformMessage creates a new INFORM message
func NewInformMessage(event string, parameters []ParameterValue) *Message {
	return &Message{
		MessageType: MessageTypeInform,
		MsgID:       generateMsgID(),
		Timestamp:   time.Now().Unix(),
		Body: InformRequest{
			Event:      event,
			Parameters: parameters,
		},
	}
}

// NewGetResponse creates a response to a GET request
func NewGetResponse(request *Message) *GetResponse {
	return &GetResponse{
		Parameters: []ParameterValue{},
		Errors:     []Error{},
	}
}

// NewSetResponse creates a response to a SET request
func NewSetResponse(request *Message) *SetResponse {
	return &SetResponse{
		Parameters: []ParameterValue{},
		Errors:     []Error{},
	}
}

// NewInformResponse creates a response to an INFORM request
func NewInformResponse(request *Message) *InformResponse {
	return &InformResponse{
		Status: "OK",
	}
}

// generateMsgID generates a unique message ID
func generateMsgID() string {
	return fmt.Sprintf("msg-%d", time.Now().UnixNano())
}

// ParseMessage parses a JSON string into a Message structure
func ParseMessage(jsonStr string) (*Message, error) {
	var msg Message
	err := json.Unmarshal([]byte(jsonStr), &msg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}

	// 根据消息类型解析body
	bodyBytes, err := json.Marshal(msg.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	switch msg.MessageType {
	case MessageTypeGet:
		var body GetRequest
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return nil, fmt.Errorf("failed to parse GET body: %w", err)
		}
		msg.Body = body
	case MessageTypeSet:
		var body SetRequest
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return nil, fmt.Errorf("failed to parse SET body: %w", err)
		}
		msg.Body = body
	case MessageTypeInform:
		var body InformRequest
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return nil, fmt.Errorf("failed to parse INFORM body: %w", err)
		}
		msg.Body = body
	}

	return &msg, nil
}

// String returns a string representation of the message
func (m *Message) String() string {
	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Sprintf("{error: %v}", err)
	}
	return string(bytes)
}