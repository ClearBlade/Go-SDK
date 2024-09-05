package GoSDK

import (
	"fmt"
)

//This file provides the interface for establishing connect collections
//that is to say, collections that are interfaced with non-platform databases
//they have to be treated a little bit differently, because a lot of configuration information
//needs to be trucked across the line during setup. enough that it's more helpful to have it in a
//struct than it is just in a map, or an endless list of function arguments.

type ConnectCollection interface {
	ToMap() map[string]interface{}
	TableName() string
	Name() string
}

// MySqlConfig houses configuration information for an MySql-backed collection
type MySqlConfig struct {
	ColName, User, Password, Host, Port, DBName, Tablename string
}

func (my MySqlConfig) TableName() string { return my.Tablename }
func (my MySqlConfig) Name() string      { return my.ColName }

func (my MySqlConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["name"] = my.ColName
	m["user"] = my.User
	m["password"] = my.Password
	m["address"] = my.Host
	m["port"] = my.Port
	m["dbname"] = my.DBName
	m["tablename"] = my.Tablename
	m["dbtype"] = "mysql"
	return m
}

// MSSqlConfig houses configuration information for an MSSql-backed collection
type MSSqlConfig struct {
	ColName, User, Password, Host, Port, DBName, Tablename string
}

func (ms MSSqlConfig) TableName() string { return ms.Tablename }
func (ms MSSqlConfig) Name() string      { return ms.Tablename }

func (ms MSSqlConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["user"] = ms.User
	m["password"] = ms.Password
	m["address"] = ms.Host
	m["port"] = ms.Port
	m["dbname"] = ms.DBName
	m["tablename"] = ms.Tablename
	m["dbtype"] = "mssql"
	m["name"] = ms.ColName
	return m
}

// PostgresqlConfig houses configuration information for an Postgresql-backed collection
type PostgresqlConfig struct {
	ColName, User, Password, Host, Port, DBName, Tablename string
}

func (pg PostgresqlConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["user"] = pg.User
	m["password"] = pg.Password
	m["address"] = pg.Host
	m["port"] = pg.Port
	m["dbname"] = pg.DBName
	m["tablename"] = pg.Tablename
	m["dbtype"] = "postgres"
	m["name"] = pg.ColName
	return m
}

func (pg PostgresqlConfig) TableName() string { return pg.Tablename }
func (pg PostgresqlConfig) Name() string      { return pg.Tablename }

type MongoDBConfig struct {
	ColName, User, Password, Host, Port, DBName, Tablename string
}

func (mg MongoDBConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["user"] = mg.User
	m["password"] = mg.Password
	m["address"] = mg.Host
	m["port"] = mg.Port
	m["dbname"] = mg.DBName
	m["tablename"] = mg.Tablename
	m["dbtype"] = "MongoDB"
	m["name"] = mg.ColName
	return m
}

func (mg MongoDBConfig) TableName() string { return mg.Tablename }
func (mg MongoDBConfig) Name() string      { return mg.Tablename }

func GenerateConnectCollection(co map[string]interface{}) (ConnectCollection, error) {
	dbtype, ok := co["dbtype"].(string)
	if !ok {
		return nil, fmt.Errorf("generateConnectCollection: dbtype field missing or is not a string")
	}
	switch dbtype {
	case "mysql":
		return &MySqlConfig{
			User:      co["user"].(string),
			Password:  co["password"].(string),
			Host:      co["address"].(string),
			Port:      co["port"].(string),
			DBName:    co["dbname"].(string),
			Tablename: co["tablename"].(string),
			ColName:   co["name"].(string),
		}, nil
	case "mssql":
		return &MSSqlConfig{
			User:      co["user"].(string),
			Password:  co["password"].(string),
			Host:      co["address"].(string),
			Port:      co["port"].(string),
			DBName:    co["dbname"].(string),
			Tablename: co["tablename"].(string),
			ColName:   co["name"].(string),
		}, nil
	case "postgresql":
		return &PostgresqlConfig{
			User:      co["user"].(string),
			Password:  co["password"].(string),
			Host:      co["address"].(string),
			Port:      co["port"].(string),
			DBName:    co["dbname"].(string),
			Tablename: co["tablename"].(string),
			ColName:   co["name"].(string),
		}, nil
	case "MongoDB":
		return &MongoDBConfig{
			User:      co["user"].(string),
			Password:  co["password"].(string),
			Host:      co["address"].(string),
			Port:      co["port"].(string),
			DBName:    co["dbname"].(string),
			Tablename: co["tablename"].(string),
			ColName:   co["name"].(string),
		}, nil
	default:
		return nil, fmt.Errorf("generateConnectCollection: Unknown connect database type: '%s'\n", dbtype)
	}
}
