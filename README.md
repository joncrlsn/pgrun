pgrun
=====

Written in GoLang. Executes each statement in a SQL file against a PostgreSQL database, stopping when any statement has an error.Contrast this behavior with the standard PostgreSQL psql executable which takes a -f argument, but continues to run even after a command fails.

Database connection properties can be specified in two ways:
  * Environment variables
  * Program flags

If you have a ~/.pgpass file, pgrun will attempt to pull password from there.

Environment variables are:
  * DBHOST : host name where database is running (default is localhost)
  * DBPORT : port database is listening on (default is 5432)
  * DBNAME : database you want to update
  * DBUSER : user in postgres you'll be executing the commands as
  * DBPASS : password for the user

Program flags (these match the psql arguments):
  * -f  : (required) file path to the SQL to run
  * -U  : user in postgres to execute the commands
  * -h  : host name where database is running -- default is localhost
  * -p  : defaults to 5432
  * -d  : database name
  * -pw : password (does not match psql flag, and not required. You will be prompted if necessary)
