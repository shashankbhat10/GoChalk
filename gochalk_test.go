package gochalk

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
)

const testString string = "Test String"

var stylesSlice = []struct {
	code   Style
	style  string
	method string
}{
	{style: "Black text", code: FgBlack, method: Black(testString)},
	{style: "Red text", code: FgRed, method: Red(testString)},
	{style: "Green text", code: FgGreen, method: Green(testString)},
	{style: "Yellow text", code: FgYellow, method: Yellow(testString)},
	{style: "Blue text", code: FgBlue, method: Blue(testString)},
	{style: "Magenta text", code: FgMagenta, method: Magenta(testString)},
	{style: "Cyan text", code: FgCyan, method: Cyan(testString)},
	{style: "White text", code: FgWhite, method: White(testString)},
	{style: "Bright black text", code: FgBrightBlack, method: BrightBlack(testString)},
	{style: "Bright red text", code: FgBrightRed, method: BrightRed(testString)},
	{style: "Bright green text", code: FgBrightGreen, method: BrightGreen(testString)},
	{style: "Bright yellow text", code: FgBrightYellow, method: BrightYellow(testString)},
	{style: "Bright blue text", code: FgBrightBlue, method: BrightBlue(testString)},
	{style: "Bright magenta text", code: FgBrightMagenta, method: BrightMagenta(testString)},
	{style: "Bright cyan text", code: FgBrightCyan, method: BrightCyan(testString)},
	{style: "Bright white text", code: FgBrightWhite, method: BrightWhite(testString)},
	{style: "Black background", code: BgBlack, method: BlackBg(testString)},
	{style: "Red background", code: BgRed, method: RedBg(testString)},
	{style: "Green background", code: BgGreen, method: GreenBg(testString)},
	{style: "Yellow background", code: BgYellow, method: YellowBg(testString)},
	{style: "Blue background", code: BgBlue, method: BlueBg(testString)},
	{style: "Magenta background", code: BgMagenta, method: MagentaBg(testString)},
	{style: "Cyan background", code: BgCyan, method: CyanBg(testString)},
	{style: "White background", code: BgWhite, method: WhiteBg(testString)},
	{style: "Bright black background", code: BgBrightBlack, method: BrightBlackBg(testString)},
	{style: "Bright red background", code: BgBrightRed, method: BrightRedBg(testString)},
	{style: "Bright green background", code: BgBrightGreen, method: BrightGreenBg(testString)},
	{style: "Bright yellow background", code: BgBrightYellow, method: BrightYellowBg(testString)},
	{style: "Bright blue background", code: BgBrightBlue, method: BrightBlueBg(testString)},
	{style: "Bright magenta background", code: BgBrightMagenta, method: BrightMagentaBg(testString)},
	{style: "Bright cyan background", code: BgBrightCyan, method: BrightCyanBg(testString)},
	{style: "Bright white background", code: BgBrightWhite, method: BrightWhiteBg(testString)},
	{style: "Bold text", code: Bold, method: TextBold(testString)},
	{style: "Dim text", code: Dim, method: TextDim(testString)},
	{style: "Italics text", code: Italics, method: TextItalics(testString)},
	{style: "Underlined text", code: Underlined, method: TextUnderlined(testString)},
}

func TestIndividualStyle(t *testing.T) {
	testString := "Test String"
	for _, style := range stylesSlice {
		actualString := getSingleStyledString(style.code, testString)
		expectedString := fmt.Sprintf("%s[%dm%s%s", escape, int(style.code), testString, resetStyle)

		if actualString != expectedString {
			t.Errorf("Expected: %s\nActual: %s'", expectedString, actualString)
		}
	}
}
func TestEscapedStyle(t *testing.T) {
	style := FgRed

	actualString := escapedStyle(style)
	expectedString := fmt.Sprintf("%s[%dm", escape, FgRed)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestEscapedStyles(t *testing.T) {
	styleSlice := []Style{FgRed, Bold, BgWhite}
	style := convertIntSliceToString(styleSlice)

	actualString := escapedStyles(style)
	expectedString := fmt.Sprintf("%s[%sm", escape, style)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestGetStyledString(t *testing.T) {
	actualString := getStyledString(FgRed, testString)
	expectedString := fmt.Sprintf("%s[%dm%s%s", escape, FgRed, testString, resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestGetMultipleStyledString(t *testing.T) {
	styleSlice := []Style{FgRed, Bold, BgWhite}
	style := convertIntSliceToString(styleSlice)

	actualString := getMultipleStyledString(style, testString)
	expectedString := fmt.Sprintf("%s[%d;%d;%dm%s%s", escape, FgRed, Bold, BgWhite, testString, resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestStyledString_Styles(t *testing.T) {
	styles := []Style{Bold, Underlined, FgCyan, BgWhite}

	actualString := StyledString(testString, styles...)
	expectedString := fmt.Sprintf("%s%s%s", escapedStyles(convertIntSliceToString(styles)), testString, resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestStyledString_Duplicate(t *testing.T) {
	styles := []Style{Bold, Underlined, FgCyan, BgWhite, Bold}

	actualString := StyledString(testString, styles...)
	expectedString := fmt.Sprintf("%s[%sm%s%s", escape, "1;4;36;47", testString, resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestStyledString_Empty(t *testing.T) {
	actualString := StyledString(testString)
	expectedString := fmt.Sprintf("%s", testString)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestStyledString_NewLine(t *testing.T) {
	str := testString + "\n"
	actualString := StyledString(str, Bold, Bold, FgRed, BgWhite)
	expectedString := fmt.Sprintf("%s[%sm%s%s\n", escape, "1;31;47", testString, resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestStyledString_MultipleFgAndBg(t *testing.T) {
	styles := []Style{Bold, Underlined, FgCyan, BgWhite, FgYellow, BgBlue}

	actualString := StyledString(testString, styles...)
	expectedString := fmt.Sprintf("%s[%sm%s%s", escape, "1;4;33;44", testString, resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestNewStyle_EmptyObject(t *testing.T) {
	newChalk := NewStyle()

	testString := "Test String"
	actualString := newChalk.ToString(testString)
	expectedString := testString

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("Expected string with no formatting. \nActual: %s\nExpected: %s", actualString, expectedString)
	}
}

func TestNewStyle_WithStyle(t *testing.T) {
	newChalk := NewStyle(FgMagenta, BgBrightBlue)

	testString := "Test String"
	actualString := newChalk.ToString(testString)
	expectedString := fmt.Sprintf("%s[%sm%s%s", escape, "35;104", testString, resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("Expected string with no formatting. \nActual: %s\nExpected: %s", actualString, expectedString)
	}
}

func TestAdd(t *testing.T) {
	newChalk := NewStyle()
	styleAdded := newChalk.Add(FgRed)

	testString := "Test String"
	styledString := styleAdded.ToString(testString)
	expectedString := fmt.Sprintf("%s[%dm%s%s", escape, int(FgRed), testString, resetStyle)

	if styledString != expectedString {
		t.Errorf("Add method did not add correct style. Actual: %s\t, Expected: %s", styledString, expectedString)
	} else if slices.Compare(newChalk.styles, styleAdded.styles) == 0 {
		t.Errorf("\nExpected: Previous chalk styles should'nt be modified\nActual: Previous chalk styles were be modified")
	}
}

func TestAdd_EmptyArgument(t *testing.T) {
	newChalk := NewStyle(FgRed, Bold)
	styleAdded := newChalk.Add()

	testString := "Test String"
	actualString := styleAdded.ToString(testString)
	expectedString := fmt.Sprintf("%s[%sm%s%s", escape, "1;31", testString, resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("Add method did not add correct style. Actual: %s\t, Expected: %s", actualString, expectedString)
	}
	if !slices.Contains(newChalk.styles, FgRed) || !slices.Contains(newChalk.styles, Bold) {
		t.Errorf("\nExpected: Previous chalk styles should'nt be modified\nActual: Previous chalk styles were be modified")
	}
	if !slices.Contains(styleAdded.styles, FgRed) || !slices.Contains(styleAdded.styles, Bold) {
		t.Errorf("\nExpected: New chalk styles should be same\nActual: New chalk styles is not same")
	}
}

func TestAdd_ReplaceFgAndBg(t *testing.T) {
	newChalk := NewStyle(FgRed, Bold, BgWhite)
	styleAdded := newChalk.Add(FgCyan, BgGreen)

	testString := "Test String"
	actualString := styleAdded.ToString(testString)
	expectedString := fmt.Sprintf("%s[%sm%s%s", escape, "1;36;42", testString, resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("Add method did not add correct style. Actual: %s\t, Expected: %s", actualString, expectedString)
	}
	if !slices.Contains(newChalk.styles, FgRed) || !slices.Contains(newChalk.styles, Bold) || !slices.Contains(newChalk.styles, BgWhite) {
		t.Errorf("\nExpected: Previous chalk styles should'nt be modified\nActual: Previous chalk styles were be modified")
	}
	if !slices.Contains(styleAdded.styles, FgCyan) || !slices.Contains(styleAdded.styles, BgGreen) {
		t.Errorf("\nExpected: New chalk styles should be same\nActual: New chalk styles is not same")
	}
}

func TestRemove(t *testing.T) {
	newChalk := NewStyle(FgRed, Bold)
	onlyRed := newChalk.Remove(Bold)

	if slices.Contains(onlyRed.styles, Bold) {
		t.Errorf("\nExpected: Remove should remove specified style\nActual: Remove did not remove specified style")
	}
	if slices.Compare(newChalk.styles, onlyRed.styles) == 0 {
		t.Errorf("\nExpected: Previous chalk styles should'nt be modified\nActual: Previous chalk styles were be modified")
	}
	if !slices.Contains(onlyRed.styles, FgRed) {
		t.Error("\nExpected: Non specified style should remain.\nActual: Non-specified style was removed")
	}
}

func TestRemove_Empty(t *testing.T) {
	newChalk := NewStyle(FgRed, Bold)
	onlyRed := newChalk.Remove()

	if !slices.Contains(newChalk.styles, FgRed) || !slices.Contains(newChalk.styles, Bold) {
		t.Error("\nExpected: Old styles should not be modified\nActual: Old styles were modified")
	}
	if !slices.Contains(onlyRed.styles, FgRed) || !slices.Contains(onlyRed.styles, Bold) {
		t.Error("\nExpected: New style object should have same styles\nActual: New style object does not have same styles")
	}
}

func TestRemoveAll(t *testing.T) {
	newChalk := NewStyle(FgRed, Bold, BgWhite)
	allRemoved := newChalk.RemoveAll()

	if len(allRemoved.styles) != 0 {
		t.Error("\nExpected: All styles should be removed.\nActual: All styled were not removed")
	} else if len(newChalk.styles) != 3 {
		t.Error("\nExpected: Old chalk object should be unchanged.\nActual: Old chalk object was changed")
	}
}

func TestPrintln(t *testing.T) {
	sc := bufio.NewScanner(os.Stdin)

	NewStyle(FgRed, Bold, BgWhite).Println(testString)
	for sc.Scan() {
		actualString := sc.Text()
		expectedString := fmt.Sprintf("%s[%sm;%s%s", escape, "1;31;47", testString, resetStyle)

		if strings.Compare(expectedString, actualString) != 0 {
			t.Errorf("\nExpected: %s\nActual:%s", expectedString, actualString)
		}
	}
}

func TestToString(t *testing.T) {
	newChalk := NewStyle(FgRed, Bold, BgWhite)

	testString := "Test string"
	actualString := newChalk.ToString(testString)
	expectedString := fmt.Sprintf("%s[%sm%s%s", escape, "1;31;47", testString, resetStyle)

	if strings.Compare(expectedString, actualString) != 0 {
		t.Errorf("\nExpected:%s\nActual:%s", expectedString, actualString)
	}
}

func TestToString_NoArgument(t *testing.T) {
	newChalk := NewStyle(FgRed, Bold, BgWhite)

	actualString := newChalk.ToString()
	expectedString := ""

	if strings.Compare(expectedString, actualString) != 0 {
		t.Errorf("\nExpected:%s\nActual:%s", expectedString, actualString)
	}
}

func TestToString_NoStyle(t *testing.T) {
	newChalk := NewStyle()

	actualString := newChalk.ToString(testString)

	if strings.Compare(testString, actualString) != 0 {
		t.Errorf("\nExpected:%s\nActual:%s", testString, actualString)
	}
}

// Test Utility Methods
func TestConvertIntSliceToString(t *testing.T) {
	slice := []Style{FgRed, Bold, BgWhite, Underlined, FgCyan}

	actualString := convertIntSliceToString(slice)
	expectedString := "31;1;47;4;36"

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestAddOrReplaceForeground_AddFg(t *testing.T) {
	prevStyles := []Style{BgRed, Bold, Underlined}
	added := addorReplaceForeground(FgWhite, prevStyles...)

	if !slices.Contains(added, FgWhite) {
		t.Error("\nExpeceted: Style should be added.\nActual: Style was not added")
	}
	if !slices.Contains(added, BgRed) {
		t.Error("\nExpeceted: Background should not be replaced.\nActual: Background was replaced")
	}
}

func TestAddOrReplaceForeground_ReplaceFg(t *testing.T) {
	prevStyles := []Style{FgRed, Bold, Underlined, BgWhite}
	added := addorReplaceForeground(FgCyan, prevStyles...)

	if !slices.Contains(added, FgCyan) {
		t.Error("\nExpeceted: Style should be added.\nActual: Style was not added")
	}
	if slices.Contains(added, FgRed) {
		t.Error("\nExpected: Foreground should be replaced.\nActual: Foreground was not replaced")
	}
	if !slices.Contains(added, BgWhite) {
		t.Error("\nExpected: Background should not be replaced.\nActual: Background was replaced")
	}
}

func TestAddOrReplaceForeground_AddBg(t *testing.T) {
	prevStyles := []Style{FgRed, Bold, Underlined}
	added := addorReplaceBackground(BgWhite, prevStyles...)

	if !slices.Contains(added, BgWhite) {
		t.Error("\nExpeceted: Style should be added.\nActual: Style was not added")
	}
	if !slices.Contains(added, FgRed) {
		t.Error("\nExpeceted: Foreground should not be replaced.\nActual: Foreground was replaced")
	}
}

func TestAddOrReplaceForeground_Replace(t *testing.T) {
	prevStyles := []Style{FgRed, Bold, Underlined, BgWhite}
	added := addorReplaceBackground(BgYellow, prevStyles...)

	if !slices.Contains(added, FgRed) {
		t.Error("\nExpeceted: Foreground was replaced.\nActual: Foreground should not be replaced")
	}
	if slices.Contains(added, BgWhite) {
		t.Error("\nExpected: Background should be replaced.\nActual: Background was not replaced")
	}
	if !slices.Contains(added, BgYellow) {
		t.Error("\nExpected: Background should be replaced.\nActual: Background was not replaced")
	}
}

func TestGetLastForeground(t *testing.T) {
	styles := []Style{FgRed, Bold, Underlined, FgWhite, BgRed}

	lastFg := getLastForeground(styles...)

	if lastFg != FgWhite {
		t.Errorf("\nExpected: %d.\nActual: %d", FgWhite, lastFg)
	}
}

func TestGetLastForeground_NoFg(t *testing.T) {
	styles := []Style{Bold, Underlined, BgRed}

	lastFg := getLastForeground(styles...)

	if lastFg != -1 {
		t.Errorf("\nExpected: %d.\nActual: %d", -1, lastFg)
	}
}

func TestGetLastBackground(t *testing.T) {
	styles := []Style{BgRed, BgBrightGreen, Bold, BgCyan, Underlined, FgWhite}

	lastBg := getLastBackground(styles...)

	if lastBg != BgCyan {
		t.Errorf("\nExpected: %d.\nActual: %d", BgCyan, lastBg)
	}
}

func TestGetLastBackground_NoBg(t *testing.T) {
	styles := []Style{Bold, Underlined, FgWhite}

	lastBg := getLastBackground(styles...)

	if lastBg != -1 {
		t.Errorf("\nExpected: %d.\nActual: %d", -1, lastBg)
	}
}

func TestFilterSlice_Filter(t *testing.T) {
	source := []Style{Bold, FgRed, BgBrightCyan, Underlined}
	remove := []Style{BgBrightCyan, Bold}

	result := filterSlice(source, remove)

	if slices.Contains(result, BgBrightCyan) || slices.Contains(result, Bold) {
		t.Error("\nExpected: Styles in remove slice should have been removed from source.\nActual: Styles in remove were not removed from source")
	}
}

func TestFilterSlice_EmptySlice(t *testing.T) {
	source := []Style{Bold, FgRed, BgBrightCyan, Underlined}
	remove := []Style{}

	result := filterSlice(source, remove)

	if !slices.Contains(result, BgBrightCyan) || !slices.Contains(result, Bold) || !slices.Contains(result, FgRed) || !slices.Contains(result, Underlined) {
		t.Error("\nExpected: Styles in remove slice should have been removed from source.\nActual: Styles in remove were not removed from source")
	}
}

func TestGetSingleStyledString_MultipleString(t *testing.T) {
	strs := []string{"This", "is", "test", "string"}

	style := FgRed

	actualString := getSingleStyledString(style, strs...)
	expectedString := fmt.Sprintf("%s[%dm%s%s", escape, FgRed, "This is test string", resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestGetSingleStyledString_SingleString(t *testing.T) {
	style := FgRed

	actualString := getSingleStyledString(style, "This is test string")
	expectedString := fmt.Sprintf("%s[%dm%s%s", escape, FgRed, "This is test string", resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestGetSingleStyledString_StringContainingStyleInside(t *testing.T) {
	middleString := Green("Green")
	style := FgRed

	actualString := getSingleStyledString(style, "This is test", middleString, "string")
	expectedString := fmt.Sprintf("%s[%dm%s%s[%dm%s%s%s[%dm%s%s", escape, FgRed, "This is test ", escape, FgGreen, "Green", resetStyle, escape, FgRed, " string", resetStyle)

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestGetSingleStyledString_Empty(t *testing.T) {
	style := FgRed

	actualString := getSingleStyledString(style)
	expectedString := fmt.Sprintf("")

	if strings.Compare(actualString, expectedString) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
	}
}

func TestRemoveNewLine(t *testing.T) {
	str := "Test String\n"

	result := removeNewLine(str)

	if strings.Compare(result, "Test String") != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", "Test String", result)
	}
}

func TestCombineStrings_Empty(t *testing.T) {
	result := combineStrings()

	if strings.Compare(result, "") != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", "", result)
	}
}

func TestCombineStrings_Combine(t *testing.T) {
	strs := []string{"Test", "String"}
	result := combineStrings(strs...)

	if strings.Compare(result, "Test String") != 0 {
		t.Errorf("\nExpected: %s\nActual: %s", "", result)
	}
}

// Test Basic Single Styling Methods

func TestBasicStyle(t *testing.T) {
	for index, style := range stylesSlice {
		if index > 10 {
			break
		}
		actualString := style.method
		expectedString := fmt.Sprintf("%s[%dm%s%s", escape, style.code, testString, resetStyle)

		if strings.Compare(actualString, expectedString) != 0 {
			t.Errorf("\nExpected: %s\nActual: %s", expectedString, actualString)
		}
	}
}
