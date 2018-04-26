#!/bin/sh
#sh createDB.sh "postgres" "postgres" "test"
user="$1"
password="$2"
dbName="$3"
export PGPASSWORD=$password

psql -U $user << END_OF_SCRIPT

DROP DATABASE $dbName; -- drop the DB

CREATE DATABASE $dbName WITH ENCODING 'UTF8' TEMPLATE template0;

END_OF_SCRIPT