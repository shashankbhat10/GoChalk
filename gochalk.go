package gochalk

import (
	"fmt"
	"slices"
	"strings"
)

const escape = "\x1b"
const resetStyle = escape + "[0m"

type Style int

// Foreground color (font colors)
const (
	FgBlack Style = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground color (font color) - Bright / Hi-intensity
const (
	FgBrightBlack Style = iota + 90
	FgBrightRed
	FgBrightGreen
	FgBrightYellow
	FgBrightBlue
	FgBrightMagenta
	FgBrightCyan
	FgBrightWhite
)

// Background Color
const (
	BgBlack Style = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Color - Bright / Hi-intensity
const (
	BgBrightBlack Style = iota + 100
	BgBrightRed
	BgBrightGreen
	BgBrightYellow
	BgBrightBlue
	BgBrightMagenta
	BgBrightCyan
	BgBrightWhite
)

// Text display style
const (
	Bold Style = iota + 1
	Dim
	Italics
	Underlined
)

// Method to return styles in escaped string format
func escapedStyle(style Style) string {
	return fmt.Sprintf("%s[%dm", escape, style)
}

// Method to return multiple styles in escaped string format. Use when a single style needs to be applied
func escapedStyles(styles string) string {
	return fmt.Sprintf("%s[%sm", escape, styles)
}

// Method to return given string encapsulated with given styling
func getStyledString(style Style, value string) string {
	return fmt.Sprintf("%s%s%s", escapedStyle(style), value, resetStyle)
}

// Method to return string wrapped in style. Use this when multiple styles might be applied
func getMultipleStyledString(styles string, value string) string {
	return fmt.Sprintf("%s%s%s", escapedStyles(styles), value, resetStyle)
}

// func getSingleString(color Style, strs ...string) string {
// 	if len(strs) == 0 {
// 		return ""
// 	}

// 	var finalString string
// 	for index, str := range strs {
// 		if strings.Contains(str, escape) && index != len(strs)-1 {
// 			finalString += str + escapedStyle(color) + " "
// 		} else if index == len(strs)-1 {
// 			finalString += str
// 		} else {
// 			finalString += str + " "
// 		}
// 	}

// 	// return getColoredString(color, finalString)
// 	return getStyledString(color, finalString)
// }

// Method to apply one or more styles to a string. If no style argument is provided then given string is returned as is.
//
// If multiple foreground / background styles are provided then the last corresponding foreground / background will be applied.
//
// examples:
//
//	gochalk.StyledString("Hello World", gochalk.FgRed) // Returns string wrapped in red foreground style
//	gochalk.StyledString("Hello World", gochalk.FgRed, gochalk.FgYellow, gochalk.FgGreen) // only green foreground will be applied
//	gochalk.StyledString("Hello World", gochalk.FgRed, gochalk.Bold) // Returns string with bold and red styles
func StyledString(val string, styles ...Style) string {
	if len(styles) == 0 {
		return val
	}

	lastForeground := getLastForeground(styles...)
	lastBackground := getLastBackground(styles...)

	stylesCopy := slices.Clone(styles)
	if lastForeground != -1 {
		stylesCopy = addorReplaceForeground(lastForeground, stylesCopy...)
	}
	if lastBackground != -1 {
		stylesCopy = addorReplaceBackground(lastBackground, stylesCopy...)
	}

	slices.Sort(stylesCopy)

	finalStyle := ""

	for _, str := range styles {
		finalStyle += fmt.Sprint(str) + ";"
	}

	finalStyle = finalStyle[0 : len(finalStyle)-1]
	styledString := ""

	if strings.HasSuffix(val, "\n") {
		stringWithNoNewLine := removeNewLine(val)
		styledString += getMultipleStyledString(finalStyle, stringWithNoNewLine) + "\n"
	}
	return styledString
}

type Chalk struct {
	styles []Style
}

// Creates a new Chalk object with the provided styles. This object can then be reused to apply required styles to strings
// If multiple foreground or background colors are provided as parameters, then the last one will be applied
//
//	error := gochalk.NewStyle(gochalk.FgRed) // Returns a Chalk object with red foreground style applied
func NewStyle(styles ...Style) *Chalk {
	if len(styles) == 0 {
		return &Chalk{}
	}

	newStyles := &Chalk{}
	return newStyles.Add(styles...)
}

// Method to add a style to current chalk object. If no parameter given then nothing happens and same object is returned.
// If Foreground or Background color is provided, it will replace any corresponding foreground / background style
// If multiple Foreground / Background styles are provided, then only the last corresponding color will be applied
//
//	errorChalk := chalk.New(gochalk.FgRed)
//	errorBold := errorChalk.Add(gochalk.Bold) // Adds Bold style to 'error' chalk object
//	warningBold := errorBold.Add(gochalk.FgYellow) // Foreground color red is replaced by yellow
//	colorful := chalk.New(gochalk.FgRed).Add(gochalk.FgYellow, gochalk.FgMagenta) // Multiple colors are provided to Add so magenta is chosen and will replace existing red
func (chalk *Chalk) Add(styles ...Style) *Chalk {
	if len(styles) == 0 {
		return chalk
	}

	lastForeground := getLastForeground(styles...)
	lastBackground := getLastBackground(styles...)

	stylesCopy := slices.Clone(chalk.styles)
	if lastForeground != -1 {
		stylesCopy = addorReplaceForeground(lastForeground, stylesCopy...)
	}
	if lastBackground != -1 {
		stylesCopy = addorReplaceBackground(lastBackground, stylesCopy...)
	}

	slices.Sort(stylesCopy)

	newChalk := Chalk{}
	newChalk.styles = append(chalk.styles, styles...)
	return &newChalk
}

// Method to remove any present styling from Chalk. If given styling is not present then method will do nothing happens.
// Method will return a new Chalk
//
//	errorBold := gochalk.NewStyle(gochalk.FgRed, gochalk.Bold)
//	errorNormal := gochalk.Remove(gochalk.Bold) // Bold styling will be removed. Will not effect errorBold object as new Chalk is returned
func (chalk *Chalk) Remove(styles ...Style) *Chalk {
	if len(styles) == 0 {
		return chalk
	}

	stylesFiltered := filterSlice(chalk.styles, styles)

	newChalk := &Chalk{styles: stylesFiltered}

	return newChalk
}

// Method to remove all styles applied to Chalk
func (chalk *Chalk) RemoveAll() *Chalk {
	chalk = nil
	return &Chalk{}
}

// Method to print string provided wrapped in styles present in Chalk
func (chalk *Chalk) Println(value ...string) {
	fmt.Println(chalk.ToString(value...))
}

// Method to get string with formatted style
func (chalk *Chalk) ToString(value ...string) string {
	combinedStyleString := convertIntSliceToString(chalk.styles)

	if len(value) == 0 {
		return ""
	}

	combinedValue := combineStrings(value...)

	// return getStyledString(combinedStyleString, combinedValue)
	return getMultipleStyledString(combinedStyleString, combinedValue)
}

// -------------------------
// Utility Methods
// -------------------------

// Method to convert int slice to a single string.
// Used to create a single string with ';' delimeter for using in styles
func convertIntSliceToString(arr []Style) string {
	return strings.Trim(strings.Replace(fmt.Sprint(arr), " ", ";", -1), "[]")
}

// Method to join all variadic string params
func combineStrings(strs ...string) string {
	if len(strs) == 0 {
		return ""
	}

	return strings.Join(strs, " ")
}

// Method to remove trailing newline symbol
func removeNewLine(value string) string {
	str, _ := strings.CutSuffix(strings.Trim(value, " "), "\n")

	return str
}

// Method to obtain a single string wrapped by the required style
func getSingleStyledString(style Style, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}

	var finalString string
	for index, str := range strs {
		if strings.Contains(str, escape) && index != len(strs)-1 {
			finalString += str + escapedStyle(style) + " "
		} else {
			finalString += str + " "
		}
	}

	return getStyledString(style, finalString[0:len(finalString)-1])
}

// Method to filter source slice by removing items present in filter slice
func filterSlice(source []Style, remove []Style) []Style {
	if len(remove) == 0 {
		return source
	}

	var sliceCopy []Style

	for _, item := range remove {
		if !slices.Contains(source, item) {
			sliceCopy = append(sliceCopy, item)
		}
	}

	return sliceCopy
}

// Method to get the last provided foreground if present from variadic styles argument
func getLastForeground(styles ...Style) Style {
	for index := range styles {
		reverseIndex := len(styles) - index - 1
		style := styles[reverseIndex]
		if (style >= 30 && style < 38) || (style >= 90 || style < 98) {
			return style
		}
	}

	return -1
}

// Method to get last background style is present from variadic styles argument
func getLastBackground(styles ...Style) Style {
	for index := range styles {
		reverseIndex := len(styles) - index - 1
		style := styles[reverseIndex]
		if (style >= 40 && style < 48) || (style >= 100 || style < 108) {
			return style
		}
	}

	return -1
}

// Method to add or replace any existing foreground style with provided color argument
func addorReplaceForeground(color Style, styles ...Style) []Style {
	replaced := false
	var stylesCopy []Style
	for index := range styles {
		style := styles[index]
		if (style >= 30 && style < 38) || (style >= 90 || style < 98) {
			if !replaced {
				// styles[index] = color
				stylesCopy = append(stylesCopy, color)
				replaced = true
			}
		} else {
			stylesCopy = append(stylesCopy, style)
		}
	}

	if !replaced {
		stylesCopy = append(stylesCopy, color)
	}

	return stylesCopy
}

// Method to add or replace any existing background style with provided color argument
func addorReplaceBackground(color Style, styles ...Style) []Style {
	replaced := false
	var stylesCopy []Style
	for index := range styles {
		style := styles[index]
		if (style >= 40 && style < 48) || (style >= 100 || style < 108) {
			if !replaced {
				// styles[index] = color
				stylesCopy = append(stylesCopy, color)
				replaced = true
			}
		} else {
			stylesCopy = append(stylesCopy, style)
		}
	}

	if !replaced {
		stylesCopy = append(stylesCopy, color)
	}

	return stylesCopy
}

// ----------------------------
// Methods for colored strings
// ----------------------------

// Method to print string with red foreground color
func Red(value ...string) string {
	return getSingleStyledString(FgRed, value...)
}

// Method to print string with green foreground color
func Green(value ...string) string {
	return getSingleStyledString(FgGreen, value...)
}

// Method to print string with yellow foreground color
func Yellow(value ...string) string {
	return getSingleStyledString(FgYellow, value...)
}

// Method to print string with blue foreground color
func Blue(value ...string) string {
	return getSingleStyledString(FgBlue, value...)
}

// Method to print string with magenta foreground color
func Magenta(value ...string) string {
	return getSingleStyledString(FgMagenta, value...)
}

// Method to print string with cyan foreground color
func Cyan(value ...string) string {
	return getSingleStyledString(FgCyan, value...)
}

// Method to print string with white foreground color
func White(value ...string) string {
	return getSingleStyledString(FgWhite, value...)
}

// Method to print string with bright red foreground color
func BrightRed(value ...string) string {
	return getSingleStyledString(FgBrightRed, value...)
}

// Method to print string with bright green foreground color
func BrightGreen(value ...string) string {
	return getSingleStyledString(FgBrightGreen, value...)
}

// Method to print string with bright yellow foreground color
func BrightYellow(value ...string) string {
	return getSingleStyledString(FgBrightYellow, value...)
}

// Method to print string with bright blue foreground color
func BrightBlue(value ...string) string {
	return getSingleStyledString(FgBrightBlue, value...)
}

// Method to print string with bright magenta foreground color
func BrightMagenta(value ...string) string {
	return getSingleStyledString(FgBrightRed, value...)
}

// Method to print string with bright cyan foreground color
func BrightCyan(value ...string) string {
	return getSingleStyledString(FgBrightCyan, value...)
}

// Method to print string with bight white foreground color
func BrightWhite(value ...string) string {
	return getSingleStyledString(FgBrightWhite, value...)
}

// ----------------------------
// Methods for colored backgrounds
// ----------------------------

// Method to print string with red background color
func RedBg(value ...string) string {
	return getSingleStyledString(BgRed, value...)
}

// Method to print string with green background color
func GreenBg(value ...string) string {
	return getSingleStyledString(BgGreen, value...)
}

// Method to print string with yellow background color
func YellowBg(value ...string) string {
	return getSingleStyledString(BgYellow, value...)
}

// Method to print string with blue background color
func BlueBg(value ...string) string {
	return getSingleStyledString(BgBlue, value...)
}

// Method to print string with magenta background color
func MagentaBg(value ...string) string {
	return getSingleStyledString(BgMagenta, value...)
}

// Method to print string with cyan background color
func CyanBg(value ...string) string {
	return getSingleStyledString(BgCyan, value...)
}

// Method to print string with white background color
func WhiteBg(value ...string) string {
	return getSingleStyledString(BgWhite, value...)
}

// Method to print string with bright red background color
func BrightRedBg(value ...string) string {
	return getSingleStyledString(BgBrightRed, value...)
}

// Method to print string with bright green background color
func BrightGreenBg(value ...string) string {
	return getSingleStyledString(BgBrightGreen, value...)
}

// Method to print string with bright yellow background color
func BrightYellowBg(value ...string) string {
	return getSingleStyledString(BgBrightYellow, value...)
}

// Method to print string with bright blue background color
func BrightBlueBg(value ...string) string {
	return getSingleStyledString(BgBrightBlue, value...)
}

// Method to print string with bright magenta background color
func BrightMagentaBg(value ...string) string {
	return getSingleStyledString(BgBrightMagenta, value...)
}

// Method to print string with bright cyan background color
func BrightCyanBg(value ...string) string {
	return getSingleStyledString(BgBrightCyan, value...)
}

// Method to print string with bright white background color
func BrightWhiteBg(value ...string) string {
	return getSingleStyledString(BgBrightWhite, value...)
}

// ----------------------------
// Methods for applying text styling
// ----------------------------

// Method to print string with bold styling
func TextBold(value ...string) string {
	return getSingleStyledString(Bold, value...)
}

// Method to print string with dim styling
func TextDim(value ...string) string {
	return getSingleStyledString(Dim, value...)
}

// Method to print string with italics styling
func TextItalics(value ...string) string {
	return getSingleStyledString(Italics, value...)
}

// Method to print string with underlined styling
func TextUnderlined(value ...string) string {
	return getSingleStyledString(Underlined, value...)
}
