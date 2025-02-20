package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TableMeta struct {
	Columns    []string
	SelectStmt string
	OutputFile string
}

type TableProcessRequest struct {
	Db        *sqlx.DB
	Query     string
	TableName string
	Record    interface{}
	Writer    *CSVWriter
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
	log := getLogger()

	for tableName, meta := range clobColumns {
		record, found := GetStructForTable(tableName)
		if !found {
			return fmt.Errorf(
				"no struct mapping found for table: %s",
				tableName,
			)
		}

		log.Printf("Exporting data for table: %s\n", tableName)

		handler := NewCSVHandler()
		writer, err := handler.CreateWriter(meta.OutputFile, record)
		if err != nil {
			return fmt.Errorf("error creating CSV writer: %w", err)
		}

		err = orc.processTableRows(&TableProcessRequest{
			Db:        sqlxDB,
			Query:     meta.SelectStmt,
			TableName: tableName,
			Record:    record,
			Writer:    writer,
		})

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

func (orc *OracleApp) processTableRows(req *TableProcessRequest) error {
	rows, err := req.Db.Queryx(req.Query)
	if err != nil {
		return fmt.Errorf("query failed for %s: %w", req.TableName, err)
	}
	defer rows.Close()

	handler := NewCSVHandler()

	for rows.Next() {
		if err := rows.StructScan(req.Record); err != nil {
			return fmt.Errorf(
				"error scanning row in %s: %w",
				req.TableName,
				err,
			)
		}

		csvrow, err := handler.ToRow(req.Record)
		if err != nil {
			return fmt.Errorf("conversion error: %w", err)
		}

		if err := req.Writer.Write(csvrow); err != nil {
			return fmt.Errorf("CSV write error: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf(
			"error in scanning rows for table %s %w",
			req.TableName,
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
	lowerTable := strings.ToLower(table)

	var idColumn string
	switch lowerTable {
	case "paragraph":
		idColumn = "PARAGRAPH_NO"
	case "plasmid":
		idColumn = "ID"
	default:
		idColumn = fmt.Sprintf("%s_ID", table) // Explicit default assignment
	}

	return fmt.Sprintf(
		"SELECT %s,%s FROM %s WHERE %s",
		idColumn,
		strings.Join(columns, ","),
		table,
		strings.Join(conditions, " OR "),
	)
}
