// Package tracerr makes error output more informative.
// It adds stack trace to error and can display error with source fragments.
//
// Check example of output here https://github.com/ztrue/tracerr
package tracerr

import (
	"fmt"
	"runtime"
)

// DefaultCap is a default cap for frames array.
// It can be changed to number of expected frames
// for purpose of performance optimisation.
var DefaultCap = 20
var DefaultPrintStackMaxDepth = 5

// Error is an error with stack trace.
type Error interface {
	Error() string
	StackTrace() []Frame
	Unwrap() error
}

type errorData struct {
	// err contains original error.
	err error
	// frames contains stack trace of an error.
	frames []Frame
}

// CustomError creates an error with provided frames.
func SetStackMaxDepth(depth int) {
	if depth <= 0 {
		depth = 5
	}
	DefaultPrintStackMaxDepth = depth
}

// CustomError creates an error with provided frames.
func CustomError(err error, frames []Frame) Error {
	return &errorData{
		err:    err,
		frames: frames,
	}
}

// Errorf creates new error with stacktrace and formatted message.
// Formatting works the same way as in fmt.Errorf.
func Errorf(message string, args ...interface{}) Error {
	return trace(fmt.Errorf(message, args...), 2)
}

// New creates new error with stacktrace.
func New(message string, skip ...int) Error {
	if len(skip) > 0 {
		return trace(fmt.Errorf(message), skip[0])
	} else {
		return trace(fmt.Errorf(message), 2)
	}
}

// Wrap adds stacktrace to existing error.
func Wrap(err error, skip ...int) Error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if ok {
		return e
	}
	if len(skip) > 0 {
		return trace(err, skip[0])
	} else {
		return trace(err, 2)
	}
}

// Unwrap returns the original error.
func Unwrap(err error) error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if !ok {
		return err
	}
	return e.Unwrap()
}

// Error returns error message.
func (e *errorData) Error() string {
	//return Sprint(e.err)
	return e.err.Error()
}

// StackTrace returns stack trace of an error.
func (e *errorData) StackTrace() []Frame {
	return e.frames
}

// Unwrap returns the original error.
func (e *errorData) Unwrap() error {
	return e.err
}

// Frame is a single step in stack trace.
type Frame struct {
	// Func contains a function name.
	Func string
	// Line contains a line number.
	Line int
	// Path contains a file path.
	Path string
}

// StackTrace returns stack trace of an error.
// It will be empty if err is not of type Error.
func StackTrace(err error) []Frame {
	e, ok := err.(Error)
	if !ok {
		return nil
	}
	return e.StackTrace()
}

// String formats Frame to string.
func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s()", f.Path, f.Line, f.Func)
}

func trace(err error, skip int) Error {
	frames := make([]Frame, 0, DefaultCap)
	catchMaxCall := DefaultPrintStackMaxDepth
	for {
		catchMaxCall--
		pc, path, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		frame := Frame{
			Func: fn.Name(),
			Line: line,
			Path: path,
		}
		frames = append(frames, frame)
		skip++
		if catchMaxCall == 0 {
			break
		}
	}
	return &errorData{
		err:    err,
		frames: frames,
	}
}

const  _=`
package fmt

import (
	"io/ioutil"
	"runtime"
	"strings"
	"sync"
)

// DefaultLinesAfter is number of source lines after traced line to display.
var DefaultLinesAfter = 2

// DefaultLinesBefore is number of source lines before traced line to display.
var DefaultLinesBefore = 3

var cache = map[string][]string{}

var mutex sync.RWMutex

type Error interface {
	Error() string
	StackTrace() []Frame
	Unwrap() error
}

type Frame struct {
	// Func contains a function name.
	Func string
	// Line contains a line number.
	Line int
	// Path contains a file path.
	Path string
}

// Print prints error message with stack trace.
//func Print(err error) {
//	fmt.Println(Sprint(err))
//}

// PrintSource prints error message with stack trace and source fragments.
//
// By default 6 lines of source code will be printed,
// see DefaultLinesAfter and DefaultLinesBefore.
//
// Pass a single number to specify a total number of source lines.
//
// Pass two numbers to specify exactly how many lines should be shown
// before and after traced line.

type errorData struct {
	// err contains original error.
	err error
	// frames contains stack trace of an error.
	frames []Frame
}

// CustomError creates an error with provided frames.
func CustomError(err error, frames []Frame) Error {
	return &errorData{
		err:    err,
		frames: frames,
	}
}

// Errorf creates new error with stacktrace and formatted message.
// Formatting works the same way as in fmt.Errorf.
//func Errorf(message string, args ...interface{}) Error {
//	return trace(Errorf(message, args...), 2)
//}

// New creates new error with stacktrace.
func NewError(message string, a ...interface{}) Error {
	return trace(Errorf(message, a...), 2)
}

// Wrap adds stacktrace to existing error.

func WrapError(err error, skip ...int) Error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if ok {
		return e
	}
	if len(skip) > 0 {
		return trace(err, skip[0])
	} else {
		return trace(err, 2)
	}
}

var DefaultCap = 20

func trace(err error, skip int) Error {
	frames := make([]Frame, 0, DefaultCap)
	for {
		pc, path, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		frame := Frame{
			Func: fn.Name(),
			Line: line,
			Path: path,
		}
		frames = append(frames, frame)
		skip++
	}
	return &errorData{
		err:    err,
		frames: frames,
	}
}

// Unwrap returns the original error.
func Unwrap(err error) error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if !ok {
		return err
	}
	return e.Unwrap()
}

// Error returns error message.
func (e *errorData) Error() string {
	return e.err.Error()
}

// StackTrace returns stack trace of an error.
func (e *errorData) StackTrace() []Frame {
	return e.frames
}

// Unwrap returns the original error.
func (e *errorData) Unwrap() error {
	return e.err
}

func PrintSource(err error, nums ...int) {
	Println(SprintSource(err, nums...))
}

// PrintSourceColor prints error message with stack trace and source fragments,
// which are in color.
// Output rules are the same as in PrintSource.
func PrintSourceColor(err error, nums ...int) {
	Println(SprintSourceColor(err, nums...))
}

// Sprint returns error output by the same rules as Print.
//func Sprint(err error) string {
//	return sprint(err, []int{0}, false)
//}

func GetErrorStack(err error) string {
	return sprint(err, []int{0}, false)
}

// SprintSource returns error output by the same rules as PrintSource.
func SprintSource(err error, nums ...int) string {
	return sprint(err, nums, false)
}

// SprintSourceColor returns error output by the same rules as PrintSourceColor.
func SprintSourceColor(err error, nums ...int) string {
	return sprint(err, nums, true)
}

func calcRows(nums []int) (before, after int, withSource bool) {
	before = DefaultLinesBefore
	after = DefaultLinesAfter
	withSource = true
	if len(nums) > 1 {
		before = nums[0]
		after = nums[1]
		withSource = true
	} else if len(nums) == 1 {
		if nums[0] > 0 {
			// Extra line goes to "before" rather than "after".
			after = (nums[0] - 1) / 2
			before = nums[0] - after - 1
		} else {
			after = 0
			before = 0
			withSource = false
		}
	}
	if before < 0 {
		before = 0
	}
	if after < 0 {
		after = 0
	}
	return before, after, withSource
}

func readLines(path string) ([]string, error) {
	mutex.RLock()
	lines, ok := cache[path]
	mutex.RUnlock()
	if ok {
		return lines, nil
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, Errorf("tracerr: file %s not found", path)
	}
	lines = strings.Split(string(b), "\n")
	mutex.Lock()
	defer mutex.Unlock()
	cache[path] = lines
	return lines, nil
}

func sourceRows(rows []string, frame Frame, before, after int, colorized bool) []string {
	lines, err := readLines(frame.Path)
	if err != nil {
		message := err.Error()
		if colorized {
			message = message // aurora.Brown(message).String()
		}
		return append(rows, message, "")
	}
	if len(lines) < frame.Line {
		message := Sprintf(
			"tracerr: too few lines, got %d, want %d",
			len(lines), frame.Line,
		)
		if colorized {
			message = message // aurora.Brown(message).String()
		}
		return append(rows, message, "")
	}
	current := frame.Line - 1
	start := current - before
	end := current + after
	for i := start; i <= end; i++ {
		if i < 0 || i >= len(lines) {
			continue
		}
		line := lines[i]
		var message string
		// TODO Pad to the same length.
		if i == frame.Line-1 {
			message = Sprintf("%d\t%s", i+1, string(line))
			if colorized {
				message = message // aurora.Red(message).String()
			}
		} else if colorized {
			message = message // aurora.Sprintf("%d\t%s", aurora.Black(i+1), string(line))
		} else {
			message = Sprintf("%d\t%s", i+1, string(line))
		}
		rows = append(rows, message)
	}
	return append(rows, "")
}

func sprint(err error, nums []int, colorized bool) string {
	if err == nil {
		return ""
	}
	e, ok := err.(Error)
	if !ok {
		return err.Error()
	}
	before, after, withSource := calcRows(nums)
	frames := e.StackTrace()
	expectedRows := len(frames) + 1
	if withSource {
		expectedRows = (before+after+3)*len(frames) + 2
	}
	rows := make([]string, 0, expectedRows)
	rows = append(rows, e.Error())
	if withSource {
		rows = append(rows, "")
	}
	for _, frame := range frames {
		message := frame.String()
		if colorized {
			message = message // aurora.Bold(message).String()
		}
		rows = append(rows, message)
		if withSource {
			rows = sourceRows(rows, frame, before, after, colorized)
		}
	}
	return strings.Join(rows, "\n")
}

func (f Frame) String() string {
	return Sprintf("%s:%d %s()", f.Path, f.Line, f.Func)
}

`
