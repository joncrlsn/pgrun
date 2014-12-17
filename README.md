# pgrun - a better way to run SQL against PostgreSQL

pgrun is (IMHO) a superior replacement of psql for running sql files against a PostgreSQL database.  It uses mostly the same argument flags as psql (-U, -h, -p, -d, -f) as well as using the standard postgreSQL environment variables like PGHOST, PGPORT, PGDATABASE, and PGPASSWORD.

Written in GoLang, pgrun executes each statement in a SQL file against a PostgreSQL database, stopping to ask you what you want to do when any statement has an error (you can Continue, Quit, or Redo the statement). Contrast this behavior with the standard psql command (which also takes a -f argument), but continues to run even after a statement fails.

Suggestions and modifications to make this more useful and "idiomatic Go" will be appreciated.

### download 
[osx64](https://github.com/joncrlsn/pgrun/raw/master/bin-osx64/pgrun "OSX 64-bit version")
[osx32](https://github.com/joncrlsn/pgrun/raw/master/bin-osx32/pgrun "OSX version")
[linux64](https://github.com/joncrlsn/pgrun/raw/master/bin-linux64/pgrun "Linux 64-bit version")
[linux32](https://github.com/joncrlsn/pgrun/raw/master/bin-linux32/pgrun "Linux version")
[win64](https://github.com/joncrlsn/pgrun/raw/master/bin-win64/pgrun.exe "Windows 64-bit version")
[win32](https://github.com/joncrlsn/pgrun/raw/master/bin-win32/pgrun.exe "Windows version")


### examples
	pgrun -U dbuser -h 10.10.41.55 -d userdb -f obfuscateUsers.sql
	PGUSER=dbuser PGHOST=10.10.41.55 PGDATABASE=userdb pgrun -f obfuscateUsers.sql

#### flags (these mostly match psql arguments):
program flag         | Explanation
-------------------: | -------------
  -f, --filename     | required. file path to the SQL
  -V, --version      | prints the version of pgrun being run
  -?, --help         | prints a summary of the commands accepted by pgrun
  -U, --user         | user in postgres to execute the commands
  -h, --host         | host name where database is running (default is localhost)
  -p, --port         | port database is listening on (default is 5432)
  -d, --dbname       | database name
  -O, --options      | postgresql connection options (like sslmode=disable)
  -w, --no-password  | Never issue a password prompt
  -W, --password     | Force a password prompt

### database connection options

  * Use environment variables (see table below)
  * Program flags (overrides environment variables)
  * ~/.pgpass file
  * Note that if password is not specified, you will be prompted.

### optional database environment variables

Name       | Explanation
---------  | -----------
PGHOST     | host name where database is running (matches psql)
PGPORT     | port database is listening on (default is 5432) (matches psql)
PGDATABASE | name of database you want to copy (matches psql)
PGUSER     | user in postgres you'll be executing the queries as (matches psql)
PGPASSWORD | password for the user (matches psql)
PGOPTION   | one or more database options (like sslmode=disable)
