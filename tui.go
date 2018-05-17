package main

import (
	"fmt"
	"math"

	"github.com/Kerrigan29a/drawille-go"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type TUI struct {
	app *tview.Application
}

func NewTUI(reader *Reader) *TUI {
	/* Config widgets */
	app := tview.NewApplication()
	logTV := tview.NewTextView().SetDynamicColors(true)
	box := tview.NewBox().SetBorder(true).SetTitle("[ " + reader.path + " ]")
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(box, 0, 3, true)
	flex.AddItem(logTV, 0, 1, false)
	app.SetRoot(flex, true)

	t := &TUI{app: app}

	/* Set function to draw canvas */
	cc := NewCanvasController(reader, box, logTV)
	box.SetDrawFunc(cc.DrawFunc)
	box.SetInputCapture(cc.InputCapture)

	return t
}

func (t *TUI) Run() error {
	if err := t.app.Run(); err != nil {
		return err
	}
	return nil
}

type CanvasController struct {
	reader *Reader
	box    *tview.Box
	logTV  *tview.TextView
	canvas drawille.Canvas
	amount int
}

func NewCanvasController(reader *Reader, box *tview.Box, logTV *tview.TextView) *CanvasController {
	canvas := drawille.NewCanvas()
	canvas.Inverse = true
	return &CanvasController{reader: reader, box: box, logTV: logTV, canvas: canvas, amount: -1}
}

func (c *CanvasController) Values(amount int) ([]float64, error) {
	prevAmount := c.reader.FrameCap()
	if prevAmount != amount {
		c.reader.SetFrameCap(amount)
		difference := amount - prevAmount
		if difference >= 0 {
			err := c.reader.ReadNext(uint(difference))
			if err != nil {
				return nil, err
			}
		}
	}
	return c.reader.Values()
}

func (c *CanvasController) DrawFunc(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	innerX := x + 1
	innerY := y + 1
	innerWidth := width - 2
	innerHeight := height - 2

	begin := 0
	amount := innerWidth * 2
	end := begin + amount

	/* Read values */
	values, err := c.Values(amount)
	if err != nil {
		panic(err)
	}
	logf(c.logTV, "[blue]DEBUG[-] values = %v\n", values)

	if len(values) == 0 {
		return innerX, innerY, innerWidth, innerHeight
	}

	/* Create graph */
	c.canvas.Clear()
	normalizedValues := normalizeFloat64(values)
	intValues := floatToInt(normalizedValues[begin:end], (innerHeight-1)*4)
	for i := begin; i < end; i++ {
		gX := i
		gY := intValues[i]
		//logf(c.logTV, "[blue]DEBUG[-] %d -> %g -> %g -> %d\n", gX, values[i], normalizedValues[i], gY)

		for j := 0; j < gY; j++ {
			c.canvas.Set(gX, j)
		}
		c.canvas.Set(gX, gY)
	}
	gWidth := c.canvas.MaxX() - c.canvas.MinX()
	gHeight := c.canvas.MaxY() - c.canvas.MinY()
	logf(c.logTV, "[blue]DEBUG[-] gWidth = %v/2 = %v\n", gWidth, gWidth/2)
	logf(c.logTV, "[blue]DEBUG[-] gHeight = %v/4 = %v\n", gHeight, gHeight/4)

	currentY := innerY
	for _, row := range c.canvas.Rows(c.canvas.MinX(), c.canvas.MinY(), c.canvas.MaxX(), c.canvas.MaxY()) {
		tview.Print(screen, row, innerX, currentY, innerWidth, tview.AlignLeft, tcell.ColorGreen)
		currentY++
	}
	return innerX, innerY, innerWidth, innerHeight
}

func (c *CanvasController) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	minAmount := uint(math.Round(0.1 * float64(c.reader.FrameCap())))
	maxAmount := uint(math.Round(0.5 * float64(c.reader.FrameCap())))
	switch event.Key() {
	case tcell.KeyRight:
		amount := minAmount
		if event.Modifiers()&tcell.ModShift != 0 {
			amount = maxAmount
		}
		c.reader.ReadNext(amount)
		return nil
	case tcell.KeyLeft:
		amount := minAmount
		if event.Modifiers()&tcell.ModShift != 0 {
			amount = maxAmount
		}
		c.reader.ReadPrev(amount)
		return nil
	case tcell.KeyRune:
		switch event.Rune() {
		case 'h':
			c.reader.ReadPrev(minAmount)
			return nil
		case 'H':
			c.reader.ReadPrev(maxAmount)
			return nil
		case 'l':
			c.reader.ReadNext(minAmount)
			return nil
		case 'L':
			c.reader.ReadNext(maxAmount)
			return nil
		default:
			return event
		}
	default:
		return event
	}
}

func logf(logTV *tview.TextView, format string, v ...interface{}) {
	txt := fmt.Sprintf(format, v...)
	logTV.Write([]byte(txt))
}
