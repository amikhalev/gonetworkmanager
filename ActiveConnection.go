package gonetworkmanager

import (
	"github.com/godbus/dbus"
)

const (
	ActiveConnectionInterface             = NetworkManagerInterface + ".Connection.Active"
	ActiveConnectionProperyConnection     = ActiveConnectionInterface + ".Connection"
	ActiveConnectionProperySpecificObject = ActiveConnectionInterface + ".SpecificObject"
	ActiveConnectionProperyID             = ActiveConnectionInterface + ".Id"
	ActiveConnectionProperyUUID           = ActiveConnectionInterface + ".Uuid"
	ActiveConnectionProperyType           = ActiveConnectionInterface + ".Type"
	ActiveConnectionProperyDevices        = ActiveConnectionInterface + ".Devices"
	ActiveConnectionProperyState          = ActiveConnectionInterface + ".State"
	ActiveConnectionProperyStateFlags     = ActiveConnectionInterface + ".StateFlags"
	ActiveConnectionProperyDefault        = ActiveConnectionInterface + ".Default"
	ActiveConnectionProperyIP4Config      = ActiveConnectionInterface + ".Ip4Config"
	ActiveConnectionProperyDHCP4Config    = ActiveConnectionInterface + ".Dhcp4Config"
	ActiveConnectionProperyDefault6       = ActiveConnectionInterface + ".Default6"
	ActiveConnectionProperyVPN            = ActiveConnectionInterface + ".Vpn"
	ActiveConnectionProperyMaster         = ActiveConnectionInterface + ".Master"
)

type ActiveConnection interface {
	// GetConnection gets connection object of the connection.
	GetConnection() Connection

	// GetSpecificObject gets a specific object associated with the active connection.
	GetSpecificObject() AccessPoint

	// GetID gets the ID of the connection.
	GetID() string

	// GetUUID gets the UUID of the connection.
	GetUUID() string

	// GetType gets the type of the connection.
	GetType() string

	// GetDevices gets array of device objects which are part of this active connection.
	GetDevices() []Device

	// GetState gets the state of the connection.
	GetState() uint32

	// GetStateFlags gets the state flags of the connection.
	GetStateFlags() uint32

	// GetDefault gets the default IPv4 flag of the connection.
	GetDefault() bool

	// GetIP4Config gets the IP4Config of the connection.
	GetIP4Config() IP4Config

	// GetDHCP4Config gets the DHCP4Config of the connection.
	GetDHCP4Config() DHCP4Config

	// GetVPN gets the VPN flag of the connection.
	GetVPN() bool

	// GetMaster gets the master device of the connection.
	GetMaster() Device
}

func NewActiveConnection(conn *dbus.Conn, objectPath dbus.ObjectPath) ActiveConnection {
	var a = &activeConnection{}
	a.init(conn, NetworkManagerInterface, objectPath)
	return a
}

type activeConnection struct {
	dbusBase
}

func (a *activeConnection) GetConnection() Connection {
	path := a.getObjectProperty(ActiveConnectionProperyConnection)
	con := NewConnection(a.conn, path)
	return con
}

func (a *activeConnection) GetSpecificObject() AccessPoint {
	path := a.getObjectProperty(ActiveConnectionProperySpecificObject)
	ap := NewAccessPoint(a.conn, path)
	return ap
}

func (a *activeConnection) GetID() string {
	return a.getStringProperty(ActiveConnectionProperyID)
}

func (a *activeConnection) GetUUID() string {
	return a.getStringProperty(ActiveConnectionProperyUUID)
}

func (a *activeConnection) GetType() string {
	return a.getStringProperty(ActiveConnectionProperyType)
}

func (a *activeConnection) GetDevices() []Device {
	paths := a.getSliceObjectProperty(ActiveConnectionProperyDevices)
	devices := make([]Device, len(paths))
	for i, path := range paths {
		devices[i] = DeviceFactory(a.conn, path)
	}
	return devices
}

func (a *activeConnection) GetState() uint32 {
	return a.getUint32Property(ActiveConnectionProperyState)
}

func (a *activeConnection) GetStateFlags() uint32 {
	return a.getUint32Property(ActiveConnectionProperyStateFlags)
}

func (a *activeConnection) GetDefault() bool {
	b := a.getProperty(ActiveConnectionProperyDefault)
	return b.(bool)
}

func (a *activeConnection) GetIP4Config() IP4Config {
	path := a.getObjectProperty(ActiveConnectionProperyIP4Config)
	r := NewIP4Config(a.conn, path)
	return r
}

func (a *activeConnection) GetDHCP4Config() DHCP4Config {
	path := a.getObjectProperty(ActiveConnectionProperyDHCP4Config)
	r := NewDHCP4Config(a.conn, path)
	return r
}

func (a *activeConnection) GetVPN() bool {
	ret := a.getProperty(ActiveConnectionProperyVPN)
	return ret.(bool)
}

func (a *activeConnection) GetMaster() Device {
	path := a.getObjectProperty(ActiveConnectionProperyMaster)
	r := DeviceFactory(a.conn, path)
	return r
}
