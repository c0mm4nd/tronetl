package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"sync"
	"testing"

	"github.com/jszwec/csvutil"
)

type sampleCSV struct {
	ID   int    `csv:"id"`
	Name string `csv:"name"`
}

// Ensures createCSVEncodeCh drains all objects and writes CSV rows.
func TestCreateCSVEncodeCh(t *testing.T) {
	buf := &bytes.Buffer{}
	csvWriter := csv.NewWriter(buf)
	enc := csvutil.NewEncoder(csvWriter)

	var wg sync.WaitGroup
	ch := createCSVEncodeCh(&wg, enc, csvWriter, 2)

	total := 5
	for i := 0; i < total; i++ {
		ch <- sampleCSV{ID: i, Name: fmt.Sprintf("n%d", i)}
	}
	close(ch)
	wg.Wait()
	csvWriter.Flush()

	rows, err := csv.NewReader(bytes.NewReader(buf.Bytes())).ReadAll()
	if err != nil {
		t.Fatalf("failed to read csv: %v", err)
	}

	if len(rows) != total+1 { // header + rows
		t.Fatalf("expected %d rows, got %d", total+1, len(rows))
	}

	if rows[1][0] != "0" || rows[1][1] != "n0" {
		t.Fatalf("unexpected first data row: %v", rows[1])
	}
}
