package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
)

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
) (*csv.Writer, error) {
	fh, err := os.Create(outputFile)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %w", err)
	}

	writer := csv.NewWriter(fh)
	if err := writer.Write(getCSVHeader(record)); err != nil {
		fh.Close()
		return nil, fmt.Errorf("error writing headers: %w", err)
	}

	return writer, nil
}
