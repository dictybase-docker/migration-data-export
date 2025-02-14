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
	CvpropID int    `db:"CVPROP_ID"` // PRIMARY KEY
	Value    string `db:"VALUE"`
}

type FeatureCvtermprop struct {
	FeatureCvtermpropID int    `db:"FEATURE_CVTERMPROP_ID"`
	Value               string `db:"VALUE"`
}

type FeaturePubprop struct {
	FeaturePubpropID int    `db:"FEATURE_PUBPROP_ID"`
	Value            string `db:"VALUE"`
}

type Organism struct {
	OrganismID int    `db:"ORGANISM_ID"`
	Comment    string `db:"COMMENT_"`
}

type Pubprop struct {
	PubpropID int    `db:"PUBPROP_ID"`
	Value     string `db:"VALUE"`
}

type Analysis struct {
	AnalysisID  int    `db:"ANALYSIS_ID"`
	Description string `db:"DESCRIPTION"`
	Sourceuri   string `db:"SOURCEURI"`
}

type ChadoLogs struct {
	ChadoLogsID int    `db:"CHADO_LOGS_ID"`
	NewValue    string `db:"NEW_VALUE"`
	OldValue    string `db:"OLD_VALUE"`
}

type Dbxref struct {
	DbxrefID    int    `db:"DBXREF_ID"`
	Description string `db:"DESCRIPTION"`
}

type Featureloc struct {
	FeaturelocID int    `db:"FEATURELOC_ID"`
	ResidueInfo  string `db:"RESIDUE_INFO"`
}

type Cv struct {
	CvID       int    `db:"CV_ID"`
	Definition string `db:"DEFINITION"`
}

type Feature struct {
	FeatureID int    `db:"FEATURE_ID"`
	Residues  string `db:"RESIDUES"`
}

type Pub struct {
	PubID       int    `db:"PUB_ID"`
	Title       string `db:"TITLE"`
	Volumetitle string `db:"VOLUMETITLE"`
}

type Arraydesign struct {
	ArraydesignID     int    `db:"ARRAYDESIGN_ID"`
	ArrayDimensions   string `db:"ARRAY_DIMENSIONS"`
	Description       string `db:"DESCRIPTION"`
	ElementDimensions string `db:"ELEMENT_DIMENSIONS"`
	Version           string `db:"VERSION"`
}

type Chadoprop struct {
	ChadopropID int    `db:"CHADOPROP_ID"`
	Value       string `db:"VALUE"`
}

type Cvterm struct {
	CvtermID   int    `db:"CVTERM_ID"`
	Definition string `db:"DEFINITION"`
}

type Featureprop struct {
	FeaturepropID int    `db:"FEATUREPROP_ID"`
	Value         string `db:"VALUE"`
}

type FeatureRelationship struct {
	FeatureRelationshipID int    `db:"FEATURE_RELATIONSHIP_ID"`
	Value                 string `db:"VALUE"`
}

type MvGeneDescription struct {
	MvGeneDescriptionID int    `db:"MV_GENE_DESCRIPTION_ID"`
	Description         string `db:"DESCRIPTION"`
}

type Phenotype struct {
	PhenotypeID int    `db:"PHENOTYPE_ID"`
	Value       string `db:"VALUE"`
}
