package test

import (
	"blog-example/log"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
	"testing"
	"time"
)

const DbName = "sample_db"

// NewIntegration will generate a
// Parameters
//   t: *testing.T
//
//   This function creates a new MySQL container for us and return three value;
//
//     1) Database connection to the newly created MySQL container
//     2) Tear down function, which we will truncate all the database after
//        each unit test case work also closes the connection to the test database
func NewIntegration(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	c := StartContainer(t)
	log.ErrChan = make(chan error, 100)
	_ = os.Setenv("BASE_PATH", getBasePath())
	_ = godotenv.Load(getTestsPath() + ".env")
	host := "tcp(" + c.Host + ")"
	url := "root:root@tcp(" + c.Host + ")"
	url = "root:root@tcp(" + c.Host + ")" + "/" + DbName + "?charset=utf8mb4&collation=utf8mb4_unicode_ci"
	db, err := sql.Open("mysql", url)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = os.Setenv("MYSQL_HOST", host)
	fmt.Println("from integration test package" + os.Getenv("MYSQL_HOST"))

	if err != nil {
		t.Errorf(err.Error())
	}

	if err != nil {
		t.Fatalf("failed to opening database connection to information_schema: %v", err)
	}

	healthCheck(db, t, c)

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		_ = Truncate(db, "sample_db")
	}

	return db, teardown
}

// Wait for the database to be ready. Wait 100ms longer between each attempt.
// Do not try more than 20 times.
func healthCheck(db *sql.DB, t *testing.T, c *Container) {
	var pingError error

	maxAttempts := 20

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = db.Ping()

		if pingError == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		DumpContainerLogs(t, c)
		t.Fatalf("waiting for database to be ready: %v", pingError)
	}
}
