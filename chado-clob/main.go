package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	go_ora "github.com/sijms/go-ora/v2"
	"github.com/urfave/cli/v2"
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
ORDER BY 
    c.table_name, 
    c.column_name`

type TableMeta struct {
	Columns    []string
	SelectStmt string
	OutputFile string
}

func main() {
	app := &cli.App{
		Name:  "oracle-ping",
		Usage: "Analyze Oracle database schema",
		Commands: []*cli.Command{
			{
				Name:   "ping",
				Usage:  "Test database connectivity",
				Action: pingAction,
			},
			{
				Name:   "clob-stats",
				Usage:  "List tables with populated CLOB columns",
				Action: clobStatsAction,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Aliases:  []string{"H"},
				Usage:    "Database host",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"P"},
				Usage:   "Database port",
				Value:   1521,
			},
			&cli.StringFlag{
				Name:     "service",
				Aliases:  []string{"s"},
				Usage:    "Database service name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "Database user",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Usage:    "Database password",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output-folder",
				Aliases: []string{"o"},
				Usage:   "Output directory for CSV files",
				Value:   ".",
			},
		},
		Action: clobStatsAction,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func clobStatsAction(cltx *cli.Context) error {
	dbh, err := setupDatabaseConnection(cltx)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to connect: %v", err), 1)
	}
	defer dbh.Close()

	rows, err := queryClobTables(dbh, cltx.String("user"))
	if err != nil {
		return cli.Exit(fmt.Sprintf("Query failed: %v", err), 1)
	}
	defer rows.Close()

	clobColumns, err := processClobRows(rows, cltx.String("output-folder"))
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error processing rows: %v", err), 1)
	}

	return nil
}

func setupDatabaseConnection(cltx *cli.Context) (*sql.DB, error) {
	connStr := go_ora.BuildUrl(
		cltx.String("host"),
		cltx.Int("port"),
		cltx.String("service"),
		cltx.String("user"),
		cltx.String("password"),
		nil,
	)

	dbh, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return dbh, nil
}

func queryClobTables(dbh *sql.DB, user string) (*sql.Rows, error) {
	rows, err := dbh.Query(clobQuery, user)
	if err != nil {
		return nil, fmt.Errorf("error in running the clob query %s", err)
	}

	return rows, nil
}

func Map[T any, U any](slice []T, f func(T) U) []U {
	result := make([]U, 0)
	for _, v := range slice {
		result = append(result, f(v))
	}

	return result
}

func processClobRows(
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

	// Assign generated select statements to table metadata
	for table, meta := range clobColumns {
		stmt := generateSelectStatement(
			table,
			meta.Columns,
		)
		clobColumns[table].SelectStmt = stmt
	}

	return clobColumns, nil
}

func generateSelectStatement(table string, columns []string) string {
	conditions := Map(columns, func(col string) string {
		return fmt.Sprintf("%s IS NOT NULL", col)
	})
	
	return fmt.Sprintf(
		"SELECT %s_ID,%s FROM %s WHERE %s",
		table,
		strings.Join(columns, ","),
		table,
		strings.Join(conditions, " OR "),
	)
}

func pingAction(cltx *cli.Context) error {
	// Build connection string
	connStr := go_ora.BuildUrl(
		cltx.String("host"),
		cltx.Int("port"),
		cltx.String("service"),
		cltx.String("user"),
		cltx.String("password"),
		nil,
	)

	// Open database connection
	dbh, err := sql.Open("oracle", connStr)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to connect: %v", err), 1)
	}
	defer dbh.Close()

	// Test connection
	err = dbh.Ping()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Connection failed: %v", err), 1)
	}

	fmt.Println("Successfully connected to Oracle database!")

	return nil
}
