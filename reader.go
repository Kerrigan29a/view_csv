package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Reader struct {
	path           string
	hasColumnNames bool
	line           uint
	column         int
	frame          *Frame
}

var EOF = errors.New("EOF")
var BOF = errors.New("BOF")

func NewReader(path string, hasColumnNames bool) *Reader {
	r := &Reader{path: path, hasColumnNames: hasColumnNames}
	if hasColumnNames {
		r.line = 1
	}
	return r
}

func (r *Reader) SelectNamedColumn(name string) error {
	if !r.hasColumnNames {
		return fmt.Errorf("the flag hasColumnNames must be enabled")
	}

	/* Read row with names*/
	rows, err := r.readLines(0, 1)
	if err != nil {
		return fmt.Errorf("unable to read a row: %s", err.Error())
	}
	row := rows[0]

	/*Find name in row */
	for i, column := range row {
		if column == name {
			r.column = i
			return nil
		}
	}
	return fmt.Errorf("unable to find the column named '%s' in '%s'", name, r.path)
}

func (r *Reader) SelectColumn(column int) {
	r.column = column
}

func (r *Reader) FrameCap() int {
	if r.frame != nil {
		return r.frame.Cap()
	}
	return 0
}

func (r *Reader) SetFrameCap(n int) {
	if r.frame != nil {
		r.frame = r.frame.Resize(n)
	} else {
		r.frame = NewFrame(n)
	}
}

func (r *Reader) normalizeColumn(amount int) (int, error) {
	if r.column >= 0 {
		if r.column >= amount {
			return -1, fmt.Errorf("out of range")
		}
		return r.column, nil
	}
	result := amount + r.column
	if result < 0 {
		return -1, fmt.Errorf("out of range")
	}
	return result, nil
}

func (r *Reader) readLines(begin, end uint) ([][]string, error) {
	if begin >= end {
		return nil, fmt.Errorf("begin can't be grater or equal to end")
	}

	f, err := os.Open(r.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := csv.NewReader(f)

	rows := [][]string{}
	for i := uint(0); i < end; i++ {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return rows, nil
			}
			return nil, err
		}
		if i >= begin {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func moveNext(line, n uint) uint {
	return line + n
}

func (r *Reader) MoveNext(n uint) {
	r.line = moveNext(r.line, n)
}

func movePrev(line, n uint, hasColumnNames bool) (uint, int) {
	newLine := int(line) - int(n)
	minValue := 0
	if hasColumnNames {
		minValue = 1
	}
	if newLine < minValue {
		delta := minValue - newLine
		return uint(minValue), delta // Begin Of File REACHED
	}
	return uint(newLine), 0 // Begin Of File NOT REACHED
}

func (r *Reader) MovePrev(n uint) {
	r.line, _ = movePrev(r.line, n, r.hasColumnNames)
}

func (r *Reader) ReadNext(n uint) error {
	/* Check frame */
	if r.frame == nil {
		return fmt.Errorf("init frame")
	}

	/* Read rows */
	rows, err := r.readLines(r.line, r.line+n)
	if err != nil {
		return fmt.Errorf("unable to read a row: %s", err.Error())
	}
	amount := uint(len(rows))

	/* Update line */
	r.MoveNext(amount)

	/* Write N values */
	for i := uint(0); i < amount; i++ {
		row := rows[i]
		/* Write value*/
		column, err := r.normalizeColumn(len(row))
		if err != nil {
			return fmt.Errorf("bad column: %s", err.Error())
		}
		value := row[column]
		number, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("unable to convert value to float: %s", err.Error())
		}
		r.frame.InsertForward(number)
	}
	err = nil
	if amount != n {
		err = EOF
	}
	return err
}

func (r *Reader) ReadPrev(n uint) error {
	/* Check frame */
	if r.frame == nil {
		return fmt.Errorf("init frame")
	}
	if (r.hasColumnNames && r.line == 1) || !r.hasColumnNames && r.line == 0 {
		return BOF
	}

	/* Read rows */
	line, delta := movePrev(r.line, uint(r.frame.Len())+n, r.hasColumnNames)
	bof := delta > 0

	/* If the new claculated N is 0 or less, it means that the reader is at the beginning of the file and can't go back more */
	newN := n - uint(delta)
	if newN <= 0 {
		return BOF
	}

	rows, err := r.readLines(line, line+newN)
	if err != nil {
		return fmt.Errorf("unable to read a row: %s", err.Error())
	}
	amount := uint(len(rows))

	/* Update line if len >= cap */
	if r.frame.Len() >= r.frame.Cap() {
		r.MovePrev(amount)
	}

	/* Write N values */
	for i := int(amount - 1); i >= 0; i-- {
		row := rows[i]
		/* Write value*/
		column, err := r.normalizeColumn(len(row))
		if err != nil {
			return fmt.Errorf("bad column: %s", err.Error())
		}
		value := row[column]
		number, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("unable to convert value to float: %s", err.Error())
		}
		r.frame.InsertBackward(number)
	}
	err = nil
	if bof || amount != n {
		err = BOF
	}
	return err
}

func (r *Reader) Values() ([]float64, error) {
	/* Check frame */
	if r.frame == nil {
		return nil, fmt.Errorf("init frame")
	}

	return r.frame.Values(), nil
}
