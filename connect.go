package GoSDK

//This file provides the interface for establishing connect collections
//that is to say, collections that are interfaced with non-platform databases
//they have to be treated a little bit differently, because a lot of configuration information
//needs to be trucked across the line during setup. enough that it's more helpful to have it in a
//struct than it is just in a map, or an endless list of function arguments.

type connectCollection interface {
	toMap() map[string]interface{}
	tableName() string
	name() string
}

//MySqlConfig houses configuration information for an MySql-backed collection
type MySqlConfig struct {
	Name, User, Password, Host, Port, DBName, Tablename string
}

func (my MySqlConfig) tableName() string { return my.Tablename }
func (my MySqlConfig) name() string      { return my.Name }

func (my MySqlConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["name"] = my.Name
	m["user"] = my.User
	m["password"] = my.Password
	m["address"] = my.Host
	m["port"] = my.Port
	m["dbname"] = my.DBName
	m["tablename"] = my.Tablename
	m["dbtype"] = "mysql"
	return m
}

//MSSqlConfig houses configuration information for an MSSql-backed collection
type MSSqlConfig struct {
	User, Password, Host, Port, DBName, Tablename string
}

func (ms MSSqlConfig) tableName() string { return ms.Tablename }
func (ms MSSqlConfig) name() string      { return ms.Tablename }

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

//PostgresqlConfig houses configuration information for an Postgresql-backed collection
type PostgresqlConfig struct {
	User, Password, Host, Port, DBName, Tablename string
}

func (pg PostgresqlConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["user"] = pg.User
	m["password"] = pg.Password
	m["address"] = pg.Host
	m["port"] = pg.Port
	m["dbname"] = pg.DBName
	m["tablename"] = pg.Tablename
	m["dbtype"] = "postgres"
	return m
}

func (pg PostgresqlConfig) tableName() string { return pg.Tablename }
func (pg PostgresqlConfig) name() string      { return pg.Tablename }
