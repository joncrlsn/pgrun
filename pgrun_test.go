package main

import (
	"fmt"
	"github.com/joncrlsn/fileutil"
	"github.com/joncrlsn/pgutil"
	"github.com/stvp/assert"
	"strings"
	"testing"
)

var testFileName string

const expectedSqlStatements int = 4

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

-- STATEMENT-BEGIN
CREATE OR REPLACE FUNCTION increment(i integer) RETURNS integer AS $$
BEGIN
    RETURN i + 1;
END;
$$ 
LANGUAGE plpgsql;
-- STATEMENT-END

INSERT INTO t_user (username, email)
 VALUES ('Jim Bob', 'jim@bob.com')  ;

 `, "\n")

	fileutil.WriteLinesSlice(lines, testFileName)
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

	assert.Equal(t, expectedSqlStatements, counter, "Incorrect number of SQL statements parsed")
}

// Runs the file against a test database
func xxxTest_runFile(t *testing.T) {
	dbInfo := pgutil.DbInfo{}
	dbInfo.DbName = "dev-cpc"
	dbInfo.DbUser = "c42"
	dbInfo.DbPass = ""
	dbInfo.DbHost = "localhost"
	dbInfo.DbOptions = "sslmode=disable"
	runFile(testFileName, &dbInfo)
}
