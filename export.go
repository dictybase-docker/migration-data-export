package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "Command line wrapper to export data from Oracle chado using modware-loader"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level",
			Usage: "logging level of the appliction with six choices, debug, info, warn, error, fatal and panic",
			Value: "warn",
		},
		cli.StringFlag{
			Name:  "log-format",
			Usage: "format of the logging out, either of json or text",
			Value: "text",
		},
		cli.StringSliceFlag{
			Name:  "hooks",
			Usage: "hook names for sending log in addition to stderr",
			Value: &cli.StringSlice{},
		},
		cli.StringFlag{
			Name:   "slack-channel",
			EnvVar: "SLACK_CHANNEL",
			Usage:  "Slack channel where the log will be posted",
		},
		cli.StringFlag{
			Name:   "slack-url",
			EnvVar: "SLACK_URL",
			Usage:  "Slack webhook url[required if slack channel is provided]",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "canonicalgff3",
			Usage:  "Export the canonical gff3 of all the dictyostelids",
			Action: CanonicalGFF3Action,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output-folder, of",
					Usage: "Output folder",
					Value: "/data/gff3",
				},
				cli.StringFlag{
					Name:  "log-folder, lf",
					Usage: "Log folder",
					Value: "/log/gff3",
				},
				cli.StringFlag{
					Name:  "config-folder, cf",
					Usage: "Folder for config files",
					Value: "/config/gff3",
				},
				cli.StringFlag{
					Name:   "dsn",
					Usage:  "dsn for oracle database server [required]",
					EnvVar: "ORACLE_DSN",
				},
				cli.StringFlag{
					Name:   "user, u",
					Usage:  "User name for oracle database [required]",
					EnvVar: "ORACLE_USER",
				},
				cli.StringFlag{
					Name:   "password, p",
					Usage:  "Password for oracle database[required]",
					EnvVar: "ORACLE_PASS",
				},
				cli.StringFlag{
					Name:   "muser, mu",
					Usage:  "User name for multigenome oracle database [required]",
					EnvVar: "MULTI_ORACLE_USER",
				},
				cli.StringFlag{
					Name:   "mpassword, mp",
					Usage:  "Password for multigenome oracle database[required]",
					EnvVar: "MULTI_ORACLE_PASS",
				},
			},
		},
		{
			Name:   "extradictygff3",
			Usage:  "Export the additional gff3 of all the D.discoideum",
			Action: ExtraGFF3Action,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output-folder, of",
					Usage: "Output folder",
					Value: "/data/gff3",
				},
				cli.StringFlag{
					Name:  "log-folder, lf",
					Usage: "Log folder",
					Value: "/log/gff3",
				},
				cli.StringFlag{
					Name:  "config-folder, cf",
					Usage: "Folder for config files",
					Value: "/config/gff3",
				},
				cli.StringFlag{
					Name:   "dsn",
					Usage:  "dsn for oracle database server [required]",
					EnvVar: "ORACLE_DSN",
				},
				cli.StringFlag{
					Name:   "user, u",
					Usage:  "User name for oracle database [required]",
					EnvVar: "ORACLE_USER",
				},
				cli.StringFlag{
					Name:   "password, p",
					Usage:  "Password for oracle database[required]",
					EnvVar: "ORACLE_PASS",
				},
			},
		},
		{
			Name:   "geneannotation",
			Usage:  "Export annotations associated with gene models",
			Action: GeneAnnoAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output-folder, of",
					Usage: "Output folder",
					Value: "/data/annotation",
				},
				cli.StringFlag{
					Name:  "log-folder, lf",
					Usage: "Log folder",
					Value: "/log/annotation",
				},
				cli.StringFlag{
					Name:  "config-folder, cf",
					Usage: "Folder for config files",
					Value: "/config/annotation",
				},
				cli.StringFlag{
					Name:   "dsn",
					Usage:  "dsn for oracle database server [required]",
					EnvVar: "ORACLE_DSN",
				},
				cli.StringFlag{
					Name:   "user, u",
					Usage:  "User name for oracle database [required]",
					EnvVar: "ORACLE_USER",
				},
				cli.StringFlag{
					Name:   "password, p",
					Usage:  "Password for oracle database[required]",
					EnvVar: "ORACLE_PASS",
				},
				cli.StringFlag{
					Name:   "legacy-user",
					Usage:  "User name for legacy oracle database[required]",
					EnvVar: "LEGACY_USER",
				},
				cli.StringFlag{
					Name:   "legacy-password",
					Usage:  "Password for legacy oracle database [required]",
					EnvVar: "LEGACY_PASS",
				},
				cli.StringFlag{
					Name:   "legacy-dsn",
					Usage:  "dsn for legacy oracle database [required]",
					EnvVar: "LEGACY_DSN",
				},
			},
		},
		{
			Name:     "dsc-orders",
			Usage:    "Export dictybase stock center orders",
			Category: "dsc",
			Action:   DscOrderAction,
			Before:   validateDscOrder,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output-folder, of",
					Usage: "Output folder",
					Value: "/data/stockcenter",
				},
				cli.StringFlag{
					Name:   "legacy-user",
					Usage:  "User name for legacy oracle database[required]",
					EnvVar: "LEGACY_USER",
				},
				cli.StringFlag{
					Name:   "legacy-password",
					Usage:  "Password for legacy oracle database [required]",
					EnvVar: "LEGACY_PASS",
				},
				cli.StringFlag{
					Name:   "legacy-dsn",
					Usage:  "dsn for legacy oracle database [required]",
					EnvVar: "LEGACY_DSN",
				},
			},
		},
		{
			Name:     "dsc-annotations",
			Usage:    "Export annotations",
			Category: "dsc",
			Before:   validateDsc,
			Action:   StockCenterAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output-folder, of",
					Usage: "Output folder",
					Value: "/data/stockcenter",
				},
				cli.StringFlag{
					Name:  "log-folder, lf",
					Usage: "Log folder",
					Value: "/log/stockcenter",
				},
				cli.StringFlag{
					Name:  "config-folder, cf",
					Usage: "Folder for config files",
					Value: "/config/stockcenter",
				},
				cli.StringFlag{
					Name:   "host",
					Usage:  "oracle database host name",
					EnvVar: "ORACLE_HOST",
				},
				cli.StringFlag{
					Name:   "port",
					Usage:  "oracle database port",
					EnvVar: "ORACLE_PORT",
				},
				cli.StringFlag{
					Name:   "sid",
					Usage:  "oracle database sid",
					EnvVar: "ORACLE_SID",
				},
				cli.StringFlag{
					Name:   "dsn",
					Usage:  "dsn for oracle database server [required]",
					EnvVar: "ORACLE_DSN",
				},
				cli.StringFlag{
					Name:   "user, u",
					Usage:  "User name for oracle database [required]",
					EnvVar: "ORACLE_USER",
				},
				cli.StringFlag{
					Name:   "password, p",
					Usage:  "Password for oracle database[required]",
					EnvVar: "ORACLE_PASS",
				},
				cli.StringFlag{
					Name:   "legacy-user",
					Usage:  "User name for legacy oracle database[required]",
					EnvVar: "LEGACY_USER",
				},
				cli.StringFlag{
					Name:   "legacy-password",
					Usage:  "Password for legacy oracle database [required]",
					EnvVar: "LEGACY_PASS",
				},
				cli.StringFlag{
					Name:   "legacy-dsn",
					Usage:  "dsn for legacy oracle database [required]",
					EnvVar: "LEGACY_DSN",
				},
			},
		},
		{
			Name:     "dsc-users",
			Usage:    "Export a map between users and their annotations",
			Category: "dsc",
			Action:   DscUsersAction,
			Before:   validateDscUsers,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "output-folder, of",
					Usage:       "Output folder of the data files",
					Destination: &outfolder,
					Value:       "/data/stockcenter",
				},
				cli.StringFlag{
					Name:   "host",
					Usage:  "oracle database host name[required]",
					EnvVar: "ORACLE_HOST",
				},
				cli.StringFlag{
					Name:   "port",
					Usage:  "oracle database port[required]",
					EnvVar: "ORACLE_PORT",
				},
				cli.StringFlag{
					Name:   "sid",
					Usage:  "oracle database sid[required]",
					EnvVar: "ORACLE_SID",
				},
				cli.StringFlag{
					Name:   "user, u",
					Usage:  "User name for oracle database [required]",
					EnvVar: "ORACLE_USER",
				},
				cli.StringFlag{
					Name:   "password, p",
					Usage:  "Password for oracle database[required]",
					EnvVar: "ORACLE_PASS",
				},
			},
		},
		{
			Name:   "literature",
			Usage:  "Export the literature and annotations",
			Action: LiteratureAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output-folder, of",
					Usage: "Output folder",
					Value: "/data/literature",
				},
				cli.StringFlag{
					Name:  "log-folder, lf",
					Usage: "Log folder",
					Value: "/log/literature",
				},
				cli.StringFlag{
					Name:  "config-folder, cf",
					Usage: "Folder for config files",
					Value: "/config/literature",
				},
				cli.StringFlag{
					Name:   "dsn",
					Usage:  "dsn for oracle database server [required]",
					EnvVar: "ORACLE_DSN",
				},
				cli.StringFlag{
					Name:   "user, u",
					Usage:  "User name for oracle database [required]",
					EnvVar: "ORACLE_USER",
				},
				cli.StringFlag{
					Name:   "password, p",
					Usage:  "Password for oracle database[required]",
					EnvVar: "ORACLE_PASS",
				},
				cli.StringFlag{
					Name:   "email",
					Usage:  "Email to use for ncbi utils[required]",
					EnvVar: "EUTILS_EMAIL",
				},
			},
		},
		{
			Name:  "clean-dbxref",
			Usage: "Remove dbxref attribute(s) from gff3 file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Name of output gff3 file, required",
				},
				cli.StringFlag{
					Name:  "input, i",
					Usage: "Name of the input gff3 file, required",
				},
				cli.StringSliceFlag{
					Name:  "n, db-name",
					Usage: "database name which will be removed along with the accession/id",
				},
			},
			Action: DbxrefCleanUpAction,
		},
		{
			Name:  "split-polypeptide",
			Usage: "Split polypeptide features to a separate file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Usage: "Name of the input gff3 file, required",
				},
			},
			Action: SplitPolypeptideAction,
		},
	}
	app.Run(os.Args)
}
