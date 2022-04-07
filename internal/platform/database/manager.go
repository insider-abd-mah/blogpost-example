package database

import (
	"blog-example/log"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"sync"
	"time"
)

const databaseDefaultTimeout = time.Minute * 2

// DBConnectionInterface base connection interface
type DBConnectionInterface interface {
	Execute(callback func(sql *sql.DB) error) error
	CloseSession()
}

// DBConnection ...
type DBConnection struct {
	Name string
}

// mysqlDBManager mysql db connection pool manager
type mysqlDBManager struct {
	Connections map[string]connection
	Initialized bool
	sync.Mutex
}

// connection single db connection
type connection struct {
	Session    *sql.DB
	UsageCount int
	Timeout    time.Time
}

var activeMysqlDBManager = &mysqlDBManager{Connections: make(map[string]connection), Initialized: false}

// Init initialize db manager loop for closing non-used connections
func Init() {
	activeMysqlDBManager.Initialized = true

	go func() {
		for {
			for con, dbCon := range activeMysqlDBManager.Connections {
				if dbCon.UsageCount == 0 && dbCon.Timeout.Before(time.Now()) {
					activeMysqlDBManager.Lock()
					log.ErrChan <- dbCon.Session.Close()
					activeMysqlDBManager.Unlock()
					delete(activeMysqlDBManager.Connections, con)
				}
			}

			time.Sleep(time.Millisecond * 50)
		}
	}()
}

// Execute on the fly db execution
func (pc *DBConnection) Execute(callback func(sql *sql.DB) error) error {
	if !activeMysqlDBManager.Initialized {
		return fmt.Errorf("database.Init() not called for application")
	}

	con, err := activeMysqlDBManager.connectOrReuse(pc.Name)

	if err != nil {
		return err
	}

	con.Timeout = time.Now().Add(databaseDefaultTimeout)
	con.incrementUsage()
	err = callback(con.Session)
	con.decrementUsage()

	return err
}

// CloseSession manually close mysql session
func (pc *DBConnection) CloseSession() {
	con, ok := activeMysqlDBManager.Connections[pc.Name]

	if !ok {
		return
	}

	con.Timeout = time.Now()
	activeMysqlDBManager.Connections[pc.Name] = con
}

// connectOrReuse initialize new connection if connection not exists on pool
func (cp *mysqlDBManager) connectOrReuse(dbName string) (connection, error) {
	cp.Lock()
	defer cp.Unlock()

	con, ok := cp.Connections[dbName]

	if !ok || con.Session.Ping() != nil {
		session, err := cp.initNewSession(dbName)

		if err != nil {
			return con, fmt.Errorf("mysql %v connection not initialized, %v", dbName, err)
		}

		con = connection{
			Session:    session,
			UsageCount: 0,
			Timeout:    time.Now().Add(databaseDefaultTimeout),
		}

		cp.Connections[dbName] = con
	}

	return con, nil
}

func GetConnection() DBConnectionInterface {
	return &DBConnection{Name: "sample_db"}
}

// initNewSession create new connection
func (cp *mysqlDBManager) initNewSession(dbName string) (*sql.DB, error) {
	url := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASS") + "@" + os.Getenv("MYSQL_HOST")
	url = url + "/" + dbName + "?charset=utf8mb4&collation=utf8mb4_unicode_ci"

	return sql.Open("mysql", url)
}

// incrementUsage increase db usage
func (dbc *connection) incrementUsage() {
	dbc.UsageCount++
}

// decrementUsage decrease db usage
func (dbc *connection) decrementUsage() {
	dbc.UsageCount--
}
