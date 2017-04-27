package connectionwrapper

import "fmt"

// ConnectionWrapperEventType represents the type of change to a path
type ConnectionWrapperEventType int

const (
	// ConnectionWrapperEventConnSuspended is called when the connection has changed to SUSPENDED
	ConnectionWrapperEventConnSuspended ConnectionWrapperEventType = iota
	// ConnectionWrapperEventConnReconnected is called when the connection has changed to RECONNECTED
	ConnectionWrapperEventConnReconnected
	// ConnectionWrapperEventConnLost is called when the connection has changed to LOST
	ConnectionWrapperEventConnLost
	// ConnectionWrapperEventInitialized is posted after the initial cache has been fully populated
	ConnectionWrapperEventInitialized

	ConnectionWrapperEventConnConnected
)

// String returns the string representation of ConnectionWrapperEventType
// "Unknown" is returned when event type is unknown
func (et ConnectionWrapperEventType) String() string {
	switch et {
	case ConnectionWrapperEventConnSuspended:
		return "ConnSuspended"
	case ConnectionWrapperEventConnReconnected:
		return "ConnReconnected"
	case ConnectionWrapperEventConnLost:
		return "ConnLost"
	case ConnectionWrapperEventInitialized:
		return "Initialized"
	case ConnectionWrapperEventConnConnected:
		return "Connected"
	default:
		return "Unknown"
	}
}

// ConnectionWrapperEvent represents a change to a path
type ConnectionWrapperEvent struct {
	Type ConnectionWrapperEventType
}

// String returns the string representation of ConnectionWrapperEvent
func (e ConnectionWrapperEvent) String() string {
	var path string
	return fmt.Sprintf("ConnectionWrapperEvent{%s %s}", e.Type, path)
}
