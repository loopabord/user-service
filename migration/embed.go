package migration

import "embed"

//go:embed *.sql
var SqlMigrations embed.FS
