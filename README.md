pgrun
=====

pgrun, intended to replace psql for running sql files, uses mostly the same argument flags as psql (-U, -h, -p, -d, -f) as well as using the standard postgres environment variables like PGHOST, PGPORT, PGDATABASE, and PGPASSWORD.

Written in GoLang, pgrun executes each statement in a SQL file against a PostgreSQL database, stopping to ask you what you want to do when any statement has an error (you can Continue, Quit, or Redo the statement). Contrast this behavior with the standard psql command (which also takes a -f argument), but continues to run even after a statement fails.  

A couple of binaries to save you the effort:
[Mac](https://github.com/joncrlsn/pgrun/raw/master/bin-osx/pgrun "OSX version")  [Linux](https://github.com/joncrlsn/pgrun/raw/master/bin-linux/pgrun "Linux version")

Database connection properties can be specified in two ways:
  * Environment variables
  * Program flags

If you have a ~/.pgpass file, pgrun will attempt to pull password from there.

#### Database connection information can be specified in up to three ways:

  * Environment variables (keeps you from typing them in often)
  * Program flags (overrides environment variables.  See above)
  * ~/.pgpass file (may contain password for the previously specified user)
  * Note that if password is not specified, you will be prompted.

#### Optional database environment variables (these match the postgres standard)

Name       | Explanation
---------  | -----------
PGHOST     | host name where database is running (matches psql)
PGPORT     | port database is listening on (default is 5432) (matches psql)
PGDATABASE | name of database you want to copy (matches psql)
PGUSER     | user in postgres you'll be executing the queries as (matches psql)
PGPASSWORD | password for the user (matches psql)
PGOPTION   | one or more database options (like sslmode=disable)

#### Program flags (these match the psql arguments):
  * -f  : (required) file path to the SQL to run
  * -U  : user in postgres to execute the commands
  * -h  : host name where database is running -- default is localhost
  * -p  : defaults to 5432
  * -d  : database name
  * -pw : password (does not match psql flag, and not required. You will be prompted if necessary)
