package glot

import (
	"fmt"
	"os"
)

func (plot *plot) plotX(pointGroup *pointGroup) error {
	f, err := os.CreateTemp(os.TempDir(), gGnuplotPrefix)
	if err != nil {
		return err
	}
	fname := f.Name()
	plot.tmpFiles[fname] = f
	for _, d := range pointGroup.castedData.([]float64) {
		f.WriteString(fmt.Sprintf("%v\n", d))
	}
	f.Close()
	cmd := plot.plotCmd
	if plot.nPlots > 0 {
		cmd = plotCommand
	}
	if pointGroup.style == "" {
		pointGroup.style = defaultStyle
	}
	var line string
	if pointGroup.name == "" {
		line = fmt.Sprintf("%s \"%s\" with %s", cmd, fname, pointGroup.style)
	} else {
		line = fmt.Sprintf("%s \"%s\" title \"%s\" with %s",
			cmd, fname, pointGroup.name, pointGroup.style)
	}
	plot.nPlots++
	return plot.cmd(line)
}

func (plot *plot) plotXY(pointGroup *pointGroup) error {
	x := pointGroup.castedData.([][]float64)[0]
	y := pointGroup.castedData.([][]float64)[1]
	npoints := min(len(x), len(y))

	f, err := os.CreateTemp(os.TempDir(), gGnuplotPrefix)
	if err != nil {
		return err
	}
	fname := f.Name()
	plot.tmpFiles[fname] = f

	for i := range npoints {
		f.WriteString(fmt.Sprintf("%v %v\n", x[i], y[i]))
	}

	f.Close()
	cmd := plot.plotCmd
	if plot.nPlots > 0 {
		cmd = plotCommand
	}

	if pointGroup.style == "" {
		pointGroup.style = "points"
	}
	var line string
	if pointGroup.name == "" {
		line = fmt.Sprintf("%s \"%s\" with %s", cmd, fname, pointGroup.style)
	} else {
		line = fmt.Sprintf("%s \"%s\" title \"%s\" with %s",
			cmd, fname, pointGroup.name, pointGroup.style)
	}
	plot.nPlots++
	return plot.cmd(line)
}

func (plot *plot) plotXYZ(points *pointGroup) error {
	x := points.castedData.([][]float64)[0]
	y := points.castedData.([][]float64)[1]
	z := points.castedData.([][]float64)[2]
	npoints := min(len(x), len(y))
	npoints = min(npoints, len(z))
	f, err := os.CreateTemp(os.TempDir(), gGnuplotPrefix)
	if err != nil {
		return err
	}
	fname := f.Name()
	plot.tmpFiles[fname] = f

	for i := range npoints {
		f.WriteString(fmt.Sprintf("%v %v %v\n", x[i], y[i], z[i]))
	}

	f.Close()
	cmd := "splot" // Force 3D plot
	if plot.nPlots > 0 {
		cmd = plotCommand
	}

	var line string
	if points.name == "" {
		line = fmt.Sprintf("%s \"%s\" with %s", cmd, fname, points.style)
	} else {
		line = fmt.Sprintf("%s \"%s\" title \"%s\" with %s",
			cmd, fname, points.name, points.style)
	}
	plot.nPlots++
	return plot.cmd(line)
}
