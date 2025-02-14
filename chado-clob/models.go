package main

import "strings"

// GetStructForTable returns the appropriate struct type for a given table name
// Returns the struct instance and a boolean indicating if the table was found
func GetStructForTable(tableName string) (interface{}, bool) {
	switch strings.ToLower(tableName) {
	case "cvprop":
		return &Cvprop{}, true
	case "feature_cvtermprop":
		return &FeatureCvtermprop{}, true
	case "feature_pubprop":
		return &FeaturePubprop{}, true
	case "organism":
		return &Organism{}, true
	case "pubprop":
		return &Pubprop{}, true
	case "analysis":
		return &Analysis{}, true
	case "chadologs":
		return &ChadoLogs{}, true
	case "dbxref":
		return &Dbxref{}, true
	case "featureloc":
		return &Featureloc{}, true
	case "cv":
		return &Cv{}, true
	case "feature":
		return &Feature{}, true
	case "pub":
		return &Pub{}, true
	case "arraydesign":
		return &Arraydesign{}, true
	case "chadoprop":
		return &Chadoprop{}, true
	case "cvterm":
		return &Cvterm{}, true
	case "featureprop":
		return &Featureprop{}, true
	case "feature_relationship":
		return &FeatureRelationship{}, true
	case "mvgenedescription":
		return &MvGeneDescription{}, true
	case "phenotype":
		return &Phenotype{}, true
	default:
		return nil, false
	}
}

type Cvprop struct {
	CvpropID int    // PRIMARY KEY
	Value    string // from VALUE column
}

type FeatureCvtermprop struct {
	FeatureCvtermpropID int
	Value               string
}

type FeaturePubprop struct {
	FeaturePubpropID int
	Value            string
}

type Organism struct {
	OrganismID int
	Comment    string `db:"COMMENT_"`
}

type Pubprop struct {
	PubpropID int
	Value     string
}

type Analysis struct {
	AnalysisID  int
	Description string
	Sourceuri   string
}

type ChadoLogs struct {
	ChadoLogsID int
	NewValue    string
	OldValue    string
}

type Dbxref struct {
	DbxrefID    int
	Description string
}

type Featureloc struct {
	FeaturelocID int
	ResidueInfo  string
}

type Cv struct {
	CvID       int
	Definition string
}

type Feature struct {
	FeatureID int
	Residues  string
}

type Pub struct {
	PubID       int
	Title       string
	Volumetitle string
}

type Arraydesign struct {
	ArraydesignID     int
	ArrayDimensions   string
	Description       string
	ElementDimensions string
	Version           string
}

type Chadoprop struct {
	ChadopropID int
	Value       string
}

type Cvterm struct {
	CvtermID   int
	Definition string
}

type Featureprop struct {
	FeaturepropID int
	Value         string
}

type FeatureRelationship struct {
	FeatureRelationshipID int
	Value                 string
}

type MvGeneDescription struct {
	MvGeneDescriptionID int
	Description         string
}

type Phenotype struct {
	PhenotypeID int
	Value       string
}
