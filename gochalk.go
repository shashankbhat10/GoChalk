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

func escapedStyles(styles string) string {
	return fmt.Sprintf("%s[%sm", escape, styles)
}

// Method to return given string encapsulated with given styling
func getStyledString(style Style, value string) string {
	return fmt.Sprintf("%s%s%s", escapedStyle(style), value, resetStyle)
}

func getMultipleStyledString(styles string, value string) string {
	return fmt.Sprintf("%s%s%s", escapedStyles(styles), value, resetStyle)
}

func getSingleString(color Style, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}

	var finalString string
	for index, str := range strs {
		if strings.Contains(str, escape) && index != len(strs)-1 {
			finalString += str + escapedStyle(color) + " "
		} else if index == len(strs)-1 {
			finalString += str
		} else {
			finalString += str + " "
		}
	}

	// return getColoredString(color, finalString)
	return getStyledString(color, finalString)
}

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

	// sort.Slice(styles, func(i, j int) bool {
	// 	return int(styles[i]) < int(styles[j])
	// })

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

// Creates a new Chalk object with the provided styles. If multiple foreground or background colors are provided as
// parameters, then the last one will be applied
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
//	colorful := chalk.New(gochalk.FgRed).Add(gochalk.FgYellow, gochalk.FgMagenta) // Multiple colors are provided to Add so magenta chosen and will replace existing red
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
	// for _, style := range styles {
	// 	if (style >= 30 && style < 38) || (style >= 90 || style < 98) {
	// 		lastForeground = style
	// 	} else if (style >= 40 && style < 48) || (style >= 100 || style < 108) {
	// 		lastBackground = style
	// 	} else {
	// 		stylesCopy = append(stylesCopy, style)
	// 	}
	// }
	// if lastForeground != -1 {
	// 	stylesCopy = append(stylesCopy, lastForeground)
	// }

	// if lastBackground != -1 {
	// 	stylesCopy = append(stylesCopy, lastBackground)
	// }

	newChalk := Chalk{}
	newChalk.styles = append(chalk.styles, styles...)
	return &newChalk
}

// Method to remove any present styling from Chalk. If given styling is not present then method will do nothing happens
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

// Method to print value provided wrapped in styles present in Chalk
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

func BrightRed(value ...string) string {
	return getSingleStyledString(FgBrightRed, value...)
}

func BrightGreen(value ...string) string {
	return getSingleStyledString(FgBrightGreen, value...)
}

func BrightYellow(value ...string) string {
	return getSingleStyledString(FgBrightYellow, value...)
}

func BrightBlue(value ...string) string {
	return getSingleStyledString(FgBrightBlue, value...)
}

func BrightMagenta(value ...string) string {
	return getSingleStyledString(FgBrightRed, value...)
}

func BrightCyan(value ...string) string {
	return getSingleStyledString(FgBrightCyan, value...)
}

func BrightWhite(value ...string) string {
	return getSingleStyledString(FgBrightWhite, value...)
}

// ----------------------------
// Methods for colored backgrounds
// ----------------------------

func RedBg(value ...string) string {
	return getSingleStyledString(BgRed, value...)
}

func GreenBg(value ...string) string {
	return getSingleStyledString(BgGreen, value...)
}

func YellowBg(value ...string) string {
	return getSingleStyledString(BgYellow, value...)
}

func BlueBg(value ...string) string {
	return getSingleStyledString(BgBlue, value...)
}

func MagentaBg(value ...string) string {
	return getSingleStyledString(BgMagenta, value...)
}

func CyanBg(value ...string) string {
	return getSingleStyledString(BgCyan, value...)
}

func WhiteBg(value ...string) string {
	return getSingleStyledString(BgWhite, value...)
}

func BrightRedBg(value ...string) string {
	return getSingleStyledString(BgBrightRed, value...)
}

func BrightGreenBg(value ...string) string {
	return getSingleStyledString(BgBrightGreen, value...)
}

func BrightYellowBg(value ...string) string {
	return getSingleStyledString(BgBrightYellow, value...)
}

func BrightBlueBg(value ...string) string {
	return getSingleStyledString(BgBrightBlue, value...)
}

func BrightMagentaBg(value ...string) string {
	return getSingleStyledString(BgBrightMagenta, value...)
}

func BrightCyanBg(value ...string) string {
	return getSingleStyledString(BgBrightCyan, value...)
}

func BrightWhiteBg(value ...string) string {
	return getSingleStyledString(BgBrightWhite, value...)
}

// ----------------------------
// Methods for applying text styling
// ----------------------------

func TextBold(value ...string) string {
	return getSingleStyledString(Bold, value...)
}

func TextDim(value ...string) string {
	return getSingleStyledString(Dim, value...)
}

func TextItalics(value ...string) string {
	return getSingleStyledString(Italics, value...)
}

func TextUnderlined(value ...string) string {
	return getSingleStyledString(Underlined, value...)
}

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
