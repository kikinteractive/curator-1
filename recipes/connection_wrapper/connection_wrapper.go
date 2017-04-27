package connectionwrapper

import (
	"errors"
	"fmt"

	curator "github.com/kikinteractive/curator-go"
	"github.com/tevino/abool"
)

// Logger provides customized logging within ConnectionWrapper.
type Logger interface {
	Printf(string, ...interface{})
	Debugf(string, ...interface{})
}

// DummyLogger is a Logger does nothing.
type DummyLogger struct{}

// Printf does nothing.
func (l DummyLogger) Printf(string, ...interface{}) {}

// Debugf does nothing.
func (l DummyLogger) Debugf(string, ...interface{}) {}

// ConnectionWrapperListenable represents a container of ConnectionWrapperListener(s).
type ConnectionWrapperListenable interface {
	curator.Listenable

	AddListener(ConnectionWrapperListener)
	RemoveListener(ConnectionWrapperListener)
}

// ConnectionWrapperListenerContainer is a container of ConnectionWrapperListener.
type ConnectionWrapperListenerContainer struct {
	curator.ListenerContainer
}

// AddListener adds a listener to the container.
func (c *ConnectionWrapperListenerContainer) AddListener(listener ConnectionWrapperListener) {
	c.Add(listener)
}

// RemoveListener removes a listener to the container.
func (c *ConnectionWrapperListenerContainer) RemoveListener(listener ConnectionWrapperListener) {
	c.Remove(listener)
}

// You can register a listener that will get notified when changes occur.
//
// NOTE: It's not possible to stay transactionally in sync. Users of this class must
// be prepared for false-positives and false-negatives. Additionally, always use the version number
// when updating data to avoid overwriting another process' change.
type ConnectionWrapper struct {
	// Tracks the number of outstanding background requests in flight. The first time this count reaches 0, we publish the initialized event.
	outstandingOps          uint64
	isInitialized           *abool.AtomicBool
	client                  curator.CuratorFramework
	listeners               ConnectionWrapperListenerContainer
	errorListeners          curator.UnhandledErrorListenerContainer
	state                   curator.State
	connectionStateListener curator.ConnectionStateListener
	logger                  Logger
}

// NewConnectionWrapper creates a ConnectionWrapper for the given client.
// Users of this type should provide a listener with callback functions that
// get call whenever the connection state changes.
func NewConnectionWrapper(client curator.CuratorFramework) *ConnectionWrapper {
	tc := &ConnectionWrapper{
		isInitialized: abool.New(),
		client:        client,
		state:         curator.LATENT,
		logger:        &DummyLogger{},
	}
	tc.connectionStateListener = curator.NewConnectionStateListener(
		func(client curator.CuratorFramework, newState curator.ConnectionState) {
			tc.handleStateChange(newState)
		})
	return tc
}

// Start starts the ConnectionWrapper.
// The cache is not started automatically. You must call this method.
func (tc *ConnectionWrapper) Start() error {
	if !tc.state.Change(curator.LATENT, curator.STARTED) {
		return errors.New("already started")
	}

	tc.client.ConnectionStateListenable().AddListener(tc.connectionStateListener)

	return nil
}

// SetLogger sets the inner Logger of ConnectionWrapper.
func (tc *ConnectionWrapper) SetLogger(l Logger) *ConnectionWrapper {
	tc.logger = l
	return tc
}

// Listenable returns the cache listeners.
func (tc *ConnectionWrapper) Listenable() ConnectionWrapperListenable {
	return &tc.listeners
}

// UnhandledErrorListenable allows catching unhandled errors in asynchornous operations.
func (tc *ConnectionWrapper) UnhandledErrorListenable() curator.UnhandledErrorListenable {
	return &tc.errorListeners
}

// callListeners calls all listeners with given event.
// Error is handled by handleException().
func (tc *ConnectionWrapper) callListeners(evt ConnectionWrapperEvent) {
	tc.listeners.ForEach(func(listener interface{}) {
		if err := listener.(ConnectionWrapperListener).ChildEvent(tc.client, evt); err != nil {
			tc.handleException(err)
		}
	})
}

// handleException sends an exception to any listeners, or else log the error if there are none.
func (tc *ConnectionWrapper) handleException(e error) {
	if tc.errorListeners.Len() == 0 {
		tc.logger.Printf("%s", e)
		return
	}
	tc.errorListeners.ForEach(func(listener interface{}) {
		listener.(curator.UnhandledErrorListener).UnhandledError(e)
	})
}

func (tc *ConnectionWrapper) handleStateChange(newState curator.ConnectionState) {
	fmt.Println("new state:", newState)
	switch newState {
	case curator.SUSPENDED:
		tc.logger.Debugf("state: Suspended")
		tc.publishEvent(ConnectionWrapperEventConnSuspended)
	case curator.LOST:
		tc.logger.Debugf("state: Lost")
		tc.isInitialized.UnSet()
		tc.publishEvent(ConnectionWrapperEventConnLost)
	case curator.CONNECTED:
		tc.logger.Debugf("state: Connected")
		tc.publishEvent(ConnectionWrapperEventConnConnected)
	case curator.RECONNECTED:
		tc.logger.Debugf("state: Reconnected")
		tc.publishEvent(ConnectionWrapperEventConnReconnected)
	default:
		tc.logger.Debugf("state: Unknown")
	}
}

// publishEvent publish an event with given type and data to all listeners.
func (tc *ConnectionWrapper) publishEvent(tp ConnectionWrapperEventType) {
	if tc.state.Value() != curator.STOPPED {
		evt := ConnectionWrapperEvent{Type: tp}
		tc.logger.Debugf("publishEvent: %v", evt)
		go tc.callListeners(evt)
	}
}
