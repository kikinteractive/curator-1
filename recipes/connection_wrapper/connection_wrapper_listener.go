package connectionwrapper

import curator "github.com/kikinteractive/curator-go"

// ConnectionWrapperListener represents listener for ConnectionWrapper changes
type ConnectionWrapperListener interface {
	// Called when a change has occurred
	ChildEvent(client curator.CuratorFramework, event ConnectionWrapperEvent) error
}

// childEventCallback is the callback type of ChildEvent within ConnectionWrapperListener
type childEventCallback func(curator.CuratorFramework, ConnectionWrapperEvent) error

// ConnectionWrapperListenerPrototype is the internal implementation of ConnectionWrapperListener
type ConnectionWrapperListenerPrototype struct {
	childEvent childEventCallback
}

// ChildEvent is called when a change has occurred
func (l *ConnectionWrapperListenerPrototype) ChildEvent(client curator.CuratorFramework, event ConnectionWrapperEvent) error {
	return l.childEvent(client, event)
}

// NewConnectionWrapperListener creates ConnectionWrapperListener with given function
func NewConnectionWrapperListener(cb childEventCallback) ConnectionWrapperListener {
	return &ConnectionWrapperListenerPrototype{cb}
}
