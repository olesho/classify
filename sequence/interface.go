package sequence

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
)

type Stats struct {
	GroupsCount int `json:"groups_count"`
	ValuesCount int `json:"values_count"`
}

type FieldDetails struct {
	Type  string
	Name  string
	Attrs map[string]string
}

type JSONResp struct {
	Values     [][]string `json:"values"`
	Stats      Stats      `json:"stats"`
	FieldTypes []string   `json:"types"`
}

func (s *Series) GetFieldTypes() []string {
	fieldTypes := make([]string, 0)
	for _, t := range s.Group.Clusters {
		for _, fs := range t.FieldSets {
			fieldTypes =  append(fieldTypes, RecognizeType(fs.Content, fs.Type))
			//fieldTypes = append(fieldTypes, fs.Type)
		}
	}
	return fieldTypes
}

func (s *Series) ToJSON(w io.Writer) error {
	fieldTypes := s.GetFieldTypes()

	return json.NewEncoder(w).Encode(JSONResp{
		Values: s.TransposedValues,
		Stats: Stats{
			GroupsCount: len(s.TransposedValues),
			ValuesCount: len(s.TransposedValues[0]),
		},
		FieldTypes: fieldTypes,
	})
}

func (s *Series) ToText(w io.Writer) {
	for _, fields := range s.TransposedValues {
		for _, field := range fields {
			w.Write([]byte(fmt.Sprintln(field)))
		}
		w.Write([]byte(fmt.Sprintln("--------------------------")))
	}
}

func (s *Series) ToTextReversed(w io.Writer) {
	for i := len(s.TransposedValues) - 1; i > -1; i-- {
		fields := s.TransposedValues[i]
		for _, field := range fields {
			w.Write([]byte(fmt.Sprintln(field)))
		}
		w.Write([]byte(fmt.Sprintln("--------------------------")))
	}
}

func (s *Series) ToCSV(w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	for _, fields := range s.TransposedValues {
		if err := csvWriter.Write(fields); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return nil
}

func (s *Series) ToXPath(w io.Writer) {
	for _, path := range s.XPaths() {
		w.Write([]byte(fmt.Sprintln(path)))
	}
}
