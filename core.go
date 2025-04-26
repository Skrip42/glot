package glot

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

var gGnuplotCmd string
var gGnuplotPrefix = "go-gnuplot-"

const defaultStyle = "points" // The default style for a curve
const plotCommand = "replot"  // The default style for a curve

// A map between os files and file names
type tempFilesDb map[string]*os.File

// Function to intialize the package and check for GNU plot installation
// This raises an error if GNU plot is not installed
func initialize() error {
	var err error

	gnuplotExecutableName := "gnuplot"

	if runtime.GOOS == "windows" {
		gnuplotExecutableName = "gnuplot.exe"
	}

	gGnuplotCmd, err = exec.LookPath(gnuplotExecutableName)
	if err != nil {
		return fmt.Errorf(
			"** could not find path to 'gnuplot':\n%v\n** set custom path to 'gnuplot' ", err)
	}
	return nil
}

type gnuplotError struct {
	err string
}

func (e *gnuplotError) Error() string {
	return e.err
}

// plotterProcess is the type for handling gnu commands.
type plotterProcess struct {
	handle *exec.Cmd
	stdin  io.WriteCloser
}

// NewPlotterProc function makes the plotterProcess struct
func newPlotterProc(persist bool) (*plotterProcess, error) {
	procArgs := []string{}
	if persist {
		procArgs = append(procArgs, "-persist")
	}
	cmd := exec.Command(gGnuplotCmd, procArgs...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	// to debug uncomment it
	// stderr, err := cmd.StderrPipe()
	// if err != nil {
	// 	return nil, err
	// }
	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	return nil, err
	// }
	// stderrReader := bufio.NewReader(stderr)
	// stdoutReader := bufio.NewReader(stdout)

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// go func() {
	// 	for {
	// 		line, _, err := stderrReader.ReadLine()
	// 		// line, _, err := bufio.NewReader(stderr).ReadLine()
	// 		fmt.Println("err>" + string(line))
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}
	// 	}
	// }()
	//
	// go func() {
	// 	for {
	// 		line, _, err := stdoutReader.ReadLine()
	// 		fmt.Println("out>" + string(line))
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}
	// 	}
	// }()

	return &plotterProcess{
		handle: cmd,
		stdin:  stdin,
	}, nil
}

// Cmd sends a command to the gnuplot subprocess and returns an error
// if something bad happened in the gnuplot process.
// ex:
//
//	fname := "foo.dat"
//	err := p.Cmd("plot %s", fname)
//	if err != nil {
//	  panic(err)
//	}
//
// func (plot *plot) cmd(command string) error {
func (plot *plot) cmd(command string) error {
	_, err := io.WriteString(plot.proc.stdin, command+"\n")

	return err
}

// Close makes sure all resources used by the gnuplot subprocess are reclaimed.
// This method is typically called when the Plotter instance is not needed
// anymore. That's usually done via a defer statement:
//
//	p, err := gnuplot.NewPlotter(...)
//	if err != nil { /* handle error */ }
//	defer p.Close()
func (plot *plot) Close() (err error) {
	if plot.proc != nil && plot.proc.handle != nil {
		plot.proc.stdin.Close()
		err = plot.proc.handle.Wait()
	}
	plot.resetPlot()
	return err
}

func (plot *plot) cleanplot() (err error) {
	plot.tmpFiles = make(tempFilesDb)
	plot.nPlots = 0
	return err
}

// ResetPlot is used to reset the whole plot.
// This removes all the PointGroup's from the plot and makes it new.
// Usage
//
//	plot.ResetPlot()
func (plot *plot) resetPlot() (err error) {
	plot.cleanplot()
	plot.pointGroup = make(map[string]*pointGroup) // Adding a mapping between a curve name and a curve
	return err
}

func SetCustomPathToGNUPlot(path string) {
	gGnuplotCmd = path
}
