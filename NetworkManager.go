package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

const (
	NetworkManagerInterface  = "org.freedesktop.NetworkManager"
	NetworkManagerObjectPath = "/org/freedesktop/NetworkManager"

	NetworkManagerGetDevices               = NetworkManagerInterface + ".GetDevices"
	NetworkManagerActivateConnection       = NetworkManagerInterface + ".ActivateConnection"
	NetworkManagerPropertyState            = NetworkManagerInterface + ".State"
	NetworkManagerPropertyActiveConnection = NetworkManagerInterface + ".ActiveConnections"
)

type NetworkManager interface {

	// GetDevices gets the list of network devices.
	GetDevices() []Device

	// GetState returns the overall networking state as determined by the
	// NetworkManager daemon, based on the state of network devices under it's
	// management.
	GetState() NmState

	// GetActiveConnections returns the active connection of network devices.
	GetActiveConnections() []ActiveConnection

	// ActivateWirelessConnection requests activating access point to network device
	ActivateWirelessConnection(connection Connection, device Device, accessPoint AccessPoint) ActiveConnection

	Subscribe() <-chan *dbus.Signal
	Unsubscribe()

	MarshalJSON() ([]byte, error)
}

func NewNetworkManager(dbusConn *dbus.Conn) NetworkManager {
	var nm = &networkManager{}
	nm.init(dbusConn, NetworkManagerInterface, NetworkManagerObjectPath)
	return nm
}

func NewNetworkManagerSystem() (NetworkManager, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return NewNetworkManager(conn), nil
}

type networkManager struct {
	dbusBase

	sigChan chan *dbus.Signal
}

func (n *networkManager) GetDevices() []Device {
	var devicePaths []dbus.ObjectPath

	n.call(&devicePaths, NetworkManagerGetDevices)
	devices := make([]Device, len(devicePaths))

	for i, path := range devicePaths {
		devices[i] = DeviceFactory(n.conn, path)
	}

	return devices
}

func (n *networkManager) GetState() NmState {
	return NmState(n.getUint32Property(NetworkManagerPropertyState))
}

func (n *networkManager) GetActiveConnections() []ActiveConnection {
	acPaths := n.getSliceObjectProperty(NetworkManagerPropertyActiveConnection)
	ac := make([]ActiveConnection, len(acPaths))

	for i, path := range acPaths {
		ac[i] = NewActiveConnection(n.conn, path)
	}

	return ac
}

func (n *networkManager) ActivateWirelessConnection(c Connection, d Device, ap AccessPoint) ActiveConnection {
	var opath dbus.ObjectPath
	n.call(&opath, NetworkManagerActivateConnection, c.GetPath(), d.GetPath(), ap.GetPath())
	return nil
}

func (n *networkManager) Subscribe() <-chan *dbus.Signal {
	if n.sigChan != nil {
		return n.sigChan
	}

	n.subscribeNamespace(NetworkManagerObjectPath)
	n.sigChan = make(chan *dbus.Signal, 10)
	n.conn.Signal(n.sigChan)

	return n.sigChan
}

func (n *networkManager) Unsubscribe() {
	n.conn.RemoveSignal(n.sigChan)
	n.sigChan = nil
}

func (n *networkManager) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"NetworkState": n.GetState().String(),
		"Devices":      n.GetDevices(),
	})
}
