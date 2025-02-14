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
  - w: An io.Writer where the output will be written.
  - prefix: The string to be printed in the specified color (or white if no color is provided).
  - secondary: The string printed immediately after the colored prefix in default formatting.
  - attrs: A variadic parameter of color attributes. If none are provided, white is used.

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
PrintColoredDynamic writes multiple colored string segments to the provided writer on the same line.

Parameters:
  - w: An io.Writer where the output will be written.
  - texts: A slice of strings to be printed sequentially.
  - colors: A slice of color attributes corresponding to each text segment.

Behavior:
  - For each text, if a color is specified at that index, that color is used.
  - If there are more texts than colors, the last provided color is used for the remaining texts.
  - All text segments are printed on the same line, and a newline is added at the end.
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

Usage:

	PrintColoredDynamicToStdout([]string{"One ", "Two ", "Three "}, []color.Attribute{color.FgRed, color.FgGreen, color.FgBlue})
*/
func PrintColoredDynamicToStdout(texts []string, colors []color.Attribute) {
	FprintColoredDynamic(os.Stdout, texts, colors)
}

/*
PrintColored is the main exported function for this util. It dynamically determines how to print colored output based on the arguments provided.

Usage:
 1. To print a simple colored line with a prefix and secondary text:
    PrintColored("Prefix: ", "Secondary", color.FgHiGreen)
 2. To print a single string in a default color:
    PrintColored("Just a string")
 3. To print multiple segments with individual colors:
    PrintColored([]string{"Segment1 ", "Segment2 ", "Segment3"},
    []color.Attribute{color.FgHiGreen, color.FgHiMagenta})

Behavior:
  - If the first argument is a []string, it expects the second argument to be a []color.Attribute and calls the dynamic multi-segment printer.
  - Otherwise, if the first argument is a string:
  - With only one argument, it prints that string in white.
  - With two arguments (both strings), it prints the first string in white followed by the second string in default formatting.
  - With additional arguments of type color.Attribute, it uses those attributes to color the prefix.
*/
func PrintColored(args ...interface{}) {
	// No arguments: nothing to print.
	if len(args) == 0 {
		return
	}

	// Check if the first argument is a slice of strings.
	if texts, ok := args[0].([]string); ok {
		// Expect a second argument as []color.Attribute, if provided.
		var colors []color.Attribute
		if len(args) > 1 {
			if cols, ok := args[1].([]color.Attribute); ok {
				colors = cols
			}
		}
		// Use the dynamic printing function.
		PrintColoredDynamicToStdout(texts, colors)
		return
	}

	// Otherwise, assume the first argument is a string.
	prefix, ok := args[0].(string)
	if !ok {
		// If not, do nothing.
		return
	}

	// Determine the secondary string.
	secondary := ""
	if len(args) >= 2 {
		if sec, ok := args[1].(string); ok {
			secondary = sec
		}
	}

	// Gather any color attributes passed.
	var attrs []color.Attribute
	if len(args) > 2 {
		// Iterate over the remaining arguments.
		for _, arg := range args[2:] {
			// Use reflection to support both a single color or a slice of colors.
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

	// If no colors are provided, default to white.
	if len(attrs) == 0 {
		attrs = append(attrs, color.FgWhite)
	}

	// Print using FprintColored to os.Stdout.
	FprintColored(os.Stdout, prefix, secondary, attrs...)
}
