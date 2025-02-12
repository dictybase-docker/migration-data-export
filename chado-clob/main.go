package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    go_ora "github.com/sijms/go-ora/v2"
    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "oracle-ping",
        Usage: "Test connectivity to Oracle databases",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:     "host",
                Aliases:  []string{"H"},
                Usage:    "Database host",
                Required: true,
            },
            &cli.IntFlag{
                Name:     "port",
                Aliases:  []string{"P"},
                Usage:    "Database port",
                Value:    1521,
                Required: true,
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
        },
        Action: func(c *cli.Context) error {
            // Build connection string
            connStr := go_ora.BuildUrl(
                c.String("host"),
                c.Int("port"),
                c.String("service"),
                c.String("user"),
                c.String("password"),
                nil,
            )

            // Open database connection
            db, err := sql.Open("oracle", connStr)
            if err != nil {
                return cli.Exit(fmt.Sprintf("Failed to connect: %v", err), 1)
            }
            defer db.Close()

            // Test connection
            err = db.Ping()
            if err != nil {
                return cli.Exit(fmt.Sprintf("Connection failed: %v", err), 1)
            }

            fmt.Println("Successfully connected to Oracle database!")
            return nil
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
