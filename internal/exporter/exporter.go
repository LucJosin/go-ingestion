package exporter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/LucJosin/go-ingestion/internal/model"
)

var companyHeaders = []string{
	"name",
	"city",
	"country",
	"founded",
	"rank",
	"profile",
}

type Exporter struct {
	outputDir string
	filename  string
}

func NewExporter(output string, filename string) *Exporter {
	return &Exporter{
		outputDir: output,
		filename:  filename,
	}
}

func (e *Exporter) ExportData(data []model.Bank) error {
	switch e.outputDir {
	case "csv":
		return e.exportToCSV(data)
	case "json":
		return e.exportToJSON(data)
	case "stdout":
		return e.exportToStdOut(data)
	default:
		return fmt.Errorf("unknown output format: %s", e.outputDir)
	}
}

func (e *Exporter) exportToStdOut(data []model.Bank) error {
	for _, company := range data {

		jsonData, err := json.Marshal(company)
		if err != nil {
			return fmt.Errorf("Error marshalling JSON: %v\n", err)
		}

		fmt.Println(string(jsonData))
	}

	return nil
}

func (e *Exporter) exportToCSV(data []model.Bank) error {
	file, err := os.Create(e.filename + ".csv")
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(companyHeaders)
	if err != nil {
		return err
	}

	for _, company := range data {
		err := writer.Write([]string{
			company.Name,
			company.City,
			company.Country,
			strconv.Itoa(company.Founded),
			strconv.Itoa(company.Rank),
			company.Profile,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Exporter) exportToJSON(data []model.Bank) error {
	file, err := os.Create(e.filename + ".json")
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	var dtos []model.Bank
	for _, company := range data {
		dtos = append(dtos, company)
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(dtos); err != nil {
		return err
	}

	return nil
}
