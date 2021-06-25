package gormdb

import (
	"fmt"
	"github.com/enorith/gormdb/plugins"
	"sync"

	orm "gorm.io/gorm"
)

var DefaultConnectionName = "default"

type Register func() (*orm.DB, error)

type Manager struct {
	connectionName string
	registers      map[string]Register
	connections    map[string]*orm.DB
	m              sync.RWMutex
}

func (m *Manager) Using(name string) *Manager {
	if m.connectionName != name {
		m.connectionName = name
	}

	return m
}

func (m *Manager) Register(name string, r Register) *Manager {
	m.m.Lock()
	defer m.m.Unlock()
	m.registers[name] = m.wrapRegister(r)
	return m
}

func (m *Manager) RegisterDefault(register Register) *Manager {
	return m.Register(DefaultConnectionName, register)
}

func (m *Manager) wrapRegister(r Register) Register {
	return func() (*orm.DB, error) {
		db, e := r()
		if e != nil {
			return nil, e
		}

		return db, db.Use(&plugins.Cast{})
	}
}

func (m *Manager) GetConnection(name ...string) (*orm.DB, error) {
	var using string
	if len(name) > 0 {
		using = name[0]
	} else {
		using = m.connectionName
	}
	m.Using(using)

	if len(m.connectionName) < 1 {
		m.connectionName = DefaultConnectionName
	}

	m.m.RLock()
	c, has := m.connections[m.connectionName]
	m.m.RUnlock()
	if has {
		return c, nil
	}
	m.m.RLock()
	register, exists := m.registers[m.connectionName]
	m.m.RUnlock()
	if !exists {
		return nil, fmt.Errorf("unregisterd connection [%s]", m.connectionName)
	}
	c, e := register()
	if e != nil {
		return nil, fmt.Errorf("register connection error: %v", e)
	}

	m.setConnection(m.connectionName, c)

	return c, nil
}

func (m *Manager) setConnection(name string, connection *orm.DB) {
	m.m.Lock()
	m.connections[name] = connection
	m.m.Unlock()
}

func NewManager() *Manager {
	return &Manager{
		registers:   make(map[string]Register),
		connections: make(map[string]*orm.DB),
		m:           sync.RWMutex{},
	}
}
