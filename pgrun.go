package main

import "fmt"
import "os"
import "strings"
import "regexp"
import "flag"
import "bufio"
import "log"
import "github.com/joncrlsn/pgutil"
import "github.com/joncrlsn/fileutil"
import "github.com/joncrlsn/misc"

// End of Statement regex
var eosRegex *regexp.Regexp = regexp.MustCompile(`;\s*$|;\s*--.*$`)
var inReader *bufio.Reader = bufio.NewReader(os.Stdin)

// Executes a file of SQL statements one statement at a time, stopping everything
// if one of them has an error
func main() {

    // -f (filename) is a required program argument
	var fileName string
	flag.StringVar(&fileName, "f", "", "path of the SQL file to run")
	dbInfo := pgutil.DbInfo{}
	dbInfo.DbOptions = "sslmode=disable"
	dbInfo.Populate()

	if len(fileName) == 0 {
        fmt.Println("Missing required filename argument (-f)")
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
    check("Opening database", err)

    // Read each statement from the file one at a time and execute them
	sqlChan := sqlStatements(fileName)
	for sql := range sqlChan {
		// Execute SQL.  If not successful, stop and ask user
		// whether or not we should continue
		fmt.Println("== Executing SQL: ", sql)
        _, err := db.Exec(sql)  

        // If there was an error, ask user whether or not to continue
		if err != nil {
            log.Println("Error:", err) 
            if misc.PromptYesNo("== SQL failed.  Do you want to continue?", false)  {
                continue
            }
            os.Exit(1)
		}
	}
}

/*
 * Reads and returns (via channel) SQL statements from the given file.
 * SQL statements must end with a semi-colon
 */
func sqlStatements(fileName string) <-chan string {
	statementChan := make(chan string)

	go func() {
		lineChan := fileutil.ReadLinesChannel(fileName)

		// TODO: Convert this to a string builder
		statement := ""
		for line := range lineChan {
			//fmt.Printf("  Line: %s\n", line)

			// ignore blank or empty lines
			if len(strings.TrimSpace(line)) == 0 {
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
	fmt.Fprintf(os.Stderr, "usage: %s -f <sqlFileName> [-host <string>] [-port <int>] [-db <string>] [-user <string>] [-password <string>] \n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func check(msg string, err error) {
    if err != nil {
        log.Fatal(msg, err)
    }
}
