package main

import "fmt"
import "os"
import "strings"
import "regexp"
import flag "github.com/ogier/pflag"
import "bufio"
import "log"
import "github.com/joncrlsn/pgutil"
import "github.com/joncrlsn/fileutil"
import "github.com/joncrlsn/misc"

// End of Statement regex
var eosRegex = regexp.MustCompile(`;\s*$|;\s*--.*$`)
var inReader = bufio.NewReader(os.Stdin)
var version = "1.0.7"

// Executes a file of SQL statements one statement at a time, stopping everything
// if one of them has an error
func main() {

	// -f (filename) is a required program argument
	var fileName string
	var verFlag bool
	var helpFlag bool
	flag.StringVar(&fileName, "f", "", "path of the SQL file to run")
	flag.BoolVarP(&verFlag, "version", "V", false, "Displays version information")
	flag.BoolVarP(&helpFlag, "help", "?", false, "Displays usage help")
	dbInfo := pgutil.DbInfo{}
	dbInfo.Populate()

	if verFlag {
		fmt.Fprintf(os.Stderr, "%s - version %s\n", os.Args[0], version)
		os.Exit(0)
	}

	if helpFlag {
		usage()
	}

	if len(fileName) == 0 {
		fmt.Fprintln(os.Stderr, "Missing required filename argument (-f)")
		usage()
	}

	exists, _ := fileutil.Exists(fileName)
	if !exists {
		fmt.Fprintf(os.Stderr, "File does not exist: %s\n", fileName)
		os.Exit(2)
	}

	runFile(fileName, &dbInfo)
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
 * Reads and returns (via channel) SQL statements from the given file.
 * SQL statements must end with a semi-colon
 */
func sqlStatements(fileName string) <-chan string {
	statementChan := make(chan string)

	go func() {
		lineChan, err := fileutil.ReadLinesChannel(fileName)
		check("reading file", err)

		// TODO: Convert this to a string builder
		statement := ""
		for line := range lineChan {
			//fmt.Printf("  Line: %s\n", line)

			// ignore blank or empty lines
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			if strings.HasPrefix(line, "--") {
				continue
			}

			statement += line + "\n"

			// look for line ending with just a semi-colon
			// or a semi-colon with a SQL comment following
			if eosRegex.MatchString(line) {
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
