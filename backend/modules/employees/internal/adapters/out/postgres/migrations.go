package postgres

import "embed"

// Migrations are embedded in the binary, so the version of the schema always
// matches the code that runs against it.
//
//go:embed migrations/*.sql
var Migrations embed.FS
