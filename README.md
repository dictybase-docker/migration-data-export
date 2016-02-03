# Migration data export
This is a source repository for [docker](http://docker.io) image to run
all data export tasks of dictybase overhaul project.

## Usage

`Build`

```docker build --rm=true -t dictybase/migration-data-export .```

`Run`

### Command line

```
docker run --rm dictybase/migration-data-export app -h

NAME:
   Command line wrapper to export data from Oracle chado using modware-loader - A new cli application

USAGE:
   Command line wrapper to export data from Oracle chado using modware-loader [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   canonicalgff3	Export the canonical gff3 of all the dictyostelids
   extradictygff3	Export the additional gff3 of all the D.discoideum
   geneannotation	Export annotations associated with gene models
   stockcenter		Export strains, plasmids and all assoicated annotations related to stock center
   literature		Export the literature and annotations
   help, h		Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
   
