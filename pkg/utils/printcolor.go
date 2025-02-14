// File: pkg/utils/printcolor.go

package utils

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/fatih/color"
)

/*
FprintColored writes a colored line to the provided writer.

Parameters:
  - w: The io.Writer where output is written.
  - prefix: The string to print in the specified color (defaults to white if no color is provided).
  - secondary: The string printed immediately after the colored prefix.
  - attrs: Variadic color attributes; if none are provided, white is used.

Usage:

	FprintColored(os.Stdout, "Loaded config from: ", configPath, color.FgHiGreen)

Notes:

	If 'secondary' is empty, only the colored prefix is printed.
*/
func FprintColored(w io.Writer, prefix, secondary string, attrs ...color.Attribute) {
	var c *color.Color
	if len(attrs) > 0 {
		c = color.New(attrs...)
	} else {
		c = color.New(color.FgWhite)
	}
	if secondary != "" {
		fmt.Fprintf(w, "%s%s\n", c.Sprint(prefix), secondary)
	} else {
		fmt.Fprintln(w, c.Sprint(prefix))
	}
}

/*
FprintColoredDynamic writes multiple colored string segments to the provided writer on one line.

Parameters:
  - w: The io.Writer for output.
  - texts: A slice of strings to print sequentially.
  - colors: A slice of color attributes corresponding to each text.
    If there are more texts than colors, the last provided color is used for the remaining texts.

Usage:

	FprintColoredDynamic(os.Stdout, []string{"A ", "B ", "C"}, []color.Attribute{color.FgHiGreen, color.FgHiMagenta})

Notes:

	All text segments are printed on the same line, followed by a newline.
*/
func FprintColoredDynamic(w io.Writer, texts []string, colors []color.Attribute) {
	for i, text := range texts {
		var attr color.Attribute
		if i < len(colors) {
			attr = colors[i]
		} else if len(colors) > 0 {
			attr = colors[len(colors)-1]
		} else {
			attr = color.FgWhite
		}
		fmt.Fprint(w, color.New(attr).Sprint(text))
	}
	fmt.Fprintln(w)
}

/*
PrintColoredDynamicToStdout is a convenience function that writes dynamic colored output to os.Stdout.
*/
func PrintColoredDynamicToStdout(texts []string, colors []color.Attribute) {
	FprintColoredDynamic(os.Stdout, texts, colors)
}

/*
PrintColored is the main exported function for this utility.
It dynamically determines how to print colored output based on the types of arguments passed.

Usage:
 1. To print a single string:
    PrintColored("Just a string")
 2. To print a prefix and secondary string with a color:
    PrintColored("Prefix: ", "Secondary", color.FgHiGreen)
 3. To print multiple segments with individual colors:
    PrintColored([]string{"Segment1 ", "Segment2 ", "Segment3"},
    []color.Attribute{color.FgHiGreen, color.FgHiMagenta})

Behavior:
  - If the first argument is a []string, it expects a second argument as []color.Attribute and calls the dynamic printer.
  - Otherwise, if the first argument is a string:
  - With one argument, prints that string in white.
  - With two arguments (both strings), prints the first in white and the second in default formatting.
  - With additional arguments of type color.Attribute (or a slice thereof), uses them to color the prefix.
*/
func PrintColored(args ...interface{}) {
	if len(args) == 0 {
		return
	}

	// Dynamic mode: if the first argument is a []string.
	if texts, ok := args[0].([]string); ok {
		var colors []color.Attribute
		if len(args) > 1 {
			if cols, ok := args[1].([]color.Attribute); ok {
				colors = cols
			}
		}
		PrintColoredDynamicToStdout(texts, colors)
		return
	}

	// Otherwise, assume the first argument is a string.
	prefix, ok := args[0].(string)
	if !ok {
		return
	}

	secondary := ""
	if len(args) >= 2 {
		if sec, ok := args[1].(string); ok {
			secondary = sec
		}
	}

	var attrs []color.Attribute
	if len(args) > 2 {
		// Collect any color attributes (supports individual values or a slice).
		for _, arg := range args[2:] {
			v := reflect.ValueOf(arg)
			switch v.Kind() {
			case reflect.Slice:
				for i := 0; i < v.Len(); i++ {
					item := v.Index(i).Interface()
					if attr, ok := item.(color.Attribute); ok {
						attrs = append(attrs, attr)
					}
				}
			default:
				if attr, ok := arg.(color.Attribute); ok {
					attrs = append(attrs, attr)
				}
			}
		}
	}

	if len(attrs) == 0 {
		attrs = append(attrs, color.FgWhite)
	}

	FprintColored(os.Stdout, prefix, secondary, attrs...)
}
