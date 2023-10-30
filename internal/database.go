package internal

import (
	support "ecommerce/utils/logger"
	"ecommerce/utils/sql"
	"fmt"
)

// NewDatabases initialises all the app's databases.
func NewDatabases(config Config, logger *support.Logger) map[string]sql.DBer {
	var res = make(map[string]sql.DBer)
	for k, v := range config.Database {
		dbModule := sql.NewDB(
			&sql.Config{
				DriverName:      v.Driver,
				Name:            k,
				URI:             v.Master,
				ReplicaURIs:     v.Slaves,
				MaxIdleConns:    v.MaxIdleConn,
				ConnMaxLifetime: v.MaxLifeTime,
				ConnMaxIdleTime: v.MaxIdleTime,
				MaxOpenConns:    v.MaxOpenConn,
			}, logger)

		errConn := dbModule.Connect()
		if errConn != nil {
			fmt.Println("error: ", errConn)
		}

		res[k] = dbModule
	}

	return res
}
