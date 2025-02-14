package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
)

type CSVWriter struct {
	*csv.Writer
	file *os.File
}

func (csvw *CSVWriter) Close() error {
	csvw.Flush()
	if err := csvw.Error(); err != nil {
		return err
	}
	return csvw.file.Close()
}

func structToCSVRow(record interface{}) ([]string, error) {
	val := reflect.ValueOf(record).Elem()
	var row []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		switch field.Kind() {
		case reflect.Int, reflect.Int64:
			row = append(row, fmt.Sprintf("%d", field.Int()))
		case reflect.String:
			row = append(row, field.String())
		default:
			return nil, fmt.Errorf("unsupported field type: %s", field.Kind())
		}
	}
	return row, nil
}

func getCSVHeader(record interface{}) []string {
	val := reflect.ValueOf(record).Elem()
	typ := val.Type()
	var headers []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if dbTag := field.Tag.Get("db"); dbTag != "" {
			headers = append(headers, dbTag)
		} else {
			headers = append(headers, field.Name)
		}
	}
	return headers
}

func createCSVWriter(
	outputFile string,
	record interface{},
) (*CSVWriter, error) {
	file, err := os.Create(outputFile)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %w", err)
	}

	csvWriter := csv.NewWriter(file)
	if err := csvWriter.Write(getCSVHeader(record)); err != nil {
		file.Close()
		return nil, fmt.Errorf("error writing headers: %w", err)
	}

	return &CSVWriter{
		Writer: csvWriter,
		file:   file,
	}, nil
}
