package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"github.com/urfave/cli"
)

var ur = regexp.MustCompile("-")
var dbRgxp = regexp.MustCompile(`Dbxref=(.+)(;)?`)

func CanonicalGFF3Action(c *cli.Context) error {
	if !ValidateArgs(c) {
		return cli.NewExitError("one or more of required arguments are not provided", 2)
	}
	if !ValidateMultiArgs(c) {
		return cli.NewExitError("one or more of required arguments are not provided", 2)
	}
	// create config folder if not exists
	CreateRequiredFolder(c.String("config-folder"))

	errChan := make(chan error)
	out := make(chan []byte)
	mopt := map[string]map[string]string{
		"purpureum": {
			"genus":   "Dictyostelium",
			"species": "purpureum",
		},
		"pallidum": {
			"genus":   "Polysphondylium",
			"species": "pallidum PN500",
		},
		"pallidum_mitochondrial": {
			"genus":   "Polysphondylium",
			"species": "pallidum CK8",
		},
		"fasciculatum": {
			"genus":                 "Dictyostelium",
			"species":               "fasciculatum SH3",
			"exclude_mitochondrial": "1",
		},
		"fasciculatum_mitochondrial": {
			"genus":              "Dictyostelium",
			"species":            "fasciculatum SH3",
			"only_mitochondrial": "1",
		},
	}
	var subf string
	var name string
	for k, v := range mopt {
		if strings.Contains(k, "_") {
			sn := strings.Split(k, "_")
			subf = sn[0]
			name = fmt.Sprintf("%s_canonical_%s", sn[0], sn[1])
		} else {
			subf = k
			name = fmt.Sprintf("%s_canonical_core", k)
		}
		conf := MakeCustomConfigFile(c, name, subf)
		CreateFolderFromYaml(conf)
		opt := make(map[string]string)
		for ik, iv := range v {
			opt[ik] = iv
		}
		opt["config"] = conf
		go RunExportCmd(opt, "chado2canonicalgff3", errChan, out)
	}
	dc := MakeDictyConfigFile(c, "canonical_core", "discoideum")
	CreateFolderFromYaml(dc)
	dopt := map[string]string{"config": dc}
	go RunExportCmd(dopt, "chado2dictycanonicalgff3", errChan, out)

	count := 1
	curr := time.Now()
	for {
		select {
		case r := <-out:
			log.Printf("\nfinished the %s\n succesfully", string(r))
			count++
			if count > 6 {
				return nil
			}
		case err := <-errChan:
			log.Printf("\nError %s in running command\n", err)
			count++
			if count > 6 {
				return nil
			}
		default:
			time.Sleep(1000 * time.Millisecond)
			elapsed := time.Since(curr)
			fmt.Printf("\r%d:%d:%d\t", int(elapsed.Hours()), int(elapsed.Minutes()), int(elapsed.Seconds()))
		}
	}
	return nil
}

func ExtraGFF3Action(c *cli.Context) error {
	if !ValidateArgs(c) {
		return cli.NewExitError("one or more of required arguments are not provided", 2)
	}
	// create config folder if not exists
	CreateRequiredFolder(c.String("config-folder"))

	errChan := make(chan error)
	out := make(chan []byte)
	cmap := map[string]string{
		"chado2dictynoncodinggff3":      "canonical_noncoding",
		"chado2dictynoncanonicalgff3":   "noncanonical_seq_center",
		"chado2dictynoncanonicalv2gff3": "noncanonical_norepred",
		"chado2dictycuratedgff3":        "noncanonical_curated",
	}
	for k, v := range cmap {
		conf := map[string]string{"config": MakeDictyConfigFile(c, v, "discoideum")}
		CreateFolderFromYaml(conf["config"])
		go RunExportCmd(conf, k, errChan, out)
	}
	ao := []map[string]string{
		map[string]string{"organism": "dicty", "reference_type": "chromosome", "feature_type": "EST"},
		map[string]string{"organism": "dicty", "reference_type": "chromosome", "feature_type": "cDNA_clone"},
		map[string]string{"organism": "dicty", "reference_type": "chromosome", "feature_type": "databank_entry", "match_type": "nucleotide_match"},
	}
	for _, m := range ao {
		m["config"] = MakeDictyConfigFile(c, m["feature_type"], "discoideum")
		go RunExportCmd(m, "chado2alignmentgff3", errChan, out)
	}

	count := 1
	//curr := time.Now()
	for {
		select {
		case r := <-out:
			log.Printf("\nfinished the %s\n succesfully", string(r))
			count++
			if count > 7 {
				return nil
			}
		case err := <-errChan:
			log.Printf("\nError %s in running command\n", err)
			count++
			if count > 7 {
				return nil
			}
		default:
			time.Sleep(100 * time.Millisecond)
			//elapsed := time.Since(curr)
			//fmt.Printf("\r%d:%d:%d\t", int(elapsed.Hours()), int(elapsed.Minutes()), int(elapsed.Seconds()))
		}
	}
	return nil
}

func LiteratureAction(c *cli.Context) error {
	if !ValidateArgs(c) {
		return cli.NewExitError("one or more of required arguments are not provided", 2)
	}
	// create config folder if not exists
	for _, f := range []string{
		c.String("config-folder"),
		c.String("output-folder"),
		c.String("log-folder"),
	} {
		CreateRequiredFolder(f)
	}
	pconf := map[string]string{
		"config": MakeLiteatureConfig(c, "chadopub2bib"),
		"email":  c.String("email"),
		"output": filepath.Join(c.String("output-folder"), "dictytemp.bib"),
	}
	dconf := map[string]string{
		"conf":   MakeLiteatureConfig(c, "dictybib"),
		"output": filepath.Join(c.String("output-folder"), "dictybib.bib"),
		"input":  filepath.Join(c.String("output-folder"), "dictytemp.bib"),
	}
	gconf := map[string]string{
		"config": MakePub2BibConfig(c, "dictygenomespub"),
		"input":  filepath.Join(c.String("output-folder"), "dictygenomes_pubid.txt"),
	}
	nconf := map[string]string{
		"conf":   MakeLiteatureConfig(c, "dictynonpub"),
		"output": filepath.Join(c.String("output-folder"), "dictynonpub.bib"),
	}
	aconf := map[string]string{
		"conf":   MakeLiteatureConfig(c, "dictypubannotation"),
		"output": filepath.Join(c.String("output-folder"), "dictypubannotation.csv"),
	}

	wg := new(sync.WaitGroup)
	wg.Add(4)
	go RunLiteratureExportCmd(pconf, "chadopub2bib", wg)
	go RunTransformCmd(gconf, "pub2bib", gconf["input"], wg)
	go RunLiteratureExportCmd(nconf, "dictynonpub2bib", wg)
	go RunLiteratureExportCmd(aconf, "dictypubannotation", wg)
	wg.Wait()
	RunLiteratureUpdateCmd(dconf, "dictybib")
	return nil
}

func GeneAnnoAction(c *cli.Context) error {
	if !ValidateExtraArgs(c) {
		return cli.NewExitError("one or more of required arguments are not provided", 2)
	}
	// create config folder if not exists
	CreateRequiredFolder(c.String("config-folder"))

	errChan := make(chan error)
	out := make(chan []byte)
	conf := make(map[string]string)
	conf["config"] = MakeGeneralConfigFile(c, "genesummary", "csv")
	CreateFolderFromYaml(conf["config"])
	for _, param := range []string{"legacy-user", "legacy-password", "legacy-dsn"} {
		opt := ur.ReplaceAllString(param, "_")
		conf[opt] = c.String(param)
	}
	go RunExportCmd(conf, "chado2genesummary", errChan, out)

	for _, param := range []string{"public", "private"} {
		conf := make(map[string]string)
		conf["config"] = MakeGeneralConfigFile(c, param, "csv")
		conf["note"] = param
		go RunExportCmd(conf, "curatornotes", errChan, out)
	}

	conf2 := make(map[string]string)
	conf2["conf"] = MakeGeneralConfigFile(c, "coll2gene", "csv")
	go RunExportCmd(conf2, "colleague2gene", errChan, out)

	count := 1
	//curr := time.Now()
	for {
		select {
		case r := <-out:
			fmt.Printf("\nfinished the %s\n succesfully", string(r))
			count++
			if count > 4 {
				return nil
			}
		case err := <-errChan:
			fmt.Printf("\nError %s in running command\n", err)
			count++
			if count > 4 {
				return nil
			}
		default:
			time.Sleep(100 * time.Millisecond)
			//elapsed := time.Since(curr)
			//fmt.Printf("\r%d:%d:%d\t", int(elapsed.Hours()), int(elapsed.Minutes()), int(elapsed.Seconds()))
		}
	}
	return nil
}

func RunExportCmd(opt map[string]string, subcmd string, errChan chan<- error, out chan<- []byte) {
	p := make([]string, 0)
	p = append(p, subcmd)
	for k, v := range opt {
		p = append(p, fmt.Sprint("--", k), v)
	}
	cmdline := strings.Join(p, " ")
	log.Printf("going to run %s\n", cmdline)
	b, err := exec.Command("modware-export", p...).CombinedOutput()
	if err != nil {
		errChan <- fmt.Errorf("Status %s message %s for cmdline %s\n", err.Error(), string(b), cmdline)
		return
	}
	out <- []byte(cmdline)
}

func RunLiteratureExportCmd(opt map[string]string, subcmd string, wg *sync.WaitGroup) {
	defer wg.Done()
	p := make([]string, 0)
	p = append(p, subcmd)
	for k, v := range opt {
		p = append(p, fmt.Sprint("--", k), v)
	}
	cmdline := strings.Join(p, " ")
	log.Printf("going to run %s\n", cmdline)
	b, err := exec.Command("modware-export", p...).CombinedOutput()
	if err != nil {
		fmt.Printf("Status %s message %s\n", err.Error(), string(b))
	} else {
		log.Printf("finished running %s\n", cmdline)
	}
}

func RunLiteratureUpdateCmd(opt map[string]string, subcmd string) {
	p := make([]string, 0)
	p = append(p, subcmd)
	for k, v := range opt {
		p = append(p, fmt.Sprint("--", k), v)
	}
	cmdline := strings.Join(p, " ")
	log.Printf("going to run %s\n", cmdline)
	b, err := exec.Command("modware-update", p...).CombinedOutput()
	if err != nil {
		fmt.Printf("Status %s message %s\n", err.Error(), string(b))
	} else {
		log.Printf("finished running %s\n", cmdline)
	}

}

func RunTransformCmd(opt map[string]string, subcmd string, file string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Write list of pubmed ids to the input file
	err := ioutil.WriteFile(file, []byte("13319664\n15867862\n17246401\n"), 0644)
	if err != nil {
		fmt.Printf("Error creating input file %s\n", file)
		return
	}
	p := make([]string, 0)
	p = append(p, subcmd)
	for k, v := range opt {
		p = append(p, fmt.Sprint("--", k), v)
	}
	cmdline := strings.Join(p, " ")
	log.Printf("going to run %s\n", cmdline)
	b, err := exec.Command("modware-transform", p...).CombinedOutput()
	if err != nil {
		fmt.Printf("Status %s message %s\n", err.Error(), string(b))
	} else {
		log.Printf("finished running %s\n", cmdline)
	}

}

func RunDumpCmd(opt map[string]string, subcmd string, errChan chan<- error, out chan<- []byte) {
	p := make([]string, 0)
	p = append(p, subcmd)
	for k, v := range opt {
		p = append(p, fmt.Sprint("--", k), v)
	}
	cmdline := strings.Join(p, " ")
	log.Printf("going to run %s\n", cmdline)
	b, err := exec.Command("modware-dump", p...).CombinedOutput()
	if err != nil {
		errChan <- fmt.Errorf("Status %s message %s\n", err.Error(), string(b))
		return
	}
	out <- []byte(cmdline)
}

func RunLiteraturePipeCmd(fopt map[string]string, sopt map[string]string, fscmd string, scmd string, wg *sync.WaitGroup) {
	defer wg.Done()
	fp := make([]string, 0)
	fp = append(fp, fscmd)
	for k, v := range fopt {
		fp = append(fp, fmt.Sprint("--", k), v)
	}
	fcmdline := strings.Join(fp, " ")

	sp := make([]string, 0)
	sp = append(sp, scmd)
	for k, v := range sopt {
		sp = append(sp, fmt.Sprint("--", k), v)
	}
	scmdline := strings.Join(sp, " ")

	fc := exec.Command("modware-export", fp...)
	sc := exec.Command("modware-update", sp...)

	// http://golang.org/pkg/io/#Pipe
	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()
	var fb bytes.Buffer
	var sb bytes.Buffer
	// first stdout go to writer
	// capture the errors of command if any
	fc.Stdout = writer
	fc.Stderr = &fb
	// second stdin is reader
	sc.Stdin = reader
	sc.Stderr = &sb

	// start and wait for the commands
	fmt.Printf("Going to run command %s\n", fcmdline)
	if err := fc.Start(); err != nil {
		fmt.Printf("Error starting first command: %s\n", fb.String())
		return
	}
	fmt.Printf("Going to command %s\n", scmdline)
	if err := sc.Start(); err != nil {
		fmt.Printf("Error starting second command: %s\n", sb.String())
		return
	}
	if err := fc.Wait(); err != nil {
		fmt.Printf("Error running first command: %s\n", fb.String())
		return
	}
	if err := sc.Wait(); err != nil {
		fmt.Printf("Error running second command: %s\n", sb.String())
		wg.Done()
		return
	}
	fmt.Println("finished both commands succesfully")
}

func DbxrefCleanUpAction(c *cli.Context) error {
	if err := ValidateCleanUpArgs(c); err != nil {
		cli.NewExitError(err.Error(), 2)
	}
	in, err := os.Open(c.String("input"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in opening file %s\n", err),
			2,
		)
	}
	defer in.Close()
	out, err := os.Create(c.String("output"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in writing file %s\n", err),
			2,
		)
	}
	defer out.Close()

	r := bufio.NewReader(in)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return cli.NewExitError(
					fmt.Sprintf("error in reading file %s\n", err),
					2,
				)
			}

		}
		if strings.Contains(line, "Dbxref") {
			gff3Slice := strings.Split(line, "\t")
			nonxref, dbxref := splitDbxref(gff3Slice[8])
			cdbxref := replaceDbxref(dbxref, c.StringSlice("db-name"))
			if len(cdbxref) == 0 {
				gff3Slice[8] = strings.TrimSuffix(nonxref, ";")
			} else {
				gff3Slice[8] = fmt.Sprintf("%s%s", nonxref, cdbxref)
			}
			fmt.Fprintf(out, "%s\n", strings.Join(gff3Slice, "\t"))
		} else {
			out.WriteString(line)
		}
	}
	return nil
}

func splitDbxref(attr string) (string, string) {
	var nonxref bytes.Buffer
	var dbxref bytes.Buffer
	for _, s := range strings.Split(attr, ";") {
		if strings.HasPrefix(s, "Dbxref") {
			dbxref.WriteString(s)
			continue
		}
		nonxref.WriteString(strings.TrimSuffix(s, "\n") + ";")
	}
	return nonxref.String(), dbxref.String()
}

func hasDb(dbxref string, db []string) bool {
	for _, n := range db {
		if strings.Contains(dbxref, n) {
			return true
		}
	}
	return false
}

func hasDbPrefix(dbxref string, db []string) bool {
	for _, n := range db {
		if strings.HasPrefix(dbxref, n) {
			return true
		}
	}
	return false
}

func replaceDbxref(dbxref string, db []string) string {
	if !hasDb(dbxref, db) {
		return dbxref
	}

	sm := dbRgxp.FindStringSubmatch(dbxref)
	if sm == nil {
		return ""
	}
	var fDbxref []string
	for _, xref := range strings.Split(sm[1], ",") {
		if !hasDbPrefix(xref, db) {
			fDbxref = append(fDbxref, xref)
		}
	}
	if len(fDbxref) == 0 {
		return ""
	}
	return fmt.Sprintf("Dbxref=%s", strings.TrimSuffix(strings.Join(fDbxref, ","), "\n"))
}

func SplitPolypeptideAction(c *cli.Context) error {
	if err := ValidatePolypetideArgs(c); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	in, err := os.Open(c.String("input"))
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in opening file %s\n", err),
			2,
		)
	}
	defer in.Close()
	genome, poly := MakeOutputName(c.String("input"))
	gr, err := os.Create(genome)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in opening %s file  for wriring %s\n", genome, err),
			2,
		)
	}
	defer gr.Close()
	pr, err := os.Create(poly)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in opening %s file  for wriring %s\n", genome, err),
			2,
		)
	}
	defer pr.Close()
	r := bufio.NewReader(in)
	fmt.Fprintf(pr, "##%s\t3\n", "gff-version")
	seenFasta := false
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return cli.NewExitError(
					fmt.Sprintf("error in reading file %s\n", genome, err),
					2,
				)
			}
		}
		switch {
		case seenFasta:
			fmt.Fprint(gr, line)
		case strings.HasPrefix(line, "###"):
			fmt.Fprint(gr, line)
			seenFasta = true
		case strings.HasPrefix(line, "##"):
			fmt.Fprint(gr, line)
		case isPolyPeptide(line):
			fmt.Fprint(pr, line)
		default:
			fmt.Fprint(gr, line)
		}
	}
	return nil
}

func isPolyPeptide(line string) bool {
	t := strings.Split(line, "\t")
	if t[2] == "polypeptide" {
		return true
	}
	return false
}

func MakeOutputName(path string) (string, string) {
	prefix := strings.Split(filepath.Base(path), ".")[0]
	dir := filepath.Dir(path)
	return filepath.Join(dir, fmt.Sprintf("%s_no_poly.gff3", prefix)),
		filepath.Join(dir, fmt.Sprintf("%s_poly.gff3", prefix))
}
