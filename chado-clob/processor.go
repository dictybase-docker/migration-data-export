package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
)

const clobQuery = `SELECT 
    c.table_name,
    c.column_name
FROM 
    all_tab_columns c
JOIN 
    all_tables t 
    ON c.owner = t.owner 
    AND c.table_name = t.table_name
WHERE 
    c.owner = :1 
    AND c.data_type = 'CLOB'
    AND t.num_rows > 0
    AND c.table_name NOT IN ('FEATURE', 'FEATUREPROP', 'CHADO_LOGS', 'CHADOPROP')
    AND NOT EXISTS (
        SELECT 1 
        FROM all_mviews mv 
        WHERE mv.owner = c.owner 
        AND mv.mview_name = c.table_name
    )
ORDER BY 
    c.table_name, 
    c.column_name`

type TableMeta struct {
	Columns    []string
	SelectStmt string
	OutputFile string
}

func (orc *OracleApp) processClobRows(
	rows *sql.Rows,
	outputFolder string,
) (map[string]*TableMeta, error) {
	clobColumns := make(map[string]*TableMeta)

	for rows.Next() {
		var table, column string
		if err := rows.Scan(&table, &column); err != nil {
			return nil, fmt.Errorf(
				"error scanning table %s: %w",
				table,
				err,
			)
		}

		if _, exists := clobColumns[table]; !exists {
			clobColumns[table] = &TableMeta{
				OutputFile: filepath.Join(
					outputFolder,
					fmt.Sprintf("%s_clob_data.csv", strings.ToLower(table)),
				),
			}
		}
		clobColumns[table].Columns = append(clobColumns[table].Columns, column)
	}
	if err := rows.Err(); err != nil {
		return clobColumns, fmt.Errorf("error in scanning rows %s", err)
	}

	for table, meta := range clobColumns {
		stmt := generateSelectStatement(
			table,
			meta.Columns,
		)
		clobColumns[table].SelectStmt = stmt
	}

	return clobColumns, nil
}

func (orc *OracleApp) processClobData(
	dbh *sql.DB,
	clobColumns map[string]*TableMeta,
) error {
	// Create single sqlx instance for all tables
	sqlxDB := sqlx.NewDb(dbh, "oracle")

	for tableName, meta := range clobColumns {
		record, found := GetStructForTable(tableName)
		if !found {
			return fmt.Errorf(
				"no struct mapping found for table: %s",
				tableName,
			)
		}

		writer, err := createCSVWriter(meta.OutputFile, record)
		if err != nil {
			return fmt.Errorf("error creating CSV writer: %w", err)
		}

		// Process rows
		err = orc.processTableRows(
			sqlxDB,
			meta.SelectStmt,
			tableName,
			record,
			writer,
		)

		// Close immediately after processing instead of deferring
		if closeErr := writer.Close(); closeErr != nil {
			if err == nil { // Preserve original error if any
				err = fmt.Errorf("error closing writer: %w", closeErr)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (orc *OracleApp) processTableRows(
	dbh *sqlx.DB,
	query string,
	tableName string,
	record interface{},
	writer *CSVWriter,
) error {
	rows, err := dbh.Queryx(query)
	if err != nil {
		return fmt.Errorf("query failed for %s: %w", tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		// Reset the struct for each row
		if err := rows.StructScan(record); err != nil {
			return fmt.Errorf("error scanning row in %s: %w", tableName, err)
		}

		csvrow, err := structToCSVRow(record)
		if err != nil {
			return fmt.Errorf("conversion error: %w", err)
		}

		if err := writer.Write(csvrow); err != nil {
			return fmt.Errorf("CSV write error: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf(
			"error in scanning rows for table %s %w",
			tableName,
			err,
		)
	}

	return nil
}

func Map[T any, U any](slice []T, f func(T) U) []U {
	result := make([]U, 0)
	for _, v := range slice {
		result = append(result, f(v))
	}
	return result
}

func buildNotNullCondition(col string) string {
	return fmt.Sprintf("%s IS NOT NULL", col)
}

func generateSelectStatement(table string, columns []string) string {
	conditions := Map(columns, buildNotNullCondition)
	return fmt.Sprintf(
		"SELECT %s_ID,%s FROM %s WHERE %s",
		table,
		strings.Join(columns, ","),
		table,
		strings.Join(conditions, " OR "),
	)
}

func (orc *OracleApp) queryClobTables(
	dbh *sql.DB,
	user string,
) (*sql.Rows, error) {
	rows, err := dbh.Query(clobQuery, user)
	if err != nil {
		return nil, fmt.Errorf("error in running the clob query %s", err)
	}

	return rows, nil
}
