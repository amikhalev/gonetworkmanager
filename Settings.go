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

func NewSettings() (Settings, error) {
	var s settings
	return &s, s.init(NetworkManagerInterface, SettingsObjectPath)
}

type settings struct {
	dbusBase
}

func (s *settings) ListConnections() []Connection {
	var connectionPaths []dbus.ObjectPath

	s.call(&connectionPaths, SettingsListConnections)
	connections := make([]Connection, len(connectionPaths))

	var err error
	for i, path := range connectionPaths {
		connections[i], err = NewConnection(path)
		if err != nil {
			panic(err)
		}
	}

	return connections
}

func (s *settings) AddConnection(settings ConnectionSettings) (con Connection, err error) {
	var path dbus.ObjectPath
	err = s.callError(&path, SettingsAddConnection, settings)
	if err != nil {
		return
	}
	con, err = NewConnection(path)
	return
}
