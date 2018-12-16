package gonetworkmanager

import "github.com/godbus/dbus"

type NmObject interface {
	GetPath() dbus.ObjectPath
	GetDbusConnection() *dbus.Conn
	Close() error
}
