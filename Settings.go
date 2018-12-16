package gonetworkmanager

import (
	"github.com/godbus/dbus"
)

const (
	SettingsInterface  = NetworkManagerInterface + ".Settings"
	SettingsObjectPath = NetworkManagerObjectPath + "/Settings"

	SettingsListConnections = SettingsInterface + ".ListConnections"
	SettingsAddConnection   = SettingsInterface + ".AddConnection"
)

type Settings interface {

	// ListConnections gets list the saved network connections known to NetworkManager
	ListConnections() []Connection

	// AddConnection call new connection and save it to disk.
	AddConnection(settings ConnectionSettings) (Connection, error)
}

func NewSettings(conn *dbus.Conn) Settings {
	var s = &settings{}
	s.init(conn, NetworkManagerInterface, SettingsObjectPath)
	return s
}

func NewSettingsSystem() (Settings, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return NewSettings(conn), nil
}

type settings struct {
	dbusBase
}

func (s *settings) ListConnections() []Connection {
	var connectionPaths []dbus.ObjectPath

	s.call(&connectionPaths, SettingsListConnections)
	connections := make([]Connection, len(connectionPaths))

	for i, path := range connectionPaths {
		connections[i] = NewConnection(s.conn, path)
	}

	return connections
}

func (s *settings) AddConnection(settings ConnectionSettings) (con Connection, err error) {
	var path dbus.ObjectPath
	err = s.callError(&path, SettingsAddConnection, settings)
	if err != nil {
		return
	}
	con = NewConnection(s.conn, path)
	return
}
