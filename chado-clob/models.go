package main

import (
	"database/sql"
	"strings"
)

// GetStructForTable returns the appropriate struct type for a given table name
// Returns the struct instance and a boolean indicating if the table was found
func GetStructForTable(tableName string) (interface{}, bool) {
	switch strings.ToLower(tableName) {
	case "organismprop":
		return &Organismprop{}, true
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
	case "phenotype":
		return &Phenotype{}, true
	default:
		return nil, false
	}
}

type Cvprop struct {
	CvpropID int            `db:"CVPROP_ID"` // PRIMARY KEY
	Value    sql.NullString `db:"VALUE"`
}

type FeatureCvtermprop struct {
	FeatureCvtermpropID int            `db:"FEATURE_CVTERMPROP_ID"`
	Value               sql.NullString `db:"VALUE"`
}

type FeaturePubprop struct {
	FeaturePubpropID int            `db:"FEATURE_PUBPROP_ID"`
	Value            sql.NullString `db:"VALUE"`
}

type Organism struct {
	OrganismID int            `db:"ORGANISM_ID"`
	Comment    sql.NullString `db:"COMMENT_"`
}

type Organismprop struct {
	OrganismpropID int            `db:"ORGANISMPROP_ID"`
	Value          sql.NullString `db:"VALUE"`
}

type Pubprop struct {
	PubpropID int            `db:"PUBPROP_ID"`
	Value     sql.NullString `db:"VALUE"`
}

type Analysis struct {
	AnalysisID  int            `db:"ANALYSIS_ID"`
	Description sql.NullString `db:"DESCRIPTION"`
	Sourceuri   sql.NullString `db:"SOURCEURI"`
}

type ChadoLogs struct {
	ChadoLogsID int            `db:"CHADO_LOGS_ID"`
	NewValue    sql.NullString `db:"NEW_VALUE"`
	OldValue    sql.NullString `db:"OLD_VALUE"`
}

type Dbxref struct {
	DbxrefID    int            `db:"DBXREF_ID"`
	Description sql.NullString `db:"DESCRIPTION"`
}

type Featureloc struct {
	FeaturelocID int            `db:"FEATURELOC_ID"`
	ResidueInfo  sql.NullString `db:"RESIDUE_INFO"`
}

type Cv struct {
	CvID       int            `db:"CV_ID"`
	Definition sql.NullString `db:"DEFINITION"`
}

type Feature struct {
	FeatureID int            `db:"FEATURE_ID"`
	Residues  sql.NullString `db:"RESIDUES"`
}

type Pub struct {
	PubID       int            `db:"PUB_ID"`
	Title       sql.NullString `db:"TITLE"`
	Volumetitle sql.NullString `db:"VOLUMETITLE"`
}

type Arraydesign struct {
	ArraydesignID     int            `db:"ARRAYDESIGN_ID"`
	ArrayDimensions   sql.NullString `db:"ARRAY_DIMENSIONS"`
	Description       sql.NullString `db:"DESCRIPTION"`
	ElementDimensions sql.NullString `db:"ELEMENT_DIMENSIONS"`
	Version           sql.NullString `db:"VERSION"`
}

type Chadoprop struct {
	ChadopropID int            `db:"CHADOPROP_ID"`
	Value       sql.NullString `db:"VALUE"`
}

type Cvterm struct {
	CvtermID   int            `db:"CVTERM_ID"`
	Definition sql.NullString `db:"DEFINITION"`
}

type Featureprop struct {
	FeaturepropID int            `db:"FEATUREPROP_ID"`
	Value         sql.NullString `db:"VALUE"`
}

type FeatureRelationship struct {
	FeatureRelationshipID int            `db:"FEATURE_RELATIONSHIP_ID"`
	Value                 sql.NullString `db:"VALUE"`
}

type MvGeneDescription struct {
	MvGeneDescriptionID int            `db:"MV_GENE_DESCRIPTION_ID"`
	Description         sql.NullString `db:"DESCRIPTION"`
}

type Phenotype struct {
	PhenotypeID int            `db:"PHENOTYPE_ID"`
	Value       sql.NullString `db:"VALUE"`
}
