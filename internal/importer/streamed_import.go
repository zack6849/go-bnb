// Streamed import utilities, used to import large datasets without loading it all into memory at once
// Allows you to read a file arbitrarily with a function to determine what the next record is, as well as a parse method per-record
// Current use intended for reading large CSV files, but wanted to make this somewhat generic
// Supports reading CSVs with columns and a way to relate the header to the individual record
//  or just ignoring the header and reading in arbitrary data at given indexes

package importer

import (
	"encoding/csv"
	"fmt"
	"io"
)

// InsertFunc generic insert function that receives an array of T and returns an error if one is thrown
type InsertFunc[T any] func(batch []T) error

// HeaderAwareRecordFunc header-aware record processing function
type HeaderAwareRecordFunc[T any] func(header []string, record []string) (T, error)

// NextFunc generic next func that returns either the next T to process, or an error.
type NextFunc[T any] func() (T, error)

func StreamHeaderAwareImport[T any](
	r io.Reader,
	parseRecord HeaderAwareRecordFunc[T],
	batchSize int,
	insertBatch InsertFunc[T],
) error {
	//open a new csv reader with the file handle
	reader := csv.NewReader(r)
	//read the header raw array from the record, eg: array of columns: [id, name, type, foo, bar, baz]
	headerData, err := reader.Read()
	//if we encounter some kind of CSV error when trying to parse the record, bubble the exception up
	if err != nil {
		return fmt.Errorf("failed to read header record: %w", err)
	}
	//define the next function to grab the next item from the CSV
	next := func() (T, error) {
		//read a record in
		record, err := reader.Read()
		//if we fail to read a record in, fail and return the zero type for T
		if err != nil {
			var zero T
			return zero, err
		}
		//call the parse record callback with the header and the record so they can be combined
		return parseRecord(headerData, record)
	}
	//wraps around the stream import with a wrapper next function basically
	return StreamImport(
		next,
		batchSize,
		insertBatch,
	)
}

func StreamImport[T any](
	next NextFunc[T],
	batchSize int,
	insertBatch InsertFunc[T],
) error {
	if batchSize < 1 {
		return fmt.Errorf("invalid batch size %d provided", batchSize)
	}
	for {
		batch, more, err := readBatch(next, batchSize)
		if err != nil {
			return err
		}
		if len(batch) > 0 {
			if err := insertBatch(batch); err != nil {
				return err
			}
		}
		if !more {
			break
		}
	}
	return nil
}

func readBatch[T any](
	next NextFunc[T],
	batchSize int,
) (batch []T, more bool, err error) {
	batch = make([]T, 0, batchSize)
	for range batchSize {
		item, err := next()
		//if we hit eof
		if err == io.EOF {
			//return what we have immediately, false for more, and
			return batch, false, nil
		}
		//if it's a real error, return nothing, no more, and the error to bubble up
		if err != nil {
			return nil, false, err
		}
		//keep pushing to the stack
		batch = append(batch, item)
	}
	//we filled our buffer, but there should be more, tell the caller to re-try
	return batch, true, nil
}

func ArrayCombine(keys []string, values []string) (map[string]string, error) {
	data := make(map[string]string)
	if len(keys) != len(values) {
		return nil, fmt.Errorf("row does not have the same number of values as the header row")
	}
	for i := range len(keys) {
		header := keys[i]
		value := values[i]
		data[header] = value
	}
	return data, nil
}
