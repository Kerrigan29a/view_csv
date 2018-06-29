package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createCSV(t *testing.T, txt string) string {
	tmpfile, err := ioutil.TempFile("", "dummy_csv")
	if err != nil {
		require.FailNowf(t, "unable to create CSV file: %s", err.Error())
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write([]byte(txt)); err != nil {
		require.FailNowf(t, "unable to write CSV file: %s", err.Error())
	}

	require.FileExists(t, tmpfile.Name())
	return tmpfile.Name()
}

func TestNewReader(t *testing.T) {
	path := createCSV(t, "hi world")
	defer os.Remove(path)

	r1 := NewReader(path, false)
	require.NotNil(t, r1)

	assert.Equal(t, path, r1.path)
	assert.False(t, r1.hasColumnNames)
	assert.EqualValues(t, 0, r1.line)
	assert.Zero(t, r1.column)
	assert.Nil(t, r1.frame)
	assert.Equal(t, 0, r1.FrameCap())

	r2 := NewReader(path, true)
	require.NotNil(t, r2)

	assert.Equal(t, path, r2.path)
	assert.True(t, r2.hasColumnNames)
	assert.EqualValues(t, 1, r2.line)
	assert.Zero(t, r2.column)
	assert.Nil(t, r2.frame)
	assert.Equal(t, 0, r2.FrameCap())
}

func TestReaderNormalizeColumn(t *testing.T) {
	path := createCSV(t, "hi world")
	defer os.Remove(path)

	r := NewReader(path, false)
	require.NotNil(t, r)

	r.column = -4
	col, err := r.normalizeColumn(3)
	assert.Error(t, err)
	assert.Equal(t, -1, col)

	r.column = -3
	col, err = r.normalizeColumn(3)
	assert.NoError(t, err)
	assert.Equal(t, 0, col)

	r.column = -2
	col, err = r.normalizeColumn(3)
	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.Equal(t, 1, col)

	r.column = -1
	col, err = r.normalizeColumn(3)
	assert.NoError(t, err)
	assert.Equal(t, 2, col)

	r.column = 0
	col, err = r.normalizeColumn(3)
	assert.NoError(t, err)
	assert.Equal(t, 0, col)

	r.column = 1
	col, err = r.normalizeColumn(3)
	assert.NoError(t, err)
	assert.Equal(t, 1, col)

	r.column = 2
	col, err = r.normalizeColumn(3)
	assert.NoError(t, err)
	assert.Equal(t, 2, col)

	r.column = 3
	col, err = r.normalizeColumn(3)
	assert.Error(t, err)
	assert.Equal(t, -1, col)
}

func TestReaderReadLines(t *testing.T) {
	path := createCSV(t, `
value1,value2
1,10
2,20
`)
	defer os.Remove(path)

	r := NewReader(path, true)
	require.NotNil(t, r)

	lines, err := r.readLines(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, [][]string{{"value1", "value2"}}, lines)

	lines, err = r.readLines(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, [][]string{{"1", "10"}}, lines)

	lines, err = r.readLines(2, 3)
	assert.NoError(t, err)
	assert.Equal(t, [][]string{{"2", "20"}}, lines)

	lines, err = r.readLines(0, 2)
	assert.NoError(t, err)
	assert.Equal(t, [][]string{{"value1", "value2"}, {"1", "10"}}, lines)

	lines, err = r.readLines(1, 3)
	assert.NoError(t, err)
	assert.Equal(t, [][]string{{"1", "10"}, {"2", "20"}}, lines)

	lines, err = r.readLines(0, 3)
	assert.NoError(t, err)
	assert.Equal(t, [][]string{{"value1", "value2"}, {"1", "10"}, {"2", "20"}}, lines)

	lines, err = r.readLines(0, 4)
	assert.NoError(t, err)
	assert.Equal(t, [][]string{{"value1", "value2"}, {"1", "10"}, {"2", "20"}}, lines)

	lines, err = r.readLines(0, 0)
	assert.Error(t, err)
	assert.Nil(t, lines)

	lines, err = r.readLines(1, 0)
	assert.Error(t, err)
	assert.Nil(t, lines)
}

func TestReaderSelectColumn(t *testing.T) {
	path := createCSV(t, "hi world")
	defer os.Remove(path)

	r := NewReader(path, false)
	require.NotNil(t, r)

	r.SelectColumn(1)
	assert.Equal(t, 1, r.column)
}

func TestReaderSelectNamedColumn(t *testing.T) {
	path := createCSV(t, `
value1,value2
1,10
2,20
3,30
4,40
5,50
6,60
7,70
8,80
9,90
10,100`)
	defer os.Remove(path)

	r := NewReader(path, true)
	require.NotNil(t, r)

	err := r.SelectNamedColumn("value1")
	assert.NoError(t, err)
	assert.Equal(t, 0, r.column)

	err = r.SelectNamedColumn("value2")
	assert.NoError(t, err)
	assert.Equal(t, 1, r.column)

	err = r.SelectNamedColumn("value3")
	assert.Error(t, err)
}

func TestNewMove(t *testing.T) {
	path := createCSV(t, `
value1,value2
1,10
2,20`)
	defer os.Remove(path)

	r1 := NewReader(path, false)
	require.NotNil(t, r1)

	r1.MoveNext(5)
	assert.EqualValues(t, 5, r1.line)
	r1.MovePrev(2)
	assert.EqualValues(t, 3, r1.line)
	r1.MovePrev(4)
	assert.EqualValues(t, 0, r1.line)
	r1.MovePrev(1)
	assert.EqualValues(t, 0, r1.line)

	r2 := NewReader(path, true)
	require.NotNil(t, r2)
	r2.MoveNext(5)
	assert.EqualValues(t, 6, r2.line)
	r2.MovePrev(2)
	assert.EqualValues(t, 4, r2.line)
	r2.MovePrev(4)
	assert.EqualValues(t, 1, r2.line)
	r2.MovePrev(1)
	assert.EqualValues(t, 1, r2.line)
}

func TestReaderReadNext(t *testing.T) {
	path := createCSV(t, `
1,10
2,20
3,30
4,40
5,50
6,60
7,70
8,80
9,90
10,100`)
	defer os.Remove(path)

	r := NewReader(path, false)
	require.NotNil(t, r)

	assert.EqualError(t, r.ReadNext(1), "init frame")
	assert.EqualValues(t, 0, r.line)

	r.SelectColumn(0)
	r.SetFrameCap(3)
	assert.Equal(t, 3, r.FrameCap())

	assert.NoError(t, r.ReadNext(1))
	values, err := r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{1}, values)
	assert.EqualValues(t, 1, r.line)

	assert.NoError(t, r.ReadNext(2))
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{1, 2, 3}, values)
	assert.EqualValues(t, 3, r.line)

	assert.NoError(t, r.ReadNext(1))
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{2, 3, 4}, values)
	assert.EqualValues(t, 4, r.line)

	assert.NoError(t, r.ReadNext(2))
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{4, 5, 6}, values)
	assert.EqualValues(t, 6, r.line)

	assert.NoError(t, r.ReadNext(3))
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{7, 8, 9}, values)
	assert.EqualValues(t, 9, r.line)

	assert.EqualError(t, r.ReadNext(4), "EOF")
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{8, 9, 10}, values)
	assert.EqualValues(t, 10, r.line)

	assert.EqualError(t, r.ReadNext(1), "EOF")
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{8, 9, 10}, values)
	assert.EqualValues(t, 10, r.line)
}

func TestReaderReadPrev(t *testing.T) {
	path := createCSV(t, `
value1,value2
1,10
2,20
3,30
4,40
5,50
6,60
7,70
8,80
9,90
10,100`)
	defer os.Remove(path)

	r := NewReader(path, true)
	require.NotNil(t, r)

	assert.EqualError(t, r.ReadPrev(1), "init frame")

	r.SelectColumn(1)
	r.SetFrameCap(3)
	assert.Equal(t, 3, r.FrameCap())

	assert.EqualValues(t, 1, r.line)
	assert.EqualError(t, r.ReadPrev(1), "BOF")
	values, err := r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{}, values)
	assert.EqualValues(t, 1, r.line)

	r.MoveNext(10)
	assert.EqualValues(t, 11, r.line)

	assert.NoError(t, r.ReadPrev(1))
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{100}, values)
	assert.EqualValues(t, 11, r.line)

	assert.NoError(t, r.ReadPrev(2))
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{80, 90, 100}, values)
	assert.EqualValues(t, 11, r.line)

	assert.NoError(t, r.ReadPrev(1))
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{70, 80, 90}, values)
	assert.EqualValues(t, 10, r.line)

	assert.NoError(t, r.ReadPrev(2))
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{50, 60, 70}, values)
	assert.EqualValues(t, 8, r.line)

	assert.NoError(t, r.ReadPrev(3))
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{20, 30, 40}, values)
	assert.EqualValues(t, 5, r.line)

	assert.EqualError(t, r.ReadPrev(4), "BOF")
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{10, 20, 30}, values)
	assert.EqualValues(t, 4, r.line)

	assert.EqualError(t, r.ReadPrev(4), "BOF")
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{10, 20, 30}, values)
	assert.EqualValues(t, 4, r.line)
}

func TestReaderRead(t *testing.T) {
	path := createCSV(t, `
value1,value2
1,10
2,20
3,30
4,40
5,50`)
	defer os.Remove(path)

	r := NewReader(path, true)
	require.NotNil(t, r)

	assert.EqualError(t, r.ReadPrev(1), "init frame")
	assert.EqualError(t, r.ReadNext(1), "init frame")

	r.SelectColumn(1)
	r.SetFrameCap(3)
	assert.Equal(t, 3, r.FrameCap())

	assert.EqualValues(t, 1, r.line)
	assert.EqualError(t, r.ReadPrev(1), "BOF")
	values, err := r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{}, values)
	assert.EqualValues(t, 1, r.line)

	assert.NoError(t, r.ReadNext(1))
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{10}, values)
	assert.EqualValues(t, 2, r.line)

	assert.EqualError(t, r.ReadPrev(2), "BOF")
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{10}, values)
	assert.EqualValues(t, 2, r.line)

	assert.NoError(t, r.ReadNext(2))
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{10, 20, 30}, values)
	assert.EqualValues(t, 4, r.line)

	assert.EqualError(t, r.ReadPrev(3), "BOF")
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{10, 20, 30}, values)
	assert.EqualValues(t, 4, r.line)

	assert.EqualError(t, r.ReadNext(3), "EOF")
	values, err = r.Values()
	assert.NoError(t, err)
	assert.Equal(t, []float64{30, 40, 50}, values)
	assert.EqualValues(t, 6, r.line)
}
