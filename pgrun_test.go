package main

import "strings"
import "testing"
import "github.com/joncrlsn/fileutil"
import "github.com/joncrlsn/pgutil"
import "fmt"

var testFileName string

const expectedSqlStatements int = 3

// Creates the testing file we'll be using
func init() {
	//testFileName = fileutil.TempFileName("sqlrun.", ".sql")
    testFileName = "temp.sql"
	lines := strings.Split(`

INSERT INTO alert (id, name) VALUES (1, 'Your zipper is open');
-- SQL comment 2
UPDATE t_user SET username = 'Sloppy Joe'
   WHERE username = 'Messy Jose;' -- comment at the end
   AND email = 'jose@sloppy.com'; -- comment at the end
INSERT INTO t_user (username, email)
 VALUES ('Jim Bob', 'jim@bob.com')  ;`, "\n")

	fileutil.WriteLinesArray(lines, testFileName)
}

// Tests that SQL statements are kept together as they should
// without getting two statements in one.
func Test_StatementGrouping(t *testing.T) {

	counter := 0
	sqlChan := sqlStatements(testFileName)
	for sql := range sqlChan {
		counter++
		fmt.Printf("=== Statement:\n%s\n", sql)
	}

	if counter == expectedSqlStatements {
		t.Log("one test passed.")
	} else {
		t.Error("Incorrect number of SQL statements found: %d instead of %d", counter, expectedSqlStatements)
	}
}

// Runs the file against a test database
func Test_runFile(t *testing.T) {
    dbInfo := pgutil.DbInfo{}
    dbInfo.DbName = "dev-cpc"
    dbInfo.DbUser = "c42"
    dbInfo.DbPass = ""
    dbInfo.DbHost = "localhost"
    dbInfo.DbOptions = "sslmode=disable"
    runFile(testFileName, &dbInfo)
}
