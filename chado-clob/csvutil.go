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

// NewCSVWriter creates a new CSVWriter instance
func NewCSVWriter(writer *csv.Writer, file *os.File) *CSVWriter {
	return &CSVWriter{
		Writer: writer,
		file:   file,
	}
}

func (csvw *CSVWriter) Close() error {
	csvw.Flush()
	if err := csvw.Error(); err != nil {
		return err
	}
	return csvw.file.Close()
}

type CSVHandler struct{}

// NewCSVHandler creates a new CSVHandler instance
func NewCSVHandler() *CSVHandler {
	return &CSVHandler{}
}

func (h *CSVHandler) ToRow(record interface{}) ([]string, error) {
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

func (h *CSVHandler) Header(record interface{}) []string {
	val := reflect.ValueOf(record).Elem()
	typ := val.Type()
	var headers []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if dbTag, ok := field.Tag.Lookup("db"); ok {
			headers = append(headers, dbTag)
		} else {
			headers = append(headers, field.Name)
		}
	}
	return headers
}

func (h *CSVHandler) CreateWriter(
	outputFile string,
	record interface{},
) (*CSVWriter, error) {
	if err := validateRecord(record); err != nil {
		return nil, err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %w", err)
	}

	csvWriter := csv.NewWriter(file)
	if err := csvWriter.Write(h.Header(record)); err != nil {
		file.Close()
		return nil, fmt.Errorf("error writing headers: %w", err)
	}

	return NewCSVWriter(csvWriter, file), nil
}

func validateRecord(record interface{}) error {
	val := reflect.ValueOf(record)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("record must be a pointer to struct")
	}
	return nil
}
