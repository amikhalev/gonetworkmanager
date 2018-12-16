package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

const (
	ConnectionInterface = SettingsInterface + ".Connection"

	ConnectionGetSettings   = ConnectionInterface + ".GetSettings"
	ConnectionUpdate        = ConnectionInterface + ".Update"
	ConnectionUpdateUnsaved = ConnectionInterface + ".UpdateUnsaved"
	ConnectionUpdate2       = ConnectionInterface + ".Update2"
)

//type ConnectionSettings map[string]map[string]interface{}
type ConnectionSettings map[string]map[string]interface{}

type ConnectionUpdate2Args map[string]interface{}
type ConnectionUpdate2Flags uint32

const (
	UpdateToDisk           ConnectionUpdate2Flags = 0x01
	UpdateInMemory                                = 0x02
	UpdateInMemoryDetached                        = 0x04
	UpdateInMemoryOnly                            = 0x08
	UpdateVolatile                                = 0x10
	UpdateBlockAutoconnect                        = 0x20
)

type Connection interface {
	GetPath() dbus.ObjectPath

	// GetSettings gets the settings maps describing this network configuration.
	// This will never include any secrets required for connection to the
	// network, as those are often protected. Secrets must be requested
	// separately using the GetSecrets() call.
	GetSettings() ConnectionSettings

	Update(settings ConnectionSettings) error
	UpdateUnsaved(settings ConnectionSettings) error
	Update2(settings ConnectionSettings, flags ConnectionUpdate2Flags, args ConnectionUpdate2Args) error

	MarshalJSON() ([]byte, error)
}

func NewConnection(conn *dbus.Conn, objectPath dbus.ObjectPath) Connection {
	var c = &connection{}
	c.init(conn, NetworkManagerInterface, objectPath)
	return c
}

type connection struct {
	dbusBase
}

func (c *connection) GetPath() dbus.ObjectPath {
	return c.obj.Path()
}

func (c *connection) GetSettings() ConnectionSettings {
	var settings map[string]map[string]dbus.Variant
	c.call(&settings, ConnectionGetSettings)

	rv := make(ConnectionSettings)

	for k1, v1 := range settings {
		rv[k1] = make(map[string]interface{})

		for k2, v2 := range v1 {
			rv[k1][k2] = v2.Value()
		}
	}

	return rv
}

func (c *connection) Update(settings ConnectionSettings) (err error) {
	return c.obj.Call(ConnectionUpdate, 0, settings).Store()
}

func (c *connection) UpdateUnsaved(settings ConnectionSettings) (err error) {
	return c.obj.Call(ConnectionUpdateUnsaved, 0, settings).Store()
}

func (c *connection) Update2(settings ConnectionSettings, flags ConnectionUpdate2Flags, args ConnectionUpdate2Args) (err error) {
	return c.obj.Call(ConnectionUpdate2, 0, settings, (uint32)(flags), args).Store()
}

func (c *connection) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.GetSettings())
}
