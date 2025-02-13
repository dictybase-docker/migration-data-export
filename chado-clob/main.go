package main

import (
	"database/sql"
	"fmt"
	"os"

	go_ora "github.com/sijms/go-ora/v2"
	"github.com/urfave/cli/v2"
)

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
	if err := processClobData(dbh, clobColumns); err != nil {
		return cli.Exit(err.Error(), 2)
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
