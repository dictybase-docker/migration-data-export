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

	"gopkg.in/codegangsta/cli.v1"
)

var ur = regexp.MustCompile("-")
var dbRgxp = regexp.MustCompile("Dbxref")

func CanonicalGFF3Action(c *cli.Context) {
	if !ValidateArgs(c) {
		log.Fatal("one or more of required arguments are not provided")
	}
	if !ValidateMultiArgs(c) {
		log.Fatal("one or more of required arguments are not provided")
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
				return
			}
		case err := <-errChan:
			log.Printf("\nError %s in running command\n", err)
			count++
			if count > 6 {
				return
			}
		default:
			time.Sleep(1000 * time.Millisecond)
			elapsed := time.Since(curr)
			fmt.Printf("\r%d:%d:%d\t", int(elapsed.Hours()), int(elapsed.Minutes()), int(elapsed.Seconds()))
		}
	}
}

func ExtraGFF3Action(c *cli.Context) {
	if !ValidateArgs(c) {
		log.Fatal("one or more of required arguments are not provided")
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
				return
			}
		case err := <-errChan:
			log.Printf("\nError %s in running command\n", err)
			count++
			if count > 7 {
				return
			}
		default:
			time.Sleep(100 * time.Millisecond)
			//elapsed := time.Since(curr)
			//fmt.Printf("\r%d:%d:%d\t", int(elapsed.Hours()), int(elapsed.Minutes()), int(elapsed.Seconds()))
		}
	}
}

func StockCenterAction(c *cli.Context) {
	if !ValidateArgs(c) {
		log.Fatal("one or more of required arguments are not provided")
	}
	if !ValidateExtraArgs(c) {
		log.Fatal("one or more of required arguments are not provided")
	}
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
				return
			}
		case err := <-errChan:
			log.Printf("\nError %s in running command\n", err)
			count++
			if count > 2 {
				return
			}
		default:
			time.Sleep(100 * time.Millisecond)
			//elapsed := time.Since(curr)
			//fmt.Printf("\r%d:%d:%d\t", int(elapsed.Hours()), int(elapsed.Minutes()), int(elapsed.Seconds()))
		}
	}
}

func LiteratureAction(c *cli.Context) {
	if !ValidateArgs(c) {
		log.Fatal("one or more of required arguments are not provided")
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
}

func GeneAnnoAction(c *cli.Context) {
	if !ValidateExtraArgs(c) {
		log.Fatal("one or more of required arguments are not provided")
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
				return
			}
		case err := <-errChan:
			fmt.Printf("\nError %s in running command\n", err)
			count++
			if count > 4 {
				return
			}
		default:
			time.Sleep(100 * time.Millisecond)
			//elapsed := time.Since(curr)
			//fmt.Printf("\r%d:%d:%d\t", int(elapsed.Hours()), int(elapsed.Minutes()), int(elapsed.Seconds()))
		}
	}
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

func DbxrefCleanUpAction(c *cli.Context) {
	if err := ValidateCleanUpArgs(c); err != nil {
		log.Fatal(err)
	}
	nRgxp, err := regexp.Compile(c.String("db-name") + `:\w+`)
	if err != nil {
		log.Fatalln("error in creating regexp for db name")
	}
	in, err := os.Open(c.String("input"))
	if err != nil {
		log.Fatalf("error in opening file %s\n", err)
	}
	defer in.Close()
	out, err := os.Create(c.String("output"))
	if err != nil {
		log.Fatalf("error in writing file %s\n", err)
	}
	defer out.Close()

	r := bufio.NewReader(in)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else {
			log.Fatalln(err)
		}
		if dbRgxp.MatchString(line) {
			gff3Slice := strings.Split(line, "\t")
			clnStr := nRgxp.ReplaceAllString(gff3Slice[8], "")
			gff3Slice[8] = clnStr
			fmt.Fprintf(out, "%s", strings.Join(gff3Slice, "\t"))
		} else {
			fmt.Fprintf(out, "%s", line)
		}
	}
}
