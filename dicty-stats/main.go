package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"

	go_ora "github.com/sijms/go-ora/v2"
	"github.com/urfave/cli/v2"
)

const (
	strainQuery = `SELECT DBXREF_ID FROM CGM_DDB.STOCK_CENTER`
	dbxrefQuery = `SELECT ACCESSION FROM DBXREF WHERE DBXREF_ID = :dbxref_id`
	phenoQuery  = `
    SELECT phen.name,p.timecreated
	FROM phenstatement pst

	LEFT JOIN genotype g on g.genotype_id = pst.genotype_id
	
	LEFT JOIN cvterm env on env.cvterm_id = pst.environment_id
	LEFT JOIN cv env_cv on env_cv.cv_id = env.cv_id
	
	LEFT JOIN phenotype p on p.phenotype_id = pst.phenotype_id
	LEFT JOIN cvterm phen on phen.cvterm_id = p.observable_id
	
	LEFT JOIN cvterm assay on assay.cvterm_id = p.assay_id
	LEFT JOIN cv assay_cv on assay_cv.cv_id = assay.cv_id
	
	LEFT JOIN pub on pub.pub_id = pst.pub_id
    WHERE g.uniquename = :accession
    `
	geneDescQuery = `
        SELECT genef.uniquename,dbxref.accession,
               fp.value description,fp.timecreated,
               fp.created_by FROM V_GENE_FEATURES genef
        JOIN featureprop fp 
            ON genef.feature_id = fp.feature_id
        JOIN cvterm 
            ON fp.type_id = cvterm.cvterm_id
        JOIN dbxref 
            ON dbxref.dbxref_id = genef.dbxref_id
        where cvterm.name = 'description'
        ORDER BY fp.timecreated
    `
)

func main() {
	app := &cli.App{
		Name:  "data-export",
		Usage: "A tool for exporting data",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "Oracle `USER` name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "pass",
				Aliases:  []string{"p"},
				Usage:    "Oracle `PASSWORD`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "host",
				Aliases:  []string{"H"},
				Usage:    "Oracle `HOST`",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"P"},
				Usage:   "Oracle `PORT`",
				Value:   1521,
			},
			&cli.StringFlag{
				Name:    "sid",
				Aliases: []string{"s"},
				Usage:   "Oracle SID",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "strain-pheno",
				Usage: "extract strain and phenotype information",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "output file name",
						Value:   "output.csv",
					},
				},
				Action: strainPhenoAction,
			},
			{
				Name:  "gene-desc",
				Usage: "extract gene description information",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "output file name",
						Value:   "output.csv",
					},
				},
				Action: geneDescAction,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func geneDescAction(cltx *cli.Context) error {
	header := []string{
		"name",
		"accession",
		"createdOn",
		"description",
		"createdOn",
		"createdBy",
	}
	writer, err := os.Create(cltx.String("output"))
	if err != nil {
		return cli.Exit(
			fmt.Sprintf(
				"error in opening file for writing %s",
				err,
			),
			2,
		)
	}
	defer writer.Close()

	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	conn, err := setupDatabaseConnection(cltx)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	defer conn.Close()

	if err := writeCSVHeader(csvWriter, header); err != nil {
		return cli.Exit(
			fmt.Sprintf("error in writing csv header %s", err),
			2,
		)
	}

	if err := processGeneSummary(conn, csvWriter); err != nil {
		return cli.Exit(err.Error(), 2)
	}

	if err := csvWriter.Error(); err != nil {
		return cli.Exit(
			fmt.Sprintf(
				"error in finish writing with csv writer %s",
				err,
			),
			2,
		)
	}

	return nil
}

func strainPhenoAction(cltx *cli.Context) error {
	header := []string{"strain_ID", "phenotype", "createdOn"}
	writer, err := os.Create(cltx.String("output"))
	if err != nil {
		return cli.Exit(
			fmt.Sprintf(
				"error in opening file for writing %s",
				err,
			),
			2,
		)
	}
	defer writer.Close()

	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	conn, err := setupDatabaseConnection(cltx)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	defer conn.Close()

	if err := writeCSVHeader(csvWriter, header); err != nil {
		return cli.Exit(
			fmt.Sprintf("error in writing csv header %s", err),
			2,
		)
	}

	if err := processStrains(conn, csvWriter); err != nil {
		return cli.Exit(err.Error(), 2)
	}

	if err := csvWriter.Error(); err != nil {
		return cli.Exit(
			fmt.Sprintf(
				"error in finish writing with csv writer %s",
				err,
			),
			2,
		)
	}

	return nil
}

func setupDatabaseConnection(cltx *cli.Context) (*sql.DB, error) {
	connStr := go_ora.BuildUrl(
		cltx.String("host"),
		cltx.Int("port"),
		cltx.String("sid"),
		cltx.String("user"),
		cltx.String("pass"),
		nil,
	)
	return sql.Open("oracle", connStr)
}

func writeCSVHeader(csvWriter *csv.Writer, header []string) error {
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("error in writing csv header %s", err)
	}
	return nil
}

func processGeneSummary(conn *sql.DB, csvWriter *csv.Writer) error {
	geneDescRows, err := conn.Query(geneDescQuery)
	if err != nil {
		return fmt.Errorf("error in running gene description query %s", err)
	}
	defer geneDescRows.Close()

	for geneDescRows.Next() {
		var name, accession, description, createdOn, createdBy string
		if err = geneDescRows.Scan(&name, &accession, &description, &createdOn, &createdBy); err != nil {
			return fmt.Errorf("error in scanning row %s", err)
		}
		row := []string{name, accession, description, createdOn, createdBy}
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("error in writing gene description row %s", err)
		}
	}

	return nil
}

func processStrains(conn *sql.DB, csvWriter *csv.Writer) error {
	dbxrefStmt, err := conn.Prepare(dbxrefQuery)
	if err != nil {
		return fmt.Errorf("error in preparing dbxref query %s", err)
	}
	defer dbxrefStmt.Close()

	phenStmt, err := conn.Prepare(phenoQuery)
	if err != nil {
		return fmt.Errorf("error in preparing phenotype query %s", err)
	}
	defer phenStmt.Close()

	geneDescRows, err := conn.Query(strainQuery)
	if err != nil {
		return fmt.Errorf("error in row query %s", err)
	}
	defer geneDescRows.Close()

	for geneDescRows.Next() {
		var dbxrefId string
		if err = geneDescRows.Scan(&dbxrefId); err != nil {
			return fmt.Errorf("error in scanning row %s", err)
		}

		if err := processAccession(dbxrefStmt, phenStmt, dbxrefId, csvWriter); err != nil {
			return err
		}
	}

	return nil
}

func processAccession(
	dbxrefStmt, phenStmt *sql.Stmt,
	dbxrefId string,
	csvWriter *csv.Writer,
) error {
	var accession, createdOn string
	if err := dbxrefStmt.QueryRow(dbxrefId).Scan(&accession); err != nil {
		return fmt.Errorf("error in running accession query %s", err)
	}

	phenoRows, err := phenStmt.Query(accession)
	if err != nil {
		return fmt.Errorf(
			"error in running phenotypes query for accession %s %s",
			accession,
			err,
		)
	}
	defer phenoRows.Close()

	for phenoRows.Next() {
		var phenotype string
		if err := phenoRows.Scan(&phenotype, &createdOn); err != nil {
			return fmt.Errorf("error in scanning phenotype row %s", err)
		}

		if err := csvWriter.Write([]string{accession, phenotype, createdOn}); err != nil {
			return fmt.Errorf("error in writing phenotype row in csv %s", err)
		}
	}

	return nil
}
