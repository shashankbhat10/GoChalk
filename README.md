<h1 align="center">GoChalk</h1>
<p align="center">A minimalist package to style your terminal output the way you want</p>

## Installation

```bash
go get github.com/shashankbhat10/GoChalk
```

## Usage

### Basic styling

```go
import (
    "github.com/shashankbhat10/gochalk"
)

func main() {
    // Prints a string in red font color
    fmt.Println(gochalk.Red("This is a red string"))

    // Print a string with multiple styles
    fmt.Println(gochalk.StyledString("String with multiple styles", gochalk.FgCyan, gochalk.Bold, gochalk.Underlined))

    // Print a string having mix of styled string
    fmt.Println(gochalk.Green("Green", gochalk.Red("Red String"), "String"))
}
```

### Create Chalk Objects

To reuse any style, create Chalk Objects. Chalk objects are immutable, and new styles being added or removed will return a new chalk object

```go
// Create a Chalk object
chalk := gochalk.NewStyle(gochalk.Bold, gochalk.FgRed)
// Will print a string with styles present in chalk object
fmt.Println(chalk.toString("Bold red string"))

// Will replace red font color with yellow color
boldWarning := chalk.Remove(gochalk.FgRed).Add(gochalk.FgYellow)

// Providing foreground or background colors will replace any existing corresponding color
warning := chalk.Add(gochalk.FgYellow)

// RemoveAll will return a new chalk with all styles removed
normalSuccess := chalk.RemoveAll().Add(gochalk.FgGreen)

// To print string using object, call the Println method of the object
normalSuccess.Println("All tests passed")
```

## Features

- Support for all basic colors supported in terminals
- No additional dependencies
- 100% test coverage

# To Do

Add support for 256 color range in supported terminals

# Author

Shashank Bhat
