package app

import "database/sql"

// DB holds the shared database connection, if configured.
var DB *sql.DB
