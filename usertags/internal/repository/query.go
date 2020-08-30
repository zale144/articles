package repository

import (
	"fmt"
	"strings"
)

func buildInsertStr(rowsL int, tableName string, cols []string) string {

	valueStrings := make([]string, 0, rowsL)
	valueArgs := make([]interface{}, 0, rowsL*len(cols))

	for i := 0; i < rowsL; i++ {

		var placeholders []string
		for j, col := range cols {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i*len(cols)+j+1))
			valueArgs = append(valueArgs, col)
		}

		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(placeholders, ",")))
	}

	stmt := fmt.Sprintf("INSERT INTO %s (%s) ", tableName, strings.Join(cols, ","))
	if rowsL > 0 {
		stmt += fmt.Sprintf("VALUES %s", strings.Join(valueStrings, ","))
	}

	return stmt
}
