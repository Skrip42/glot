package glot

import (
	"fmt"
)

// A pointGroup refers to a set of points that need to plotted.
// It could either be a set of points or a function of co-ordinates.
// For Example z = Function(x,y)(3 Dimensional) or  y = Function(x) (2-Dimensional)
type pointGroup struct {
	name       string // Name of the curve
	dimensions int    // dimensions of the curve
	style      string // current plotting style
	data       any    // Data inside the curve in any integer/float format
	castedData any    // The data inside the curve typecasted to float64
	set        bool   //
}

type number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func toFloat64[T number](data []T) []float64 {
	result := make([]float64, len(data))
	for i, val := range data {
		result[i] = float64(val)
	}
	return result
}

func to2DFloat64[T number](data [][]T) [][]float64 {
	result := make([][]float64, len(data))
	for i, val := range data {
		result[i] = toFloat64(val)
	}
	return result
}

func (plot *plot) addMultiDimensionPointGroup(name string, style Style, points [][]float64) error {
	if plot.dimensions != len(points) {
		return &gnuplotError{fmt.Sprintf("The dimensions of this PointGroup are not compatible with the dimensions of the plot.\nIf you want to make a 2-d curve you must specify a 2-d plot.")}
	}

	curve := &pointGroup{name: name, dimensions: plot.dimensions, data: points, set: true, style: string(style)}

	curve.castedData = points
	if plot.dimensions == 2 {
		plot.plotXY(curve)
	} else {
		plot.plotXYZ(curve)
	}
	plot.pointGroup[name] = curve

	return nil
}

func (plot *plot) addOnewDimensionPointGroup(name string, style Style, points []float64) error {
	curve := &pointGroup{name: name, dimensions: plot.dimensions, data: points, set: true, style: string(style)}
	curve.castedData = points
	plot.plotX(curve)
	plot.pointGroup[name] = curve
	return nil
}

// AddPointGroup function adds a group of points to a plot.
//
// Usage
//
//	dimensions := 2
//	persist := false
//	debug := false
//	plot, _ := glot.NewPlot(dimensions, persist, debug)
//	plot.AddPointGroup("Sample1", "points", []int32{51, 8, 4, 11})
//	plot.AddPointGroup("Sample2", "points", []int32{1, 2, 4, 11})
//	plot.SavePlot("1.png")
func (plot *plot) AddPointGroup(name string, style Style, data any) (err error) {
	_, exists := plot.pointGroup[name]
	if exists {
		return &gnuplotError{fmt.Sprintf("A PointGroup with the name %s already exists, please use another name of the curve or remove this curve before using another one with the same name.", name)}
	}

	switch v := data.(type) {
	case [][]float64:
		plot.addMultiDimensionPointGroup(name, style, v)
	case [][]float32:
		plot.addMultiDimensionPointGroup(name, style, to2DFloat64(v))
	case [][]int:
		plot.addMultiDimensionPointGroup(name, style, to2DFloat64(v))
	case [][]int8:
		plot.addMultiDimensionPointGroup(name, style, to2DFloat64(v))
	case [][]int16:
		plot.addMultiDimensionPointGroup(name, style, to2DFloat64(v))
	case [][]int32:
		plot.addMultiDimensionPointGroup(name, style, to2DFloat64(v))
	case [][]int64:
		plot.addMultiDimensionPointGroup(name, style, to2DFloat64(v))
	case []float64:
		plot.addOnewDimensionPointGroup(name, style, v)
	case []float32:
		plot.addOnewDimensionPointGroup(name, style, toFloat64(v))
	case []int:
		plot.addOnewDimensionPointGroup(name, style, toFloat64(v))
	case []int8:
		plot.addOnewDimensionPointGroup(name, style, toFloat64(v))
	case []int16:
		plot.addOnewDimensionPointGroup(name, style, toFloat64(v))
	case []int32:
		plot.addOnewDimensionPointGroup(name, style, toFloat64(v))
	case []int64:
		plot.addOnewDimensionPointGroup(name, style, toFloat64(v))
	default:
		return &gnuplotError{fmt.Sprintf("invalid number of dims ")}

	}
	return err
}

// RemovePointGroup helps to remove a particular point group from the plot.
// This way you can remove a pointgroup if it's un-necessary.
//
// Usage
//
//	dimensions := 3
//	persist := false
//	debug := false
//	plot, _ := glot.NewPlot(dimensions, persist, debug)
//	plot.AddPointGroup("Sample1", "points", []int32{51, 8, 4, 11})
//	plot.AddPointGroup("Sample2", "points", []int32{1, 2, 4, 11})
//	plot.RemovePointGroup("Sample1")
func (plot *plot) RemovePointGroup(name string) {
	delete(plot.pointGroup, name)
	plot.cleanplot()
	for _, pointGroup := range plot.pointGroup {
		plot.plotX(pointGroup)
	}
}

// ResetPointGroupStyle helps to reset the style of a particular point group in a plot.
// Using both AddPointGroup and RemovePointGroup you can add or remove point groups.
// And dynamically change the plots.
//
// Usage
//
//	dimensions := 2
//	persist := false
//	debug := false
//	plot, _ := glot.NewPlot(dimensions, persist, debug)
//	plot.AddPointGroup("Sample1", "points", []int32{51, 8, 4, 11})
//	plot.ResetPointGroupStyle("Sample1", "points")
func (plot *plot) ResetPointGroupStyle(name string, style string) (err error) {
	pointGroup, exists := plot.pointGroup[name]
	if !exists {
		return &gnuplotError{fmt.Sprintf("A curve with name %s does not exist.", name)}
	}
	plot.RemovePointGroup(name)
	pointGroup.style = style
	plot.plotX(pointGroup)
	return err
}
