package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	go_ora "github.com/sijms/go-ora/v2"
	"github.com/urfave/cli/v2"
)

type OracleApp struct {
	cltx *cli.Context
}

func main() {
	app := &cli.App{
		Name:  "oracle-ping",
		Usage: "Analyze Oracle database schema",
		Commands: []*cli.Command{
			{
				Name:  "ping",
				Usage: "Test database connectivity",
				Action: func(cltx *cli.Context) error {
					orc := &OracleApp{cltx: cltx}
					return orc.pingAction()
				},
			},
			{
				Name:  "clob-stats",
				Usage: "List tables with populated CLOB columns",
				Action: func(cltx *cli.Context) error {
					orc := &OracleApp{cltx: cltx}
					return orc.clobStatsAction()
				},
			},
			{
				Name:  "list-tables",
				Usage: "Export all user-owned table names to a file",
				Action: func(cltx *cli.Context) error {
					orc := &OracleApp{cltx: cltx}
					return orc.listTablesAction()
				},
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
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"f"},
				Usage:   "Output file path for table list",
				Value:   "tables.txt",
			},
		},
		Action: func(c *cli.Context) error {
			orc := &OracleApp{cltx: c}
			return orc.clobStatsAction()
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func (orc *OracleApp) clobStatsAction() error {
	dbh, err := orc.setupDatabaseConnection()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to connect: %v", err), 1)
	}
	defer dbh.Close()

	rows, err := orc.queryClobTables(dbh, orc.cltx.String("user"))
	if err != nil {
		return cli.Exit(fmt.Sprintf("Query failed: %v", err), 1)
	}
	defer rows.Close()

	clobColumns, err := orc.processClobRows(
		rows,
		orc.cltx.String("output-folder"),
	)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error processing rows: %v", err), 1)
	}
	for tableName, meta := range clobColumns {
		fmt.Printf("table: %s | statement: %s\n", tableName, meta.SelectStmt)
	}
	if err := orc.processClobData(dbh, clobColumns); err != nil {
		return cli.Exit(err.Error(), 2)
	}

	return nil
}

func getLogger() *log.Logger {
	return log.New(
		os.Stderr,
		"[chado-clob] ",
		log.LstdFlags|log.Lmsgprefix,
	)
}

func (orc *OracleApp) listTablesAction() error {
	dbh, err := orc.setupDatabaseConnection()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to connect: %v", err), 1)
	}
	defer dbh.Close()

	rows, err := dbh.Query(tableListQuery, orc.cltx.String("user"))
	if err != nil {
		return cli.Exit(fmt.Sprintf("Query failed: %v", err), 1)
	}
	defer rows.Close()

	outputFile := orc.cltx.String("output")
	f, err := os.Create(outputFile)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to create output file: %v", err), 1)
	}
	defer f.Close()

	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return cli.Exit(fmt.Sprintf("Error scanning row: %v", err), 1)
		}
		if _, err := fmt.Fprintln(f, tableName); err != nil {
			return cli.Exit(fmt.Sprintf("Error writing to file: %v", err), 1)
		}
	}
	
	fmt.Printf("Successfully exported table names to %s\n", outputFile)
	return nil
}

func (orc *OracleApp) pingAction() error {
	cltx := orc.cltx
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
