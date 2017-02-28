package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v1"
)

type Pub2BibConfig struct {
	Output    string `yaml:"output"`
	XMLOutput string `yaml:"xml_output"`
	Email     string `yaml:"email"`
	LogLevel  string `yaml:"log_level"`
	LogFile   string `yaml:"logfile"`
}

type LiteratureConfig struct {
	Dsn      string `yaml:"dsn"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	LogLevel string `yaml:"log_level"`
	LogFile  string `yaml:"logfile"`
}

type GeneralConfig struct {
	Dsn      string `yaml:"dsn"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Output   string `yaml:"output"`
	LogLevel string `yaml:"log_level"`
	LogFile  string `yaml:"logfile"`
}

type GFF3Config struct {
	Dsn      string `yaml:"dsn"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Output   string `yaml:"output"`
	LogLevel string `yaml:"log_level"`
	LogFile  string `yaml:"logfile"`
}

type StockCenterConfig struct {
	Dsn       string `yaml:"dsn"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Ldsn      string `yaml:"legacy_dsn"`
	Luser     string `yaml:"legacy_user"`
	Lpassword string `yaml:"legacy_password"`
	LogLevel  string `yaml:"log_level"`
	LogFile   string `yaml:"logfile"`
	Dir       string `yaml:"dir"`
}

func CreateYamlFile(in interface{}, c *cli.Context, name string) string {
	b, err := yaml.Marshal(in)
	if err != nil {
		log.Fatal(err)
	}
	p := filepath.Join(c.String("config-folder"), fmt.Sprint(name, ".yaml"))
	cf, err := os.Create(p)
	defer cf.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprint(cf, string(b))
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func MakePub2BibConfig(c *cli.Context, name string) string {
	pconf := Pub2BibConfig{
		Output:    filepath.Join(c.String("output-folder"), fmt.Sprint(name, ".bib")),
		XMLOutput: filepath.Join(c.String("output-folder"), fmt.Sprint(name, ".xml")),
		Email:     c.String("email"),
		LogLevel:  "info",
		LogFile:   filepath.Join(c.String("log-folder"), fmt.Sprint(name, ".log")),
	}
	return CreateYamlFile(pconf, c, name)
}

func MakeCustomConfigFile(c *cli.Context, name string, subfolder string) string {
	gconf := GFF3Config{
		Dsn:      c.String("dsn"),
		User:     c.String("muser"),
		Password: c.String("mpassword"),
		Output:   filepath.Join(c.String("output-folder"), subfolder, fmt.Sprint(name, ".gff3")),
		LogLevel: "debug",
		LogFile:  filepath.Join(c.String("log-folder"), subfolder, fmt.Sprint(name, ".log")),
	}
	return CreateYamlFile(gconf, c, name)
}

func MakeSCConfig(c *cli.Context, name string) string {
	gconf := StockCenterConfig{
		Dsn:       c.String("dsn"),
		User:      c.String("user"),
		Password:  c.String("password"),
		Ldsn:      c.String("legacy-dsn"),
		Lpassword: c.String("legacy-password"),
		Luser:     c.String("legacy-user"),
		LogFile:   filepath.Join(c.String("log-folder"), fmt.Sprint(name, ".log")),
		LogLevel:  "info",
	}
	return CreateYamlFile(gconf, c, name)
}

func MakeLiteatureConfig(c *cli.Context, name string) string {
	gconf := LiteratureConfig{
		Dsn:      c.String("dsn"),
		User:     c.String("user"),
		Password: c.String("password"),
		LogFile:  filepath.Join(c.String("log-folder"), fmt.Sprint(name, ".log")),
		LogLevel: "info",
	}
	return CreateYamlFile(gconf, c, name)
}

func MakeDictyConfigFile(c *cli.Context, name string, subfolder string) string {
	gconf := GFF3Config{
		Dsn:      c.String("dsn"),
		User:     c.String("user"),
		Password: c.String("password"),
		Output:   filepath.Join(c.String("output-folder"), subfolder, fmt.Sprint(name, ".gff3")),
		LogFile:  filepath.Join(c.String("log-folder"), subfolder, fmt.Sprint(name, ".log")),
		LogLevel: "info",
	}
	return CreateYamlFile(gconf, c, name)
}

func MakeGeneralConfigFile(c *cli.Context, name string, ext string) string {
	gconf := GeneralConfig{
		Dsn:      c.String("dsn"),
		User:     c.String("user"),
		Password: c.String("password"),
		Output:   filepath.Join(c.String("output-folder"), fmt.Sprint(name, ".", ext)),
		LogFile:  filepath.Join(c.String("log-folder"), fmt.Sprint(name, ".log")),
		LogLevel: "info",
	}
	return CreateYamlFile(gconf, c, name)
}

func MakeConfigFile(c *cli.Context, name string) string {
	gconf := GFF3Config{
		Dsn:      c.String("dsn"),
		User:     c.String("user"),
		Password: c.String("password"),
		Output:   filepath.Join(c.String("output-folder"), fmt.Sprint(name, ".gff3")),
		LogFile:  filepath.Join(c.String("log-folder"), fmt.Sprint(name, ".log")),
		LogLevel: "info",
	}
	return CreateYamlFile(gconf, c, name)
}
