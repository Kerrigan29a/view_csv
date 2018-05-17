package drawille

import (
	"fmt"
	"math"
)

// Braille chars start at 0x2800
var brailleStartOrdinal = 0x2800

func internalPosition(n, base int) int {
	if n >= 0 {
		return n % base
	}
	result := n % base
	if result == 0 {
		return -base
	}
	return result
}

func getDot(y, x int, inverse bool) int {
	y = internalPosition(y, 4)
	x = internalPosition(x, 2)

	/* x>=0 && y>=0 && !inverse */
	if y == 0 && x == 0 && !inverse {
		return 0x1
	}
	if y == 1 && x == 0 && !inverse {
		return 0x2
	}
	if y == 2 && x == 0 && !inverse {
		return 0x4
	}
	if y == 3 && x == 0 && !inverse {
		return 0x40
	}
	if y == 0 && x == 1 && !inverse {
		return 0x8
	}
	if y == 1 && x == 1 && !inverse {
		return 0x10
	}
	if y == 2 && x == 1 && !inverse {
		return 0x20
	}
	if y == 3 && x == 1 && !inverse {
		return 0x80
	}

	/* x>=0 && y>=0 && inverse */
	if y == 0 && x == 0 && inverse {
		return 0x40
	}
	if y == 1 && x == 0 && inverse {
		return 0x4
	}
	if y == 2 && x == 0 && inverse {
		return 0x2
	}
	if y == 3 && x == 0 && inverse {
		return 0x1
	}
	if y == 0 && x == 1 && inverse {
		return 0x80
	}
	if y == 1 && x == 1 && inverse {
		return 0x20
	}
	if y == 2 && x == 1 && inverse {
		return 0x10
	}
	if y == 3 && x == 1 && inverse {
		return 0x8
	}

	/* x<0 && y<0 && !inverse */
	if y == -1 && x == -1 && !inverse {
		return 0x80
	}
	if y == -2 && x == -1 && !inverse {
		return 0x20
	}
	if y == -3 && x == -1 && !inverse {
		return 0x10
	}
	if y == -4 && x == -1 && !inverse {
		return 0x8
	}
	if y == -1 && x == -2 && !inverse {
		return 0x40
	}
	if y == -2 && x == -2 && !inverse {
		return 0x4
	}
	if y == -3 && x == -2 && !inverse {
		return 0x2
	}
	if y == -4 && x == -2 && !inverse {
		return 0x1
	}

	/* x<0 && y<0 && inverse */
	if y == -1 && x == -1 && inverse {
		return 0x8
	}
	if y == -2 && x == -1 && inverse {
		return 0x10
	}
	if y == -3 && x == -1 && inverse {
		return 0x20
	}
	if y == -4 && x == -1 && inverse {
		return 0x80
	}
	if y == -1 && x == -2 && inverse {
		return 0x1
	}
	if y == -2 && x == -2 && inverse {
		return 0x2
	}
	if y == -3 && x == -2 && inverse {
		return 0x4
	}
	if y == -4 && x == -2 && inverse {
		return 0x40
	}

	/* x>=0 && y<0 && !inverse */
	if y == -1 && x == 0 && !inverse {
		return 0x40
	}
	if y == -2 && x == 0 && !inverse {
		return 0x4
	}
	if y == -3 && x == 0 && !inverse {
		return 0x2
	}
	if y == -4 && x == 0 && !inverse {
		return 0x1
	}
	if y == -1 && x == 1 && !inverse {
		return 0x80
	}
	if y == -2 && x == 1 && !inverse {
		return 0x20
	}
	if y == -3 && x == 1 && !inverse {
		return 0x10
	}
	if y == -4 && x == 1 && !inverse {
		return 0x8
	}

	/* x>=0 && y<0 && inverse */
	if y == -1 && x == 0 && inverse {
		return 0x1
	}
	if y == -2 && x == 0 && inverse {
		return 0x2
	}
	if y == -3 && x == 0 && inverse {
		return 0x4
	}
	if y == -4 && x == 0 && inverse {
		return 0x40
	}
	if y == -1 && x == 1 && inverse {
		return 0x8
	}
	if y == -2 && x == 1 && inverse {
		return 0x10
	}
	if y == -3 && x == 1 && inverse {
		return 0x20
	}
	if y == -4 && x == 1 && inverse {
		return 0x80
	}

	/* x<0 && y>=0 && !inverse */
	if y == 0 && x == -1 && !inverse {
		return 0x8
	}
	if y == 1 && x == -1 && !inverse {
		return 0x10
	}
	if y == 2 && x == -1 && !inverse {
		return 0x20
	}
	if y == 3 && x == -1 && !inverse {
		return 0x80
	}
	if y == 0 && x == -2 && !inverse {
		return 0x1
	}
	if y == 1 && x == -2 && !inverse {
		return 0x2
	}
	if y == 2 && x == -2 && !inverse {
		return 0x4
	}
	if y == 3 && x == -2 && !inverse {
		return 0x40
	}

	/* x<0 && y>=0 && inverse */
	if y == 0 && x == -1 && inverse {
		return 0x80
	}
	if y == 1 && x == -1 && inverse {
		return 0x20
	}
	if y == 2 && x == -1 && inverse {
		return 0x10
	}
	if y == 3 && x == -1 && inverse {
		return 0x8
	}
	if y == 0 && x == -2 && inverse {
		return 0x40
	}
	if y == 1 && x == -2 && inverse {
		return 0x4
	}
	if y == 2 && x == -2 && inverse {
		return 0x2
	}
	if y == 3 && x == -2 && inverse {
		return 0x1
	}

	panic(fmt.Sprintf("Unknown values: y=%d x=%d inverse=%t", y, x, inverse))
}

func externalPosition(n, base int) int {
	return int(math.Floor(float64(n) / float64(base)))
}

// Convert x,y to cols, rows
func getPos(x, y int) (int, int) {
	c := externalPosition(x, 2)
	r := externalPosition(y, 4)
	return c, r
}

type Canvas struct {
	LineEnding string
	Inverse    bool
	chars      map[int]map[int]int
}

// Make a new canvas
func NewCanvas() Canvas {
	c := Canvas{LineEnding: "\n", Inverse: false}
	c.Clear()
	return c
}

func (c Canvas) MaxY() int {
	max := 0
	for k, _ := range c.chars {
		if k > max {
			max = k
		}
	}
	return max * 4
}

func (c Canvas) MinY() int {
	min := 0
	for k, _ := range c.chars {
		if k < min {
			min = k
		}
	}
	return min * 4
}

func (c Canvas) MaxX() int {
	max := 0
	for _, v := range c.chars {
		for k, _ := range v {
			if k > max {
				max = k
			}
		}
	}
	return max * 2
}

func (c Canvas) MinX() int {
	min := 0
	for _, v := range c.chars {
		for k, _ := range v {
			if k < min {
				min = k
			}
		}
	}
	return min * 2
}

// Clear all pixels
func (c *Canvas) Clear() {
	c.chars = make(map[int]map[int]int)
}

// Set a pixel of c
func (c *Canvas) Set(x, y int) {
	col, row := getPos(x, y)
	if m := c.chars[row]; m == nil {
		c.chars[row] = make(map[int]int)
	}
	val := c.chars[row][col]
	mapv := getDot(y, x, c.Inverse)
	c.chars[row][col] = val | mapv
}

// Unset a pixel of c
func (c *Canvas) UnSet(x, y int) {
	col, row := getPos(x, y)
	if m := c.chars[row]; m == nil {
		c.chars[row] = make(map[int]int)
	}
	c.chars[row][col] &^= getDot(y, x, c.Inverse)
}

// Toggle a point
func (c *Canvas) Toggle(x, y int) {
	col, row := getPos(x, y)
	if m := c.chars[row]; m == nil {
		c.chars[row] = make(map[int]int)
	}
	c.chars[row][col] ^= getDot(y, x, c.Inverse)
}

// Set text to the given coordinates
func (c *Canvas) SetText(x, y int, text string) {
	col, row := getPos(x, y)
	if m := c.chars[row]; m == nil {
		c.chars[row] = make(map[int]int)
	}
	for i, char := range text {
		c.chars[row][col+i] = int(char) - brailleStartOrdinal
	}
}

// Get pixel at the given coordinates
func (c Canvas) Get(x, y int) bool {
	dot := getDot(y, x, c.Inverse)
	col, row := getPos(x, y)
	char := c.chars[row][col]
	return (char & dot) != 0
}

// Get character at the given screen coordinates
func (c Canvas) GetScreenCharacter(x, y int) rune {
	return rune(c.chars[y][x] + brailleStartOrdinal)
}

// Get character for the given pixel
func (c Canvas) GetCharacter(x, y int) rune {
	return c.GetScreenCharacter(x/4, y/4)
}

// Retrieve the rows from a given view
func (c Canvas) Rows(minX, minY, maxX, maxY int) []string {
	minRow, maxRow := minY/4, maxY/4
	minCol, maxCol := minX/2, maxX/2

	txts := make([]string, 0)
	if c.Inverse {
		for row := maxRow; row >= minRow; row-- {
			txts = append(txts, c.line(row, minCol, maxCol))
		}
	} else {
		for row := minRow; row <= maxRow; row++ {
			txts = append(txts, c.line(row, minCol, maxCol))
		}
	}
	return txts
}

func (c Canvas) line(row, minCol, maxCol int) string {
	txt := []rune{}
	for col := minCol; col <= maxCol; col++ {
		char := c.chars[row][col]
		txt = append(txt, rune(char+brailleStartOrdinal))
	}
	return string(txt)
}

// Retrieve a string representation of the frame at the given parameters
func (c Canvas) Frame(minX, minY, maxX, maxY int) string {
	var txt string
	for _, row := range c.Rows(minX, minY, maxX, maxY) {
		txt += row
		txt += c.LineEnding
	}
	return txt
}

func (c Canvas) String() string {
	return c.Frame(c.MinX(), c.MinY(), c.MaxX(), c.MaxY())
}

func (c *Canvas) DrawLine(x1, y1, x2, y2 float64) {
	xdiff := math.Abs(x1 - x2)
	ydiff := math.Abs(y2 - y1)

	var xdir, ydir float64
	if x1 <= x2 {
		xdir = 1
	} else {
		xdir = -1
	}
	if y1 <= y2 {
		ydir = 1
	} else {
		ydir = -1
	}

	r := math.Max(xdiff, ydiff)

	for i := 0; i < round(r)+1; i = i + 1 {
		x, y := x1, y1
		if ydiff != 0 {
			y += (float64(i) * ydiff) / (r * ydir)
		}
		if xdiff != 0 {
			x += (float64(i) * xdiff) / (r * xdir)
		}
		c.Toggle(round(x), round(y))
	}
}

func (c *Canvas) DrawPolygon(center_x, center_y, sides, radius float64) {
	degree := 360 / sides
	for n := 0; n < int(sides); n = n + 1 {
		a := float64(n) * degree
		b := float64(n+1) * degree

		x1 := (center_x + (math.Cos(radians(a)) * (radius/2 + 1)))
		y1 := (center_y + (math.Sin(radians(a)) * (radius/2 + 1)))
		x2 := (center_x + (math.Cos(radians(b)) * (radius/2 + 1)))
		y2 := (center_y + (math.Sin(radians(b)) * (radius/2 + 1)))

		c.DrawLine(x1, y1, x2, y2)
	}
}

func radians(d float64) float64 {
	return d * (math.Pi / 180)
}

func round(x float64) int {
	return int(x + 0.5)
}
