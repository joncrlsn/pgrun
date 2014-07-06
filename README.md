pgrun
=====

Executes each statement in a SQL file against a PostgreSQL database, stopping when any statement has an error.Contrast this behavior with the standard PostgreSQL psql executable which takes a -f argument, but continues to run even after a command fails.

Database connection properties can be specified in multiple ways:
  * Environment variables
  * Program flags

Environment variables are:
  DBHOST : host name where database is running (default is localhost)
  DBPORT : port database is listening on (default is 5432)
  DBNAME : database you want to update
  DBUSER : user in postgres you'll be executing the commands as
  DBPASS : password for the user

Program flags:
  -f    : file path to the SQL to run
  -host : host name where database is running--default is localhost
  -port : defaults to 5432
  -db   : database name
  -user : user in postgres to execute the commands
  -pw   : password for the user
