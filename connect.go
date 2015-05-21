package GoSDK

//This file provides the interface for establishing connect collections
//that is to say, collections that are interfaced with non-platform databases
//they have to be treated a little bit differently, because a lot of configuration information
//needs to be trucked across the line during setup. enough that it's more helpful to have it in a
//struct than it is just in a map, or an endless list of function arguments.

type connectCollection interface {
	toMap() map[string]interface{}
	tableName() string
}

type MSSqlConfig struct {
	User, Password, Host, Port, DBName, Tablename string
}

func (ms MSSqlConfig) tableName() string { return ms.Tablename }

func (ms MSSqlConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["user"] = ms.User
	m["password"] = ms.Password
	m["address"] = ms.Host
	m["port"] = ms.Port
	m["dbname"] = ms.DBName
	m["tablename"] = ms.Tablename
	m["dbtype"] = "mssql"
	return m
}
