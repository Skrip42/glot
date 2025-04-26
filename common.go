package glot

import (
	"fmt"
	"strconv"
)

// SetTitle sets the title for the plot
//
// Usage
//
//	dimensions := 3
//	persist := false
//	debug := false
//	plot, _ := glot.NewPlot(dimensions, persist, debug)
//	plot.AddPointGroup("Sample 1", "lines", []float64{2, 3, 4, 1})
//	plot.SetTitle("Test Results")
func (plot *plot) SetTitle(title string) error {
	return plot.cmd(fmt.Sprintf("set title \"%s\" ", title))
}

// SetXLabel changes the label for the x-axis
//
// Usage
//
//	 dimensions := 3
//	 persist := false
//	 debug := false
//	 plot, _ := glot.NewPlot(dimensions, persist, debug)
//	 plot.AddPointGroup("Sample 1", "lines", []float64{2, 3, 4, 1})
//	 plot.SetTitle("Test Results")
//		plot.SetXLabel("X-Axis")
func (plot *plot) SetXLabel(label string) error {
	return plot.cmd(fmt.Sprintf("set xlabel '%s'", label))
}

// SetYLabel changes the label for the y-axis
//
// Usage
//
//	 dimensions := 3
//	 persist := false
//	 debug := false
//	 plot, _ := glot.NewPlot(dimensions, persist, debug)
//	 plot.AddPointGroup("Sample 1", "lines", []float64{2, 3, 4, 1})
//	 plot.SetTitle("Test Results")
//		plot.SetYLabel("Y-Axis")
func (plot *plot) SetYLabel(label string) error {
	return plot.cmd(fmt.Sprintf("set ylabel '%s'", label))
}

// SetZLabel changes the label for the z-axis
//
// Usage
//
//	 dimensions := 3
//	 persist := false
//	 debug := false
//	 plot, _ := glot.NewPlot(dimensions, persist, debug)
//	 plot.AddPointGroup("Sample 1", "lines", []float64{2, 3, 4, 1})
//	 plot.SetTitle("Test Results")
//		plot.SetZLabel("Z-Axis")
func (plot *plot) SetZLabel(label string) error {
	return plot.cmd(fmt.Sprintf("set zlabel '%s'", label))
}

func (plot *plot) SetGrid() error {
	return plot.cmd("set grid")
}

// SetLabels Functions helps to set labels for x, y, z axis  simultaneously
//
// Usage
//
//	 dimensions := 3
//	 persist := false
//	 debug := false
//	 plot, _ := glot.NewPlot(dimensions, persist, debug)
//	 plot.AddPointGroup("Sample 1", "lines", []float64{2, 3, 4, 1})
//	 plot.SetTitle("Test Results")
//		plot.SetLabels("X-axis","Y-Axis","Z-Axis")
func (plot *plot) SetLabels(labels ...string) error {
	ndims := len(labels)
	if ndims > 3 || ndims <= 0 {
		return &gnuplotError{fmt.Sprintf("invalid number of dims '%v'", ndims)}
	}
	slabelFunc := []func(string) error{plot.SetXLabel, plot.SetYLabel, plot.SetZLabel}

	for i, label := range labels {
		err := slabelFunc[i](label)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetXrange changes the label for the x-axis
//
// Usage
//
//	 dimensions := 3
//	 persist := false
//	 debug := false
//	 plot, _ := glot.NewPlot(dimensions, persist, debug)
//	 plot.AddPointGroup("Sample 1", "lines", []float64{2, 3, 4, 1})
//	 plot.SetTitle("Test Results")
//		plot.SetXrange(-2,2)
func (plot *plot) SetXrange(start int, end int) error {
	return plot.cmd(fmt.Sprintf("set xrange [%d:%d]", start, end))
}

// SetLogscale changes the label for the x-axis
//
// Usage
//
//	dimensions := 3
//	persist := false
//	debug := false
//	plot, _ := glot.NewPlot(dimensions, persist, debug)
//	plot.SetYrange(-2, 18)
//	plot.AddPointGroup("rates", "circle", [][]float64{{2, 4, 8, 16, 32}, {4, 7, 4, 10, 3}})
//	plot.SetLogscale("x", 2)
func (plot *plot) SetLogscale(axis string, base int) error {
	return plot.cmd(fmt.Sprintf("set logscale %s %d", axis, base))
}

// SetYrange changes the label for the y-axis
//
// Usage
//
//	 dimensions := 3
//	 persist := false
//	 debug := false
//	 plot, _ := glot.NewPlot(dimensions, persist, debug)
//	 plot.AddPointGroup("Sample 1", "lines", []float64{2, 3, 4, 1})
//	 plot.SetTitle("Test Results")
//		plot.SetYrange(-2,2)
func (plot *plot) SetYrange(start int, end int) error {
	return plot.cmd(fmt.Sprintf("set yrange [%d:%d]", start, end))
}

// SetZrange changes the label for the z-axis
//
// Usage
//
//	 dimensions := 3
//	 persist := false
//	 debug := false
//	 plot, _ := glot.NewPlot(dimensions, persist, debug)
//	 plot.AddPointGroup("Sample 1", "lines", []float64{2, 3, 4, 1})
//	 plot.SetTitle("Test Results")
//		plot.SetZrange(-2,2)
func (plot *plot) SetZrange(start int, end int) error {
	return plot.cmd(fmt.Sprintf("set zrange [%d:%d]", start, end))
}

// SavePlot function is used to save the plot at this point.
// The plot is dynamic and additional pointgroups can be added and removed and different versions
// of the same plot can be saved.
//
// Usage
//
//	 dimensions := 3
//	 persist := false
//	 debug := false
//	 plot, _ := glot.NewPlot(dimensions, persist, debug)
//	 plot.AddPointGroup("Sample 1", "lines", []float64{2, 3, 4, 1})
//	 plot.SetTitle("Test Results")
//		plot.SetZrange(-2,2)
//	 plot.SavePlot("1.jpeg")
func (plot *plot) SavePlot(filename string, weight, height int) error {
	if plot.nPlots == 0 {
		return &gnuplotError{fmt.Sprintf("This plot has 0 curves and therefore its a redundant plot and it can't be printed.")}
	}
	outputFormat := "set terminal " + string(plot.format) +
		" size " + strconv.Itoa(weight) + ", " + strconv.Itoa(height)

	_ = plot.cmd(outputFormat)
	outputFileCommand := "set output " + "'" + filename + "'"
	_ = plot.cmd(outputFileCommand)
	_ = plot.cmd("replot  ")
	return nil
}

// SetFormat function is used to save the plot at this point.
// The plot is dynamic and additional pointgroups can be added and removed and different versions
// of the same plot can be saved.
//
// Usage
//
//	 dimensions := 3
//	 persist := false
//	 debug := false
//	 plot, _ := glot.NewPlot(dimensions, persist, debug)
//	 plot.AddPointGroup("Sample 1", "lines", []float64{2, 3, 4, 1})
//	 plot.SetTitle("Test Results")
//		plot.SetFormat("pdf")
//	 plot.SavePlot("1.pdf")
//
// NOTE: png is default format for saving files.
func (plot *plot) SetFormat(newformat Format) error {
	plot.format = newformat
	return nil
}

func (plot *plot) SetKeyOutside() error {
	return plot.cmd("set key outside")
}
