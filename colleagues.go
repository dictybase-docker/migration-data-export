package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

func validateColleagues(c *cli.Context) error {
	if !ValidateExtraArgs(c) {
		return cli.NewExitError("one or more of required arguments are not provided", 2)
	}
	return nil
}

func ColleaguesAction(c *cli.Context) error {
	if err := CreateFolder(c.String("output-folder")); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	if err := exportColleagues(c); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	return nil
}

func exportColleagues(c *cli.Context) error {
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

func makeColleaguesExportCmd(c *cli.Context) []string {
	return []string{
		"colleague",
		"--dsn",
		c.String("legacy-dsn"),
		"-u",
		c.String("legacy-user"),
		"-p",
		c.String("legacy-password"),
		"--crel",
		filepath.Join(c.String("output-folder"), "user_relations.csv"),
		"--cout",
		filepath.Join(c.String("output-folder"), "users.csv"),
	}
}
