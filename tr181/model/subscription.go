package model

const (
	ValueChange       = "ValueChange"
	ObjectCreation    = "ObjectCreation"
	ObjectDeletion    = "ObjectDeletion"
	OperationComplete = "OperationComplete"
	Event             = "Event"
)

// Listener defines a callback for changes
type Listener func(interface{})

type TR181DataModel struct {
	Parameters map[string]interface{}
	Listeners  map[string][]Listener
}
