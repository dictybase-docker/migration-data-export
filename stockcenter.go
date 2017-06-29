package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	_ "gopkg.in/rana/ora.v4"
	"gopkg.in/urfave/cli.v1"
)

var outfolder string

const layout = "2006-01-02 15:04:05"

func validateDscUsers(c *cli.Context) error {
	for _, p := range []string{"host", "sid", "user", "password"} {
		if !c.IsSet(p) {
			return cli.NewExitError(
				fmt.Sprintf("flag %s is required", p),
				2,
			)
		}
	}
	return nil
}

func validateDsc(c *cli.Context) error {
	if !ValidateArgs(c) {
		return cli.NewExitError("one or more of required arguments are not provided", 2)
	}
	if !ValidateExtraArgs(c) {
		return cli.NewExitError("one or more of required arguments are not provided", 2)
	}
	return nil
}

func validateDscOrder(c *cli.Context) error {
	if !ValidateExtraArgs(c) {
		return cli.NewExitError("one or more of required arguments are not provided", 2)
	}
	return nil
}

func getOracleConnection(c *cli.Context) (*sql.DB, error) {
	dataSource := fmt.Sprintf(
		"%s/%s@%s:%s/%s",
		c.String("user"),
		c.String("password"),
		c.String("host"),
		c.String("port"),
		c.String("sid"),
	)
	return sql.Open("ora", dataSource)
}

func DscUsersAction(c *cli.Context) error {
	CreateRequiredFolder(outfolder)
	if err := exportPlasmidUsers(c); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	if err := exportStrainUsers(c); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	return nil
}

func exportPlasmidUsers(c *cli.Context) error {
	log := getLogger(c)
	outfile := filepath.Join(outfolder, "plasmid_user_annotations.csv")
	writer, err := os.Create(outfile)
	if err != nil {
		return fmt.Errorf("unable to open file %s", err)
	}
	defer writer.Close()
	csv := csv.NewWriter(writer)

	dbh, err := getOracleConnection(c)
	if err != nil {
		return fmt.Errorf("error in connecting to database %s", err)
	}
	defer dbh.Close()
	log.Infof("started writing annotation in %s", outfile)
	var (
		plasmidId  int
		createdBy  string
		createdOn  time.Time
		modifiedOn time.Time
	)
	stmt := `
	SELECT plasmid.id,plasmid.created_by,plasmid.date_created, plasmid.date_modified
	FROM CGM_DDB.PLASMID plasmid
	`
	rows, err := dbh.Query(stmt)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("unable to run query %s", err), 2)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&plasmidId, &createdBy, &createdOn, &modifiedOn)
		if err != nil {
			return fmt.Errorf("unable to scan the next row %s", err)
		}
		err = csv.Write(
			[]string{
				fmt.Sprintf("DBP%07d", plasmidId),
				createdBy,
				createdOn.Format(layout),
				modifiedOn.Format(layout),
			},
		)
		if err != nil {
			return fmt.Errorf("unable to write csv row %s", err)
		}
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("unable to close the rows %s", err)
	}
	csv.Flush()
	if err := csv.Error(); err != nil {
		return fmt.Errorf("unable to finish csv writing %s", err)
	}
	log.Infof("finished writing annotations to %s", outfile)
	return nil
}

func exportStrainUsers(c *cli.Context) error {
	log := getLogger(c)
	outfile := filepath.Join(outfolder, "strain_user_annotations.csv")
	writer, err := os.Create(outfile)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("unable to open file %s", err), 2)
	}
	defer writer.Close()
	csv := csv.NewWriter(writer)

	dbh, err := getOracleConnection(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error in connecting to database %s", err), 2)
	}
	defer dbh.Close()
	log.Infof("started writing annotation in %s", outfile)
	var (
		strainId   string
		createdBy  string
		createdOn  time.Time
		modifiedOn time.Time
	)
	stmt := `
	SELECT dbxref.accession,sc.created_by, sc.date_created, sc.date_modified
	FROM CGM_DDB.STOCK_CENTER sc
	JOIN CGM_CHADO.DBXREF dbxref
	ON sc.dbxref_id = dbxref.dbxref_id
	`
	rows, err := dbh.Query(stmt)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("unable to run query %s", err), 2)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&strainId, &createdBy, &createdOn, &modifiedOn)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("unable to scan the next row %s", err), 2)
		}
		err = csv.Write(
			[]string{
				strainId,
				createdBy,
				createdOn.Format(layout),
				modifiedOn.Format(layout),
			},
		)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("unable to write csv row %s", err), 2)
		}
	}
	err = rows.Err()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("unable to close the rows %s", err), 2)
	}
	csv.Flush()
	if err := csv.Error(); err != nil {
		return cli.NewExitError(fmt.Sprintf("unable to finish csv writing %s", err), 2)
	}
	log.Infof("finished writing annotations to %s", outfile)
	return nil
}

func StockCenterAction(c *cli.Context) error {
	CreateRequiredFolder(c.String("config-folder"))
	CreateRequiredFolder(c.String("output-folder"))
	errChan := make(chan error)
	out := make(chan []byte)
	for _, scmd := range []string{"dictystrain", "dictyplasmid"} {
		yc := MakeSCConfig(c, scmd)
		CreateSCFolder(yc)
		conf := make(map[string]string)
		conf["config"] = yc
		conf["dir"] = c.String("output-folder")
		go RunDumpCmd(conf, scmd, errChan, out)
	}

	count := 1
	//curr := time.Now()
	for {
		select {
		case r := <-out:
			log.Printf("\nfinished the %s\n succesfully", string(r))
			count++
			if count > 2 {
				break
			}
		case err := <-errChan:
			log.Printf("\nError %s in running command\n", err)
			count++
			if count > 2 {
				break
			}
		default:
			time.Sleep(100 * time.Millisecond)
			//elapsed := time.Since(curr)
			//fmt.Printf("\r%d:%d:%d\t", int(elapsed.Hours()), int(elapsed.Minutes()), int(elapsed.Seconds()))
		}
	}
	return nil
}

func DscOrderAction(c *cli.Context) error {
	CreateRequiredFolder(outfolder)
	if err := exportDscOrders(c); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	return nil
}

func exportDscOrders(c *cli.Context) error {
	log := getLogger(c)
	mainCmd, err := exec.LookPath("modware-export")
	if err != nil {
		return fmt.Errorf("could not find binary %s", err)
	}
	subCmd := makeOrderExportCmd(c)
	log.Infof("running the command %s", strings.Join(subCmd, " "))
	_, err = exec.Command(mainCmd, subCmd...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error %s running the command %s", err, strings.Join(subCmd, " "))
	}
	log.Info("successfully ran the command")
	return nil
}

func makeOrderExportCmd(c *cli.Context) []string {
	return []string{
		"dscorders",
		"--dsn",
		c.String("legacy-dsn"),
		"-u",
		c.String("legacy-user"),
		"-p",
		c.String("legacy-password"),
		"--so",
		filepath.Join(c.String("output-folder"), "strain_orders.csv"),
		"--po",
		filepath.Join(c.String("output-folder"), "plasmid_orders.csv"),
	}
}
