# pgrun - a customized psql clone

pgrun is a customized replacement of psql for running sql files against a PostgreSQL database when you want to be notified that an update or alter failed before continuing on.  It uses mostly the same argument flags as psql (-U, -h, -p, -d, -f) as well as the standard postgreSQL environment variables like PGHOST, PGPORT, PGDATABASE, and PGPASSWORD.

Written in GoLang, pgrun executes each statement from the given SQL file against a PostgreSQL database, stopping on error to ask you what you want to do (you can Continue, Quit, or Redo the statement). Contrast this behavior with the standard psql command which is all or nothing (continue if any fail, or quit if any fail)

### download 
[osx](https://github.com/joncrlsn/pgrun/raw/master/bin-osx/pgrun "OSX version")
[linux](https://github.com/joncrlsn/pgrun/raw/master/bin-linux/pgrun "Linux version")
[win](https://github.com/joncrlsn/pgrun/raw/master/bin-win/pgrun.exe "Windows version")
[arm64](https://github.com/joncrlsn/pgrun/raw/master/bin-arm64/pgrun.exe "Mac arm64 version")


### examples
	pgrun -U dbuser -h 10.10.41.55 -d userdb -f obfuscateUsers.sql
	PGUSER=dbuser PGHOST=10.10.41.55 PGDATABASE=userdb pgrun -f obfuscateUsers.sql

#### flags/options (these mostly match psql arguments):
program flag/option  | explanation
-------------------: | -------------
  -f, --filename     | required. file path to the SQL
  -V, --version      | prints the version of pgrun being run
  -?, --help         | prints a summary of the commands accepted by pgrun
  -U, --user         | user in postgres to execute the commands
  -h, --host         | host name where database is running (default is localhost)
  -p, --port         | port database is listening on (default is 5432)
  -d, --dbname       | database name
  -O, --options      | postgresql connection options (like sslmode=disable)
  -w, --no-password  | Never issue a db password prompt.  Fail if none found.
  -W, --password     | Force a db password prompt

### database connection options

  * Use environment variables (see table below)
  * Program flags (overrides environment variables)
  * ~/.pgpass file
  * Note that if password is not specified, you will be prompted.

### optional database environment variables

name       | explanation
---------  | -----------
PGHOST     | host name where database is running (matches psql)
PGPORT     | port database is listening on (default is 5432) (matches psql)
PGDATABASE | name of database you want to copy (matches psql)
PGUSER     | user in postgres you'll be executing the queries as (matches psql)
PGPASSWORD | password for the user (matches psql)
PGOPTION   | one or more database options (like sslmode=disable)

### todo
1. ~~Fix bug where Ctrl-C in the password entry field messes up the console.~~ Fixed in version 1.0.7
1. ~~Fix -? and -V flags that are not working.~~ Fixed in version 1.0.8
1. Allow editing of a failed SQL statement before rerunning.
1. Improve the accuracy of parsing ~/.pgpass
