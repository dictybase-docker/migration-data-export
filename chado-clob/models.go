package main

import "strings"

// GetStructForTable returns the appropriate struct type for a given table name
func GetStructForTable(tableName string) interface{} {
    switch strings.ToLower(tableName) {
    case "cvprop":
        return &Cvprop{}
    case "featurecvtermprop":
        return &FeatureCvtermprop{}
    case "featurepubprop":
        return &FeaturePubprop{}
    case "organism":
        return &Organism{}
    case "pubprop":
        return &Pubprop{}
    case "analysis":
        return &Analysis{}
    case "chadologs":
        return &ChadoLogs{}
    case "dbxref":
        return &Dbxref{}
    case "featureloc":
        return &Featureloc{}
    case "cv":
        return &Cv{}
    case "feature":
        return &Feature{}
    case "pub":
        return &Pub{}
    case "arraydesign":
        return &Arraydesign{}
    case "chadoprop":
        return &Chadoprop{}
    case "cvterm":
        return &Cvterm{}
    case "featureprop":
        return &Featureprop{}
    case "featurerelationship":
        return &FeatureRelationship{}
    case "mvgenedescription":
        return &MvGeneDescription{}
    case "phenotype":
        return &Phenotype{}
    default:
        return nil
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
    Comment    string `gorm:"column:COMMENT_"`
}

type Pubprop struct {
    PubpropID int
    Value     string
}

type Analysis struct {
    AnalysisID   int
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
    Version          string
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
