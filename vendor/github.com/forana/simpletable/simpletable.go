package simpletable

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
)

// SimpleTable represents the data to be printed, in tabular form
type SimpleTable struct {
	headers      []string
	rows         [][]string
	columnWidths []int
}

var errInputNotSlice = errors.New("The input was not a slice.")
var errTypeNotConvertible = errors.New("SimpleTable can only print structs or slices.")

// New builds a new SimpleTable from your dataset, converting to string slices on the fly if the passed objects are not string slices.
// Will panic if `data` is not a slice.
func New(headers []string, data interface{}) (*SimpleTable, error) {
	table := SimpleTable{
		headers:      headers,
		rows:         [][]string{},
		columnWidths: nil,
	}

	dv := reflect.ValueOf(data)
	for n := 0; n < dv.Len(); n++ {
		item := dv.Index(n)
		row := []string{}
		t := item.Type()
		if t.Kind() == reflect.Slice {
			row = convertSlice(item)
		} else if t.Kind() == reflect.Struct {
			row = convertStruct(item)
		} else {
			return nil, errTypeNotConvertible
		}

		if table.columnWidths == nil {
			table.columnWidths = make([]int, len(row))
			for i, column := range headers {
				if len(column) > table.columnWidths[i] {
					table.columnWidths[i] = len(column)
				}
			}
		}
		for i, column := range row {
			if len(column) > table.columnWidths[i] {
				table.columnWidths[i] = len(column)
			}
		}
		table.rows = append(table.rows, row)
	}

	return &table, nil
}

func convertSlice(row reflect.Value) []string {
	result := []string{}
	for i := 0; i < row.Len(); i++ {
		result = append(result, fmt.Sprintf("%v", row.Index(i).Interface()))
	}
	return result
}

func convertStruct(row reflect.Value) []string {
	result := []string{}
	for i := 0; i < row.NumField(); i++ {
		result = append(result, fmt.Sprintf("%v", row.Field(i).Interface()))
	}
	return result
}

// Write prints the table, including headers, to the specified Writer.
func (table *SimpleTable) Write(w io.Writer) (err error) {
	_, err = w.Write([]byte(formatRow(table.headers, table.columnWidths)))
	if err != nil {
		return
	}
	for _, row := range table.rows {
		_, err = w.Write([]byte(formatRow(row, table.columnWidths)))
		if err != nil {
			return
		}
	}
	return
}

// Print prints the table, including headers, to stdout.
func (table *SimpleTable) Print() error {
	return table.Write(os.Stdout)
}

func formatRow(row []string, widths []int) string {
	columns := make([]string, len(row))
	for i, column := range row {
		columns[i] = pad(column, widths[i])
	}
	return strings.Join(columns, " ") + "\n"
}

func pad(s string, l int) string {
	return s + strings.Repeat(" ", l-len(s))
}

type sortingSlice struct {
	column int
	slice  *[][]string
}

func (ss sortingSlice) Len() int {
	return len(*ss.slice)
}

func (ss sortingSlice) Less(i, j int) bool {
	return strings.ToLower((*ss.slice)[i][ss.column]) < strings.ToLower((*ss.slice)[j][ss.column])
}

func (ss sortingSlice) Swap(i, j int) {
	(*ss.slice)[i], (*ss.slice)[j] = (*ss.slice)[j], (*ss.slice)[i]
}

// Sort sorts the table in place by the specified column index, ascending, case-insensitive.
func (table *SimpleTable) Sort(columnIndex int) {
	sort.Sort(sortingSlice{
		column: columnIndex,
		slice:  &table.rows,
	})
}

var errNotStruct = errors.New("The passed value was not a struct.")

// HeadersForType iterates through struct fields, returning a usable header slice for table printing.
func HeadersForType(obj interface{}) []string {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Struct {
		panic(errNotStruct)
	}
	headers := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		headers[i] = strings.ToUpper(t.Field(i).Name)
	}
	return headers
}
