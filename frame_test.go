package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFrame(t *testing.T) {
	f := NewFrame(3)
	require.NotNil(t, f)
	assert.Equal(t, 3, f.cap)
	assert.Equal(t, true, f.empty)
	assert.Zero(t, f.begin)
	assert.Zero(t, f.end)
	assert.Equal(t, []float64{0, 0, 0}, f.buffer)
	assert.Equal(t, 0, f.Len())
	assert.Equal(t, []float64{}, f.Values())
}

func TestFrameIndex(t *testing.T) {
	f := NewFrame(3)
	require.NotNil(t, f)

	assert.Equal(t, 0, f.index(-3))
	assert.Equal(t, 1, f.index(-2))
	assert.Equal(t, 2, f.index(-1))
	assert.Equal(t, 0, f.index(0))
	assert.Equal(t, 1, f.index(1))
	assert.Equal(t, 2, f.index(2))
	assert.Equal(t, 0, f.index(3))
}

func TestFrameInsertForward(t *testing.T) {
	f := NewFrame(3)
	require.NotNil(t, f)

	f.InsertForward(1)
	assert.Equal(t, false, f.empty)
	assert.Equal(t, 0, f.begin)
	assert.Equal(t, 1, f.end)
	assert.Equal(t, []float64{1, 0, 0}, f.buffer)
	assert.Equal(t, 1, f.Len())
	assert.Equal(t, []float64{1}, f.Values())

	f.InsertForward(2)
	assert.Equal(t, false, f.empty)
	assert.Equal(t, 0, f.begin)
	assert.Equal(t, 2, f.end)
	assert.Equal(t, []float64{1, 2, 0}, f.buffer)
	assert.Equal(t, 2, f.Len())
	assert.Equal(t, []float64{1, 2}, f.Values())

	f.InsertForward(3)
	assert.Equal(t, false, f.empty)
	assert.Equal(t, 0, f.begin)
	assert.Equal(t, 0, f.end)
	assert.Equal(t, []float64{1, 2, 3}, f.buffer)
	assert.Equal(t, 3, f.Len())
	assert.Equal(t, []float64{1, 2, 3}, f.Values())

	f.InsertForward(4)
	assert.Equal(t, false, f.empty)
	assert.Equal(t, 1, f.begin)
	assert.Equal(t, 1, f.end)
	assert.Equal(t, []float64{4, 2, 3}, f.buffer)
	assert.Equal(t, 3, f.Len())
	assert.Equal(t, []float64{2, 3, 4}, f.Values())

}

func TestFrameInsertBackward(t *testing.T) {
	f := NewFrame(3)
	require.NotNil(t, f)

	f.InsertBackward(1)
	assert.Equal(t, false, f.empty)
	assert.Equal(t, 2, f.begin)
	assert.Equal(t, 0, f.end)
	assert.Equal(t, []float64{0, 0, 1}, f.buffer)
	assert.Equal(t, 1, f.Len())
	assert.Equal(t, []float64{1}, f.Values())

	f.InsertBackward(2)
	assert.Equal(t, false, f.empty)
	assert.Equal(t, 1, f.begin)
	assert.Equal(t, 0, f.end)
	assert.Equal(t, []float64{0, 2, 1}, f.buffer)
	assert.Equal(t, 2, f.Len())
	assert.Equal(t, []float64{2, 1}, f.Values())

	f.InsertBackward(3)
	assert.Equal(t, false, f.empty)
	assert.Equal(t, 0, f.begin)
	assert.Equal(t, 0, f.end)
	assert.Equal(t, []float64{3, 2, 1}, f.buffer)
	assert.Equal(t, 3, f.Len())
	assert.Equal(t, []float64{3, 2, 1}, f.Values())

	f.InsertBackward(4)
	assert.Equal(t, false, f.empty)
	assert.Equal(t, 2, f.begin)
	assert.Equal(t, 2, f.end)
	assert.Equal(t, []float64{3, 2, 4}, f.buffer)
	assert.Equal(t, 3, f.Len())
	assert.Equal(t, []float64{4, 3, 2}, f.Values())

}

func TestFrameInsert(t *testing.T) {
	f := NewFrame(3)
	require.NotNil(t, f)

	f.InsertForward(1)
	assert.Equal(t, 3, f.Cap())
	assert.Equal(t, 1, f.Len())
	assert.Equal(t, []float64{1}, f.Values())

	f.InsertBackward(2)
	assert.Equal(t, 3, f.Cap())
	assert.Equal(t, 2, f.Len())
	assert.Equal(t, []float64{2, 1}, f.Values())

	f.InsertForward(3)
	assert.Equal(t, 3, f.Cap())
	assert.Equal(t, 3, f.Len())
	assert.Equal(t, []float64{2, 1, 3}, f.Values())

	f.InsertBackward(4)
	assert.Equal(t, 3, f.Cap())
	assert.Equal(t, 3, f.Len())
	assert.Equal(t, []float64{4, 2, 1}, f.Values())

	f.InsertForward(5)
	assert.Equal(t, 3, f.Cap())
	assert.Equal(t, 3, f.Len())
	assert.Equal(t, []float64{2, 1, 5}, f.Values())

	f.InsertBackward(6)
	assert.Equal(t, 3, f.Cap())
	assert.Equal(t, 3, f.Len())
	assert.Equal(t, []float64{6, 2, 1}, f.Values())
}

func TestFrameResize(t *testing.T) {
	f1 := NewFrame(5)
	require.NotNil(t, f1)

	f1.InsertForward(1)
	f1.InsertForward(2)
	f1.InsertForward(3)
	f1.InsertForward(4)
	assert.Equal(t, []float64{1, 2, 3, 4}, f1.Values())

	f2 := f1.Resize(3)
	require.NotNil(t, f2)
	assert.Equal(t, []float64{1, 2, 3}, f2.Values())

	f3 := f1.Resize(7)
	require.NotNil(t, f3)
	assert.Equal(t, []float64{1, 2, 3, 4}, f3.Values())

}
