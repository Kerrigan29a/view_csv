package main

type Frame struct {
	buffer []float64
	cap    int
	begin  int
	end    int
	empty  bool
}

func NewFrame(n int) *Frame {
	return &Frame{
		buffer: make([]float64, n, n),
		cap:    n,
		begin:  0,
		end:    0,
		empty:  true,
	}
}

func (f *Frame) Reset() {
	f.begin = 0
	f.end = 0
	f.empty = true
}

func (f *Frame) Resize(n int) *Frame {
	if n == f.cap {
		return f
	}
	newF := NewFrame(n)
	oldValues := f.Values()
	nOldValues := len(oldValues)
	minCap := nOldValues
	if n < nOldValues {
		minCap = n
	}
	for i := 0; i < minCap; i++ {
		newF.InsertForward(oldValues[i])
	}
	return newF
}

func (f *Frame) index(n int) int {
	result := n % f.cap
	if result >= 0 {
		return result
	}
	return result + f.cap
}

func (f *Frame) InsertForward(value float64) {
	f.buffer[f.end] = value
	/* Before update end, move forward begin when buffer is full */
	if f.Len() >= f.cap {
		f.begin = f.index(f.begin + 1)
	}
	f.end = f.index(f.end + 1)
	f.empty = false
}

func (f *Frame) InsertBackward(value float64) {
	begin := f.index(f.begin - 1)
	f.buffer[begin] = value
	/* Before update begi, move backward begin when buffer is full */
	if f.Len() >= f.cap {
		f.end = f.index(f.end - 1)
	}
	f.begin = begin
	f.empty = false
}

func (f *Frame) Len() int {
	if f.begin == f.end {
		if f.empty {
			return 0
		}
		return f.cap
	}
	if f.begin < f.end {
		return f.end - f.begin
	}
	return f.cap - f.begin + f.end
}

func (f *Frame) Cap() int {
	return f.cap
}

func (f *Frame) Values() []float64 {
	if f.empty {
		return []float64{}
	}
	if f.begin < f.end {
		return append([]float64{}, f.buffer[f.begin:f.end]...)
	}
	if f.begin == f.end && f.begin == 0 {
		return append([]float64{}, f.buffer...)
	}
	values := append([]float64{}, f.buffer[f.begin:]...)
	return append(values, f.buffer[:f.end]...)
}
