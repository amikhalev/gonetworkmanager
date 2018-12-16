package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

const (
	DeviceInterface = NetworkManagerInterface + ".Device"

	DevicePropertyInterface            = DeviceInterface + ".Interface"
	DevicePropertyIpInterface          = DeviceInterface + ".IpInterface"
	DevicePropertyState                = DeviceInterface + ".State"
	DevicePropertyIP4Config            = DeviceInterface + ".Ip4Config"
	DevicePropertyDeviceType           = DeviceInterface + ".DeviceType"
	DevicePropertyAvailableConnections = DeviceInterface + ".AvailableConnections"
	DevicePropertyDhcp4Config          = DeviceInterface + ".Dhcp4Config"
)

func DeviceFactory(conn *dbus.Conn, objectPath dbus.ObjectPath) Device {
	d := NewDevice(conn, objectPath).(*device)

	switch d.GetDeviceType() {
	case NmDeviceTypeWifi:
		return &wirelessDevice{*d}
	}

	return d
}

type Device interface {
	GetPath() dbus.ObjectPath

	// GetInterface gets the name of the device's control (and often data)
	// interface.
	GetInterface() string

	// GetIpInterface gets the IP interface name of the device.
	GetIpInterface() string

	// GetState gets the current state of the device.
	GetState() NmDeviceState

	// GetIP4Config gets the Ip4Config object describing the configuration of the
	// device. Only valid when the device is in the NM_DEVICE_STATE_ACTIVATED
	// state.
	GetIP4Config() IP4Config

	// GetDHCP4Config gets the Dhcp4Config object describing the configuration of the
	// device. Only valid when the device is in the NM_DEVICE_STATE_ACTIVATED
	// state.
	GetDHCP4Config() DHCP4Config

	// GetDeviceType gets the general type of the network device; ie Ethernet,
	// WiFi, etc.
	GetDeviceType() NmDeviceType

	// GetAvailableConnections gets an array of object paths of every configured
	// connection that is currently 'available' through this device.
	GetAvailableConnections() []Connection

	MarshalJSON() ([]byte, error)
}

func NewDevice(conn *dbus.Conn, objectPath dbus.ObjectPath) Device {
	var d = &device{}
	d.init(conn, NetworkManagerInterface, objectPath)
	return d
}

type device struct {
	dbusBase
}

func (d *device) GetPath() dbus.ObjectPath {
	return d.obj.Path()
}

func (d *device) GetInterface() string {
	return d.getStringProperty(DevicePropertyInterface)
}

func (d *device) GetIpInterface() string {
	return d.getStringProperty(DevicePropertyIpInterface)
}

func (d *device) GetState() NmDeviceState {
	return NmDeviceState(d.getUint32Property(DevicePropertyState))
}

func (d *device) GetIP4Config() IP4Config {
	path := d.getObjectProperty(DevicePropertyIP4Config)
	if path == "/" {
		return nil
	}

	cfg := NewIP4Config(d.conn, path)
	return cfg
}

func (d *device) GetDHCP4Config() DHCP4Config {
	path := d.getObjectProperty(DevicePropertyDhcp4Config)
	if path == "/" {
		return nil
	}

	cfg := NewDHCP4Config(d.conn, path)
	return cfg
}

func (d *device) GetDeviceType() NmDeviceType {
	return NmDeviceType(d.getUint32Property(DevicePropertyDeviceType))
}

func (d *device) GetAvailableConnections() []Connection {
	connPaths := d.getSliceObjectProperty(DevicePropertyAvailableConnections)
	conns := make([]Connection, len(connPaths))

	for i, path := range connPaths {
		conns[i] = NewConnection(d.conn, path)
	}

	return conns
}

func (d *device) marshalMap() map[string]interface{} {
	return map[string]interface{}{
		"Interface":            d.GetInterface(),
		"IP interface":         d.GetIpInterface(),
		"State":                d.GetState().String(),
		"IP4Config":            d.GetIP4Config(),
		"DHCP4Config":          d.GetDHCP4Config(),
		"DeviceType":           d.GetDeviceType().String(),
		"AvailableConnections": d.GetAvailableConnections(),
	}
}

func (d *device) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.marshalMap())
}
