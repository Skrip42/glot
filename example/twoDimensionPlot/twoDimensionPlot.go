package main

import (
	"math"
	"time"

	"github.com/Skrip42/glot"
)

func main() {
	plot, err := glot.NewPlot(2, true)
	if err != nil {
		panic(err)
	}
	_ = plot.SetXLabel("X")
	_ = plot.SetYLabel("Y")
	_ = plot.SetGrid()
	_ = plot.SetYrange(0, 500)
	_ = plot.SetXrange(0, 100)
	_ = plot.SetTitle("this is 2d plot example")
	_ = plot.SetKeyOutside()

	styles := []glot.Style{
		glot.StyleLines,
		glot.StylePoints,
		// glot.StyleLinepoints,
		glot.StyleImpulses,
		glot.StyleDots,
		// glot.StyleBar,
		// glot.StyleFillSolid,
		// glot.StyleHistogram,
		glot.StyleCircle,
		glot.StyleErrorBars,
		glot.StyleBoxErrorBars,
		glot.StyleBoxes,
		glot.StyleLp,
	}

	for i, style := range styles {
		points := make([][]float64, 2)
		points[0] = make([]float64, 100)
		points[1] = make([]float64, 100)
		for x := range 100 {
			points[0][x] = float64(x)
			points[1][x] = (math.Pow(float64(x), 2) / 10) * float64(i)
		}
		err := plot.AddPointGroup(string(style), style, points)
		if err != nil {
			panic(err)
		}
	}
	err = plot.SavePlot("2dplot.png", 800, 1200)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
}
