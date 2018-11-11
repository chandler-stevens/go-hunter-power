package main

// Import necessary libraries for file input, output, and string manipulation
import (
"bufio"
"fmt"
"log"
"os"
"strings"
)

// Globally-defined legal character strings
const letter = "abcdefghijklmnopqrstuvwxyz"
const digit = "0123456789"
const symbol = " \""

func main() {
	// Check if input filename given as passed argument from prompt
	if len(os.Args) < 2 {
		log.Fatal("No Input File given!\n")
	}
	filename := os.Args[1]
	// Check if given filenamen actually opens an input file for reading
	inputFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Cannot find Input File!\n", err)
	}
	// Report initial processing of input filename in console
	fmt.Println("Processing Input File:", filename)
	// Check if an output file can be created for writing
	// with same name as input file
	periodIndex := strings.IndexAny(filename, ".")
	outputFilename := (filename[:periodIndex])+".out"
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal("Cannot create Output File!\n", err)
	}
	// Initialize local variables for reading, counting, and error checking
	line := ""
	tokenCount := 0
	spaceIndex := -1
	nextSpaceIndex := -1
	lineIndex := 0
	errorFound := false
	// Setup buffers for reading from input file and writing to output file
	scan := bufio.NewScanner(inputFile)
	write := bufio.NewWriter(outputFile)
	// Read each line from the input file buffer
  for scan.Scan() {
    line = scan.Text()
		// Find the first space in the line
    spaceIndex = strings.IndexAny(line, " ")
		nextSpaceIndex = -1
		// Reset lineIndex to first character in the line
		lineIndex = 0
		// Iterate through each character in the line
		for lineIndex < len(line) {
			// Check if 5 consecutive characters match keyword BEGIN
			if ((lineIndex+4) < len(line)) &&
			(line[lineIndex:(lineIndex+5)] == "BEGIN") &&
			// and nothing besides a newline, space, or tab follow
			(((lineIndex+5) == len(line)) ||
			(((lineIndex+5) < len(line)) &&
			(line[(lineIndex+5)] == ' ') || (line[(lineIndex+5)] == '\t'))) {
				fmt.Fprintf(outputFile, "BEGIN\n")
				tokenCount++
				// Increment character to index after the 'N'
				lineIndex += 5
			// Check if 5 consecutive characters match keyword WRITE
			} else if ((lineIndex+4) < len(line)) &&
				(line[lineIndex:(lineIndex+5)] == "WRITE") &&
			// and nothing besides a newline, space, or tab follow
				(((lineIndex+5) == len(line)) ||
				(((lineIndex+5) < len(line)) &&
				(line[(lineIndex+5)] == ' ') || (line[(lineIndex+5)] == '\t'))) {
					fmt.Fprintf(outputFile, "WRITE\n")
					tokenCount++
					// Increment character to index after the 'E'
					lineIndex += 5
			// Check if 3 consecutive characters match keyword END
			} else if ((lineIndex+2) < len(line)) &&
				(line[lineIndex:(lineIndex+3)] == "END") &&
			// and nothing besides a newline, space, or tab follow
				(((lineIndex+3) == len(line)) ||
				(((lineIndex+3) < len(line)) &&
				(line[(lineIndex+3)] == ' ') || (line[(lineIndex+3)] == '\t'))) {
					fmt.Fprintf(outputFile, "END\n")
					tokenCount++
					// Increment character to index after the 'D'
					lineIndex += 3
			// Check if character matches the terminal symbol . {point}
			} else if lineIndex < len(line) &&
			 	line[lineIndex] == '.' &&
				// and nothing besides a newline, space, or tab follow
				(((lineIndex+1) == len(line)) ||
				(((lineIndex+1) < len(line)) &&
				(line[(lineIndex+1)] == ' ') || (line[(lineIndex+1)] == '\t')))	{
					fmt.Fprintf(outputFile, "POINT\n")
					tokenCount++
					// Increment character to index after the '.'
					lineIndex += 1
			// Check if character matches the terminal symbol $ {string datatype}
			} else if lineIndex < len(line) && line[lineIndex] == '$' {
					// Attempt to find a space after the '$'
					nextSpaceIndex = strings.IndexAny(line[lineIndex:], " ")
					// Check if a space was found after the '$'
					if nextSpaceIndex != -1 {
						// Find the first non-space or non-tab in the line
						firstCharIndex := -1
						index := len(line)-1
						for index >= 0 {
							if line[index] != ' ' || line[index] != '\t' {
								firstCharIndex = index
							}
							index--
						}
						// Check if the character is not the first non-space or non-tab
						if firstCharIndex != -1  && firstCharIndex != lineIndex {
							// Adjust the index of the found space to correct position
							nextSpaceIndex += spaceIndex+1
						}
						// Set identifier name starting after the '$' until the found space
						ident := line[(lineIndex+1):nextSpaceIndex]
						// Check if the identifier name has only legal terminal symbol
						// letters and/or digits
						if identifierCheck(ident, outputFile) {
							fmt.Fprintf(outputFile, ("ID[STRING]: "+ident+"\n"))
							tokenCount++
							// Increment character to index after the found space
							lineIndex = nextSpaceIndex+1
							// Increment space to found space
							spaceIndex = nextSpaceIndex
						} else {
								errorFound = true
						}
					} else {
							// Attempt to find a tab after the '$'
							ident := ""
							tabIndex := strings.IndexAny(line[lineIndex:], "\t")
							// Check if a tab was found after the '$'
							if tabIndex != -1 {
								// Set identifier name starting after '$' until the found tab
								ident = line[(lineIndex+1):tabIndex]
								// Increment character to index after the tab
								lineIndex = tabIndex+1
							} else {
									// Set identifier name starting after '$' until end of line
									ident = line[(lineIndex+1):]
									lineIndex = len(line)
							}
							// Check if the identifier name has only legal terminal symbol
							// letters and/or digits
							if identifierCheck(ident, outputFile) {
								fmt.Fprintf(outputFile, ("ID[STRING]: "+ident+"\n"))
								tokenCount++
							} else {
									errorFound = true
							}
					}
			// Check if character matches the terminal symbol # {integer datatype}
			} else if lineIndex < len(line) && line[lineIndex] == '#' {
					// Attempt to find a space after the '#'
					nextSpaceIndex = strings.IndexAny(line[lineIndex:], " ")
					// Check if a space was found after the '#'
					if nextSpaceIndex != -1 {
						// Find the first non-space or non-tab in the line
						firstCharIndex := -1
						index := len(line)-1
						for index >= 0 {
							if line[index] != ' ' || line[index] != '\t' {
								firstCharIndex = index
							}
							index--
						}
						// Check if the character is not the first non-space or non-tab
						if firstCharIndex != -1  && firstCharIndex != lineIndex {
							// Adjust the index of the found space to correct position
							nextSpaceIndex += spaceIndex+1
						}
						// Set identifier name starting after the '#' until the found space
						ident := line[(lineIndex+1):nextSpaceIndex]
						// Check if the identifier name has only legal terminal symbol
						// letters and/or digits
						if identifierCheck(ident, outputFile) {
							fmt.Fprintf(outputFile, ("ID[INT]: "+ident+"\n"))
							tokenCount++
							// Increment character to index after the found space
							lineIndex = nextSpaceIndex+1
							// Increment space to found space
							spaceIndex = nextSpaceIndex
						} else {
								errorFound = true
						}
					} else {
						// Attempt to find a tab after the '#'
							ident := ""
							tabIndex := strings.IndexAny(line[lineIndex:], "\t")
							// Check if a tab was found after the '#'
							if tabIndex != -1 {
								// Set identifier name starting after '#' until the found tab
								ident = line[(lineIndex+1):tabIndex]
								// Increment character to index after the tab
								lineIndex = tabIndex+1
							} else {
									// Set identifier name starting after '#' until end of line
									ident = line[(lineIndex+1):]
									lineIndex = len(line)
							}
							// Check if the identifier name has only legal terminal symbol
							// letters and/or digits
							if identifierCheck(ident, outputFile) {
								fmt.Fprintf(outputFile, ("ID[INT]: "+ident+"\n"))
								tokenCount++
							} else {
									errorFound = true
							}
					}
			// Check if character matches the terminal symbol % {real number datatype}
			} else if lineIndex < len(line) && line[lineIndex] == '%' {
					// Attempt to find a space after the '%'
					nextSpaceIndex = strings.IndexAny(line[lineIndex:], " ")
					// Check if a space was found after the '%'
					if nextSpaceIndex != -1 {
						// Find the first non-space or non-tab in the line
						firstCharIndex := -1
						index := len(line)-1
						for index >= 0 {
							if line[index] != ' ' || line[index] != '\t' {
								firstCharIndex = index
							}
							index--
						}
						// Check if the character is not the first non-space or non-tab
						if firstCharIndex != -1  && firstCharIndex != lineIndex {
							// Adjust the index of the found space to correct position
							nextSpaceIndex += spaceIndex+1
						}
						// Set identifier name starting after the '%' until the found space
						ident := line[(lineIndex+1):nextSpaceIndex]
						// Check if the identifier name has only legal terminal symbol
						// letters and/or digits
						if identifierCheck(ident, outputFile) {
							fmt.Fprintf(outputFile, ("ID[REAL]: "+ident+"\n"))
							tokenCount++
							// Increment character to index after the found space
							lineIndex = nextSpaceIndex+1
							// Increment space to found space
							spaceIndex = nextSpaceIndex
						} else {
								errorFound = true
						}
					} else {
							// Attempt to find a tab after the '%'
							ident := ""
							tabIndex := strings.IndexAny(line[lineIndex:], "\t")
							// Check if a tab was found after the '%'
							if tabIndex != -1 {
								// Set identifier name starting after '#' until the found tab
								ident = line[(lineIndex+1):tabIndex]
								// Increment character to index after the tab
								lineIndex = tabIndex+1
							} else {
									// Set identifier name starting after '%' until end of line
									ident = line[(lineIndex+1):]
									lineIndex = len(line)
							}
							// Check if the identifier name has only legal terminal symbol
							// letters and/or digits
							if identifierCheck(ident, outputFile) {
								fmt.Fprintf(outputFile, ("ID[REAL]: "+ident+"\n"))
								tokenCount++
							} else {
									errorFound = true
							}
					}
			// Check if character matches the terminal symbol [0-9] {digit}
			} else if lineIndex < len(line) &&
				(strings.Contains(digit, string(line[lineIndex])) ||
				// or if character matches terminal symbol: + {positive sign)
				// or if if character matches terminal symbol: - {negative sign)
				((line[lineIndex] == '+' || line[lineIndex] == '-') &&
				// with at least one digit after the sign
				(lineIndex+1) != len(line) &&
				strings.Contains(digit, string(line[(lineIndex+1)])))) {
					// Attempt to find a space after the character
					nextSpaceIndex = strings.IndexAny(line[lineIndex:], " ")
					// Check if a space was found after the character
					if nextSpaceIndex != -1 {
						// Find the first non-space or non-tab in the line
						firstCharIndex := -1
						index := len(line)-1
						for index >= 0 {
							if line[index] != ' ' || line[index] != '\t' {
								firstCharIndex = index
							}
							index--
						}
						// Check if the character is not the first non-space or non-tab
						if firstCharIndex != -1  && firstCharIndex != lineIndex {
							// Adjust the index of the found space to correct position
							nextSpaceIndex += spaceIndex+1
						}
						// Set number starting from character until the found space
						number := line[lineIndex:nextSpaceIndex]
						// Attempt to find a decimal point in the number
						decimalIndex := strings.IndexAny(number, ".")
						// Check if number has a sign given
						signAdjust := 0
						if line[lineIndex] == '+' || line[lineIndex] == '-' {
							signAdjust = 1
						}
						// Check if a decimal point was found in the number
						if decimalIndex != -1 {
							// Check if the characters before and after the decimal point
							// are only legal terminal symbol digits
							if numberCheck(number[signAdjust:decimalIndex], outputFile) &&
							numberCheck(number[decimalIndex+1:], outputFile) {
								fmt.Fprintf(outputFile, ("REAL_CONST: "+number+"\n"))
								tokenCount++
								// Increment character to index after the found space
								lineIndex = nextSpaceIndex+1
								// Increment space to found space
								spaceIndex = nextSpaceIndex
							} else {
									errorFound = true
							}
						// Check if the number has only legal terminal symbol digits
						} else if numberCheck(number, outputFile) {
								fmt.Fprintf(outputFile, ("INT_CONST: "+number+"\n"))
								tokenCount++
								// Increment character to index after the found space
								lineIndex = nextSpaceIndex+1
								// Increment space to found space
								spaceIndex = nextSpaceIndex
						} else {
								errorFound = true
						}
				} else {
						// Check if number has a sign given
						signAdjust := 0
						if line[lineIndex] == '+' || line[lineIndex] == '-' {
							signAdjust = 1
						}
						number := ""
						// Attempt to find a tab after the character
						tabIndex := strings.IndexAny(line[lineIndex:], "\t")
						// Check if a tab was found after the character
						if tabIndex != -1 {
							// Set number starting from character until the found tab
							number = line[lineIndex:tabIndex]
							// Increment character to index after the tab
							lineIndex = tabIndex+1
						} else {
							// Set number starting from character until end of line
							number = line[lineIndex:]
							lineIndex = len(line)
						}
						// Attempt to find a decimal point in the number
						decimalIndex := strings.IndexAny(number, ".")
						// Check if a decimal point was found in the number
						if decimalIndex != -1 {
							// Check if the characters before and after the decimal point
							// are only legal terminal symbol digits
							if numberCheck(number[signAdjust:decimalIndex], outputFile) &&
							(decimalIndex+1 == len(number) ||
							(decimalIndex+1 < len(number) &&
							numberCheck(number[decimalIndex+1:], outputFile))) {
								fmt.Fprintf(outputFile, ("REAL_CONST: "+number+"\n"))
								tokenCount++
							} else {
									errorFound = true
							}
						// Check if the number has only legal terminal symbol digits
						} else if numberCheck(number, outputFile) {
								fmt.Fprintf(outputFile, ("INT_CONST: "+number+"\n"))
								tokenCount++
						} else {
								errorFound = true
						}
				}
			// Check if character matches the terminal symbol " {double quotation mark}
			} else if lineIndex < len(line) && line[lineIndex] == '"' {
					// Find the last double quotation mark in the line
					nextQuotationIndex := strings.LastIndexAny(line, "\"")
					// Check if there actually is a second double quotation mark
					if nextQuotationIndex != -1 &&
					nextQuotationIndex != lineIndex &&
					// and nothing besides a newline, space, or tab follow
					(((nextQuotationIndex+1) == len(line)) ||
					(((nextQuotationIndex+1) < len(line)) &&
					(line[nextQuotationIndex+1] == ' ') ||
					(line[nextQuotationIndex+1] == '\t'))) {
						// Set phrase starting from first double quotation mark
						// through second double quotation mark
						phrase := line[lineIndex:(nextQuotationIndex+1)]
						// Check if the phrase has only legal terminal symbol
						// letters, digits, spaces, and/or tabs
						if stringCheck(phrase, outputFile) {
							fmt.Fprintf(outputFile, ("STRING: "+phrase+"\n"))
							tokenCount++
						} else {
								errorFound = true
						}
						// Increment character to index after second double quotation mark
						lineIndex = nextQuotationIndex+1
					} else {
							// Check if a space was found after the character
							nextSpaceIndex = strings.IndexAny(line[lineIndex:], " ")
							if nextSpaceIndex != -1 {
								// Find the first non-space or non-tab in the line
								firstCharIndex := -1
								index := len(line)-1
								for index >= 0 {
									if line[index] != ' ' || line[index] != '\t' {
										firstCharIndex = index
									}
									index--
								}
								// Check if the character is not the first non-space or non-tab
								if firstCharIndex != -1  && firstCharIndex != lineIndex {
									// Adjust the index of the found space to correct position
									nextSpaceIndex += spaceIndex+1
								}
								// Increment character to index after the second quotation mark
								lineIndex = nextSpaceIndex+1
							} else {
									// Attempt to find a tab after the character
									tabIndex := strings.IndexAny(line[lineIndex:], "\t")
									// Check if a tab was found after the character
									if tabIndex != -1 {
										// Increment character to index after the tab
										lineIndex = tabIndex+1
									// Report lexical error in output file with failed lexeme
									} else {
										fmt.Fprintf(outputFile,
											"Lexical Error, unrecognized symbol "+
											string(line[lineIndex:])+"\n")
										lineIndex = len(line)
										errorFound = true
									}
							}
					}
			// Check if character matches the terminal symbol <= {assignment operator}
			} else if lineIndex < len(line) && line[lineIndex] == '<' &&
				strings.IndexAny(line, "=") == (lineIndex+1) &&
				// and nothing besides a newline, space, or tab follow
				(((lineIndex+2) == len(line)) ||
				(((lineIndex+2) < len(line)) &&
				(line[lineIndex+2] == ' ') ||
				(line[lineIndex+2] == '\t'))) {
					fmt.Fprintf(outputFile, "ASSIGN\n")
					tokenCount++
					// Increment character to index after assignment operator
					lineIndex += 2
			// Check if character matches the terminal symbol + {addition operator}
			} else if lineIndex < len(line) && line[lineIndex] == '+' &&
				// and nothing besides a newline, space, or tab follow
				(((lineIndex+1) == len(line)) ||
				(((lineIndex+1) < len(line)) &&
				(line[lineIndex+1] == ' ') ||
				(line[lineIndex+1] == '\t'))) {
					fmt.Fprintf(outputFile, "PLUS\n")
					tokenCount++
					// Increment character to index after addition operator
					lineIndex++
			// Check if character matches the terminal symbol - {subtraction operator}
			} else if lineIndex < len(line) && line[lineIndex] == '-' &&
				// and nothing besides a newline, space, or tab follow
				(((lineIndex+1) == len(line)) ||
				(((lineIndex+1) < len(line)) &&
				(line[lineIndex+1] == ' ') ||
				(line[lineIndex+1] == '\t'))) {
					fmt.Fprintf(outputFile, "MINUS\n")
					tokenCount++
					// Increment character to index after subtraction operator
					lineIndex++
			// Check if character matches the
			// terminal symbol * {multiplication operator}
			} else if lineIndex < len(line) && line[lineIndex] == '*' &&
				// and nothing besides a newline, space, or tab follow
				(((lineIndex+1) == len(line)) ||
				(((lineIndex+1) < len(line)) &&
				(line[lineIndex+1] == ' ') ||
				(line[lineIndex+1] == '\t'))) {
					fmt.Fprintf(outputFile, "TIMES\n")
					tokenCount++
					// Increment character to index after multiplication operator
					lineIndex++
			// Check if character matches the terminal symbol / {division operator}
			} else if lineIndex < len(line) && line[lineIndex] == '/' &&
				// and nothing besides a newline, space, or tab follow
				(((lineIndex+1) == len(line)) ||
				(((lineIndex+1) < len(line)) &&
				(line[lineIndex+1] == ' ') ||
				(line[lineIndex+1] == '\t'))) {
					fmt.Fprintf(outputFile, "DIVISION\n")
					tokenCount++
					// Increment character to index after division operator
					lineIndex++
			// Check if character matches the terminal symbol ^ {exponent operator}
			} else if lineIndex < len(line) && line[lineIndex] == '^' &&
				// and nothing besides a newline, space, or tab follow
				(((lineIndex+1) == len(line)) ||
				(((lineIndex+1) < len(line)) &&
				(line[lineIndex+1] == ' ') ||
				(line[lineIndex+1] == '\t'))) {
					fmt.Fprintf(outputFile, "POWER\n")
					tokenCount++
					// Increment character to index after exponent operator
					lineIndex++
			// Check if character matches the terminal symbol ( {opening parenthesis}
			} else if lineIndex < len(line) && line[lineIndex] == '(' &&
				// and nothing besides a newline, space, or tab follow
				(((lineIndex+1) == len(line)) ||
				(((lineIndex+1) < len(line)) &&
				(line[lineIndex+1] == ' ') ||
				(line[lineIndex+1] == '\t'))) {
					fmt.Fprintf(outputFile, "LPAREN\n")
					tokenCount++
					// Increment character to index after opening parenthesis
					lineIndex++
			// Check if character matches the terminal symbol ) {closing parenthesis}
			} else if lineIndex < len(line) && line[lineIndex] == ')' &&
				// and nothing besides a newline, space, or tab follow
				(((lineIndex+1) == len(line)) ||
				(((lineIndex+1) < len(line)) &&
				(line[lineIndex+1] == ' ') ||
				(line[lineIndex+1] == '\t'))) {
					fmt.Fprintf(outputFile, "RPAREN\n")
					tokenCount++
					// Increment character to index after closing parenthesis
					lineIndex++
			// Check if character matches the terminal symbol : {colon}
			} else if lineIndex < len(line) && line[lineIndex] == ':' &&
				// and nothing besides a newline, space, or tab follow
				(((lineIndex+1) == len(line)) ||
				(((lineIndex+1) < len(line)) &&
				(line[lineIndex+1] == ' ') ||
				(line[lineIndex+1] == '\t'))) {
					fmt.Fprintf(outputFile, "COLON\n")
					tokenCount++
					// Increment character to index after colon
					lineIndex++
			// Check if character is a space or a tab
			} else if lineIndex < len(line) &&
				(line[lineIndex] == ' ' || line[lineIndex] == '\t') {
				// Increment space to found space or tab
				spaceIndex = lineIndex
				// Increment next space to found space or tab
				nextSpaceIndex = lineIndex
				// Increment character to index after space or tab
				lineIndex++
			} else if lineIndex < len(line) {
				// Attempt to find a space after the character
				nextSpaceIndex = strings.IndexAny(line[lineIndex:], " ")
				// Check if a space was found after the character
				if nextSpaceIndex != -1 {
					// Find the first non-space or non-tab in the line
					firstCharIndex := -1
					index := len(line)-1
					for index >= 0 {
						if line[index] != ' ' || line[index] != '\t' {
							firstCharIndex = index
						}
						index--
					}
					// Check if the character is not the first non-space or non-tab
					if firstCharIndex != -1  && firstCharIndex != lineIndex {
						// Adjust the index of the found space to correct position
						nextSpaceIndex += spaceIndex+1
					}
					// Report lexical error in output file with failed lexeme
					fmt.Fprintf(outputFile, "Lexical Error, unrecognized symbol "+
						string(line[lineIndex:nextSpaceIndex])+"\n")
					// Increment character to index after the found space
					lineIndex = nextSpaceIndex+1
					errorFound = true
			} else {
					// Attempt to find a tab after the character
					tabIndex := strings.IndexAny(line[lineIndex:], "\t")
					// Check if a tab was found after the character
					if tabIndex != -1 {
						// Report lexical error in output file with failed lexeme
						fmt.Fprintf(outputFile, "Lexical Error, unrecognized symbol "+
							string(line[lineIndex:tabIndex])+"\n")
						// Increment character to index after the found tab
						lineIndex = tabIndex+1
						errorFound = true
					} else {
							// Report lexical error in output file with failed lexeme
							fmt.Fprintf(outputFile, "Lexical Error, unrecognized symbol "+
								string(line[lineIndex:])+"\n")
							errorFound = true
					}
				}
			}
		}
	}
	// Check if the reading buffer failed to read a line from the input file
  if err := scan.Err(); err != nil {
      log.Fatal("Could not read line from Input File!\n", err)
  }
	// Clear the writing buffer
	write.Flush()
	// Close the input file and output file
	inputFile.Close()
	outputFile.Close()
	// Report total number of successful tokens produced in console
	fmt.Println(tokenCount, "Tokens produced")
	// Report output filename in console
	fmt.Println("Results in Output File:", outputFilename)
	// Report lexical error in console
	if errorFound {
		fmt.Println("There was a Lexical error processing the file")
	}
}

// Supplementary function to check a if a phrase
// has only legal terminal symbol letters, digits, spaces, and/or tabs
func stringCheck(phrase string, outputFile *os.File) bool {
	charIndex := 0
	success := true
	for charIndex < len(phrase) {
		if !((strings.Contains(letter, string(phrase[charIndex]))) ||
		(strings.Contains(digit, string(phrase[charIndex]))) ||
		(strings.Contains(symbol, string(phrase[charIndex])))) {
			success = false
		}
		charIndex++
	}
	// Report lexical error in output file with failed lexeme
	if !(success) {
		fmt.Fprintf(outputFile, ("Lexical Error, unrecognized symbol "+
			phrase+"\n"))
	}
	return success
}

// Supplementary function to check a if a number
// has only legal terminal symbol digits
func numberCheck(number string, outputFile *os.File) bool {
	charIndex := 0
	success := true
	for charIndex < len(number) {
		if !(strings.Contains(digit, string(number[charIndex]))) {
			success = false
		}
		charIndex++
	}
	// Report lexical error in output file with failed lexeme
	if !(success) {
		fmt.Fprintf(outputFile, ("Lexical Error, unrecognized symbol "+
			number+"\n"))
	}
	return success
}

// Supplementary function to check a if an identifier
// has only legal terminal symbol letters and/or digits
func identifierCheck(identifier string, outputFile *os.File) bool {
	charIndex := 0
	success := true
	for charIndex < len(identifier) {
		if !((strings.Contains(letter, string(identifier[charIndex]))) ||
		(strings.Contains(digit, string(identifier[charIndex])))) {
			success = false
		}
		charIndex++
	}
	// Report lexical error in output file with failed lexeme
	if !(success) {
		fmt.Fprintf(outputFile, ("Lexical Error, unrecognized symbol "+
			identifier+"\n"))
	}
	return success
}
