package main

//
// Copyright (c) 2016 Jon Carlson.  All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.
//

import (
	"bufio"
	"fmt"
	"github.com/joncrlsn/fileutil"
	"github.com/joncrlsn/misc"
	"github.com/joncrlsn/pgutil"
	flag "github.com/ogier/pflag"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	version = "1.0.9"
)

var (

	// There are times when we need to mark the statement beginning and end
	// with something other than a semi-colon.  For example, functions can
	// have semi-colons at the end of a line without it signifying the end
	// of the statement.
	stmtBeginRegex = regexp.MustCompile(`^\s*--\s*STATEMENT-BEGIN\s*$`)
	stmtEndRegex   = regexp.MustCompile(`^\s*--\s*STATEMENT-END\s*$`)

	// eosRegex is End of Statement regex
	eosRegex = regexp.MustCompile(`;\s*$|;\s*--.*$`)

	inReader = bufio.NewReader(os.Stdin)
)

// Executes a file of SQL statements one statement at a time, stopping everything
// if one of them has an error
func main() {

	// -f (filename) is a required program argument
	var fileName = flag.StringP("file", "f", "", "path of the SQL file to run")
	dbInfo := pgutil.DbInfo{}
	verFlag, helpFlag := dbInfo.Populate()

	if verFlag {
		fmt.Fprintf(os.Stderr, "%s version %s\n", os.Args[0], version)
		fmt.Fprintln(os.Stderr, "Copyright (c) 2016 Jon Carlson.  All rights reserved.")
		fmt.Fprintln(os.Stderr, "Use of this source code is governed by the MIT license")
		fmt.Fprintln(os.Stderr, "that can be found here: http://opensource.org/licenses/MIT")
		os.Exit(1)
	}

	if helpFlag {
		usage()
	}

	if len(*fileName) == 0 {
		fmt.Fprintln(os.Stderr, "Missing required filename argument (-f)")
		usage()
	}

	exists, _ := fileutil.Exists(*fileName)
	if !exists {
		fmt.Fprintf(os.Stderr, "File does not exist: %s\n", *fileName)
		os.Exit(2)
	}

	runFile(*fileName, &dbInfo)
}

// Reads the file and runs the SQL statements one by one
func runFile(fileName string, dbInfo *pgutil.DbInfo) {
	// Open connection to the database
	db, err := dbInfo.Open()
	check("opening database", err)

	// Read each statement from the file one at a time and execute them
	sqlChan := sqlStatements(fileName)
	for sql := range sqlChan {
		// Execute SQL.  If not successful, stop and ask user
		// whether or not we should continue
		fmt.Println("\n---")
		log.Print("Executing SQL: ", sql)

		runSql := true // Let's us loop for rerunning error
		for runSql {
			result, err := db.Exec(sql)

			// If there was an error, ask user whether or not we should continue
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: ", err)
				action := misc.ChooseOne("Continue, Quit, or Redo?  Enter c, q, or r: ", "c", "q", "r")
				// runSql has to be true here
				if action == "c" {
					runSql = false
				}
				if action == "q" {
					os.Exit(1)
				}
			} else {
				runSql = false
				rowCnt, err := result.RowsAffected()
				check("getting rows affected count", err)
				log.Printf("Rows affected: %d\n", rowCnt)
			}
		}

	}
	log.Println("Done!")
}

/*
 * sqlStatements chunks SQL lines from the given file into complete statements and returns
 * them a statement at a time (via a channel).  Most SQL statements end with a semi-colon.
 * However, some statements (like CREATE FUNCTION) are tricky because they can have statements
 * inside the main statement.  So this also looks for begin/end notations:
 * -- STATEMENT-BEGIN
 * -- STATEMENT-END
 */
func sqlStatements(fileName string) <-chan string {
	statementChan := make(chan string)

	go func() {
		lineChan, err := fileutil.ReadLinesChannel(fileName)
		check("reading file", err)

		// delimitedStatement is true when we are in a specially delimited statement:
		// -- STATEMENT-BEGIN
		// -- STATEMENT-END
		var delimitedStatement bool = false

		// TODO: Convert this to a string builder
		statement := ""
		for line := range lineChan {
			//fmt.Printf("  Line: %s\n", line)

			// Remove whitespace from the beginning and the end of the line
			line = strings.TrimSpace(line)

			// Ignore empty line
			if len(line) == 0 {
				continue
			}

			if stmtBeginRegex.MatchString(line) {
				delimitedStatement = true
				statement = "" // lose the current line
				continue
			}

			if stmtEndRegex.MatchString(line) {
				delimitedStatement = false
				if len(statement) > 0 {
					statementChan <- statement
					statement = "" // lose the current line
				}
				continue
			}

			// ignore lines that are fully commented out
			if strings.HasPrefix(line, "--") {
				continue
			}

			statement += line + "\n"

			// When we are not in a specially delimited statement, a line
			// ending with a semi-colon denotes the end of the statement.
			if !delimitedStatement && eosRegex.MatchString(line) {
				statementChan <- statement
				statement = ""
			}
		}

		close(statementChan)
	}()

	return statementChan
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s -f <sqlFileName> [-U <user>] [-h <host>] [-p <port>] [-d <dbname>] ... \n", os.Args[0])
	fmt.Fprintln(os.Stderr, `
Database connection properties can be specified in two ways:
  * Environment variables
  * Program flags (override environment variables)
  * ~/.pgpass file (for the password)

Environment variables are:
  PGHOST     : host name where database is running (default is localhost)
  PGPORT     : port database is listening on (default is 5432)
  PGDATABASE : name of database you want to update
  PGUSER     : user in postgres you'll be executing the commands as
  PGPASSWORD : password for the user
  PGOPTION   : postgresql options (like sslmode=disable)

Program flags are:
  -f, --filename     : required. file path to the SQL
  -V, --version      : prints the version of pgrun being run
  -?, --help         : prints a summary of the commands accepted by pgrun
  -U, --user         : user in postgres to execute the commands
  -h, --host         : host name where database is running (default is localhost)
  -p, --port         : port database is listening on (default is 5432)
  -d, --dbname       : database name
  -O, --options      : postgresql connection options (like sslmode=disable)
  -w, --no-password  : Never issue a password prompt
  -W, --password     : Force a password prompt
`)

	os.Exit(2)
}

func check(msg string, err error) {
	if err != nil {
		log.Fatal("Error "+msg, err)
	}
}
