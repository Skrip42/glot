// Glot is a library for having simplified 1,2,3 Dimensional points/line plots
// It's built on top of Gnu plot and offers the ability to use Raw Gnu plot commands
// directly from golang.
// See the gnuplot documentation page for the exact semantics of the gnuplot
// commands.
//  http://www.gnuplot.info/

package glot

import (
	"fmt"
	"sync"
)

type Style string

const (
	StyleLines        Style = "lines"
	StylePoints       Style = "points"
	StyleLinepoints   Style = "linepoints"
	StyleImpulses     Style = "impulses"
	StyleDots         Style = "dots"
	StyleBar          Style = "bar"
	StyleFillSolid    Style = "fill solid"
	StyleHistogram    Style = "histogram"
	StyleCircle       Style = "circle"
	StyleErrorBars    Style = "errorbars"
	StyleBoxErrorBars Style = "boxerrorbars"
	StyleBoxes        Style = "boxes"
	StyleLp           Style = "lp"
)

type Format string

const (
	FormatPng Format = "png"
	FormatPdf Format = "pdf"
)

// Plot is the basic type representing a plot.
// Every plot has a set of Pointgroups that are simultaneously plotted
// on a 2/3 D plane given the plot type.
// The Plot dimensions must be specified at the time of construction
// and can't be changed later.  All the Pointgroups added to a plot must
// have same dimensions as the dimension specified at the
// the time of plot construction.
// The Pointgroups can be dynamically added and removed from a plot
// And style changes can also be made dynamically.
// Plot is an interface for plotting data in 1, 2 or 3 dimensions
type Plot interface {
	// AddPointGroup adds a new point group with the given name and style
	AddPointGroup(name string, style Style, points any) error

	// RemovePointGroup helps to remove a particular point group from the plot.
	RemovePointGroup(name string)

	// ResetPointGroupStyle helps to reset the style of a particular point group in a plot.
	ResetPointGroupStyle(name string, style string) error

	// SetTitle sets the title for the plot
	SetTitle(title string) error

	// SetXLabel sets the label for the x-axis
	SetXLabel(label string) error

	// SetYLabel sets the label for the y-axis
	SetYLabel(label string) error

	// SetZLabel sets the label for the z-axis
	SetZLabel(label string) error

	// SetLabels sets labels for x, y, z axes simultaneously
	SetLabels(labels ...string) error

	// SetXrange sets the range for the x-axis
	SetXrange(start int, end int) error

	// SetLogscale changes the label for the x-axis
	SetLogscale(axis string, base int) error

	// SetYrange changes the label for the y-axis
	SetYrange(start int, end int) error

	// SetZrange changes the label for the z-axis
	SetZrange(start int, end int) error

	// SavePlot function is used to save the plot at this point.
	SavePlot(filename string, w, h int) error

	// SetFormat sets the output format (png, pdf, etc)
	SetFormat(format Format) error

	SetKeyOutside() error

	SetGrid() error
}

// plot implements the Plot interface
type plot struct {
	proc       *plotterProcess
	plotCmd    string
	nPlots     int                    // number of currently active plots
	tmpFiles   tempFilesDb            // A temporary file used for saving data
	dimensions int                    // dimensions of the plot
	pointGroup map[string]*pointGroup // A map between Curve name and curve type. This maps a name to a given curve in a plot. Only one curve with a given name exists in a plot.
	format     Format                 // The saving format of the plot. This could be PDF, PNG, JPEG and so on.
	style      string                 // style of the plot
	title      string                 // The title of the plot.
}

// NewPlot Function makes a new plot with the specified dimensions.
//
// Usage
//
//	dimensions := 3
//	persist := false
//	debug := false
//	plot, _ := glot.NewPlot(dimensions, persist, debug)
//
// Variable definitions
//
//	dimensions  :=> refers to the dimensions of the plot.
//	persist     :=> used to make the gnu plot window stay open.
func NewPlot(dimensions int, persist bool) (Plot, error) {
	err := sync.OnceValue(initialize)()
	if err != nil {
		return nil, err
	}
	p := &plot{proc: nil, plotCmd: "plot",
		nPlots: 0, dimensions: dimensions, style: "points", format: "png"}
	p.pointGroup = make(map[string]*pointGroup) // Adding a mapping between a curve name and a curve
	p.tmpFiles = make(tempFilesDb)
	proc, err := newPlotterProc(persist)
	if err != nil {
		return nil, err
	}
	// Only 1,2,3 Dimensional plots are supported
	if dimensions > 3 || dimensions < 1 {
		return nil, &gnuplotError{fmt.Sprintf("invalid number of dims '%v'", dimensions)}
	}
	p.proc = proc
	return p, nil
}
