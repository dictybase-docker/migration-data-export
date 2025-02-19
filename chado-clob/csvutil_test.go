package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type taggedStruct struct {
	Field1 int    `db:"custom_name"`
	Field2 string // No tag
}

type untaggedStruct struct {
	ID   int
	Name string
}

type validStruct struct {
	ID   int    `db:"test_id"`
	Name string `db:"test_name"`
	Age  int64
}

type invalidStruct struct {
	Value float64
}

func TestCSVHandler_ToRow(t *testing.T) {
	t.Run("valid struct with all supported types", func(t *testing.T) {
		h := &CSVHandler{}
		record := &validStruct{
			ID:   123,
			Name: "test",
			Age:  30,
		}

		got, err := h.ToRow(record)
		assert.NoError(t, err)
		assert.Equal(t, []string{"123", "test", "30"}, got)
	})

	t.Run("unsupported field type", func(t *testing.T) {
		h := &CSVHandler{}
		record := &invalidStruct{Value: 3.14}

		_, err := h.ToRow(record)
		assert.Error(t, err)
	})
}

func TestCSVHandler_Header(t *testing.T) {
	t.Run("struct with db tags", func(t *testing.T) {
		h := &CSVHandler{}
		got := h.Header(&taggedStruct{})
		assert.Equal(t, []string{"custom_name", "Field2"}, got)
	})

	t.Run("struct without tags", func(t *testing.T) {
		h := &CSVHandler{}
		got := h.Header(&untaggedStruct{})
		assert.Equal(t, []string{"ID", "Name"}, got)
	})
}
