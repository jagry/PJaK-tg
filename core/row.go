package core

import "database/sql"

func rowsClose(rows **sql.Rows) {
	if *rows != nil {
		(*rows).Close()
	}
}
