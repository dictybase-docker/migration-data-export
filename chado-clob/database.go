package main

import (
	"database/sql"
	"fmt"
	
	go_ora "github.com/sijms/go-ora/v2"
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

func (orc *OracleApp) setupDatabaseConnection() (*sql.DB, error) {
	connStr := go_ora.BuildUrl(
		orc.cltx.String("host"),
		orc.cltx.Int("port"),
		orc.cltx.String("service"),
		orc.cltx.String("user"),
		orc.cltx.String("password"),
		nil,
	)

	dbh, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	return dbh, nil
}

func (orc *OracleApp) queryClobTables(dbh *sql.DB, user string) (*sql.Rows, error) {
	rows, err := dbh.Query(clobQuery, user)
	if err != nil {
		return nil, fmt.Errorf("error in running the clob query %s", err)
	}
	return rows, nil
}
