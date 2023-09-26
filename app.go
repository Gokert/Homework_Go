package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Flags struct {
	countStrings   bool
	showDuplicates bool
	haveRegister   bool
	showUnique     bool
	countSymbols   int
	countRow       int
}

type Data struct {
	command    *Flags
	inputFile  io.Reader
	outputFile io.Writer
}

func Uniq(data *Data) error {
	var err error = nil
	var prevValue string
	originPrevValue := ""

	lenghtStrings := 0
	lastCount := 0

	haveRows := (*data.command).countRow
	haveSymbols := (*data.command).countSymbols
	haveRegister := (*data.command).haveRegister

	scanner := bufio.NewScanner((*data).inputFile)
	for scanner.Scan() {
		line := scanner.Text()

		if lenghtStrings > 0 {
			if haveSymbols < len(originPrevValue) {
				prevValue = originPrevValue[haveSymbols:]
			}
		}
		if haveSymbols > 0 && haveSymbols < len(line) {
			line = line[haveSymbols:]
		}
		if haveRegister {
			line = strings.ToLower(line)
			prevValue = strings.ToLower(prevValue)
		}
		if haveRows > 0 {
			words := strings.Split(line, " ")
			if haveRows < len(words) {
				line = strings.Join(words[haveRows:], " ")
			}

			words = strings.Split(prevValue, " ")
			if haveRows < len(words) {
				prevValue = strings.Join(words[haveRows:], " ")
			}
		}

		if lenghtStrings == 0 {
			originPrevValue = scanner.Text()
			lenghtStrings += 1
			lastCount = 1
		} else if prevValue != line {
			if err := Writer(data, originPrevValue, lastCount); err != nil {
				return err
			}

			originPrevValue = scanner.Text()
			lastCount = 1
		} else if prevValue == line {
			lenghtStrings += 1
			lastCount += 1
		}
	}

	if err := Writer(data, originPrevValue, lastCount); err != nil {
		return err
	}

	return err
}

func Writer(data *Data, line string, count int) error {
	var err error = nil
	uniqFlag := (*data.command).showUnique
	сountFlag := (*data.command).countStrings
	duplicateFlag := (*data.command).showDuplicates

	if !сountFlag && !duplicateFlag && !uniqFlag {
		io.WriteString((*data).outputFile, fmt.Sprintf("%s\n", line))
	} else if сountFlag && !duplicateFlag && !uniqFlag {
		io.WriteString(data.outputFile, fmt.Sprintf("%d %s\n", count, line))
	} else if !сountFlag && duplicateFlag && !uniqFlag {
		if count > 1 {
			io.WriteString(data.outputFile, fmt.Sprintf("%s\n", line))
		}
	} else if !сountFlag && !duplicateFlag && uniqFlag {
		if count == 1 {
			io.WriteString(data.outputFile, fmt.Sprintf("%s\n", line))
		}
	}

	return err
}

func argsFlags(data *Data) error {
	var err error = nil

	flag.BoolVar(&(*data.command).showUnique, "u", false, "Show only unique ones")
	flag.BoolVar(&(*data.command).showDuplicates, "d", false, "Show only duplicates")
	flag.BoolVar(&(*data.command).countStrings, "c", false, "Show number of repetitions")
	flag.BoolVar(&(*data.command).haveRegister, "i", false, "Case sensitive")
	flag.IntVar(&(*data.command).countRow, "f", 0, "How many words to ignore")
	flag.IntVar(&(*data.command).countSymbols, "s", 0, "How many words to ignore characters")

	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		if args[0] == "input.txt" {
			input, err := os.Open("input.txt")
			if err != nil {
				return err
			}
			data.inputFile = input
		} else if args[0] == "output.txt" {
			output, err := os.OpenFile("output.txt", os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			data.outputFile = output
		}

	}
	if len(args) > 1 {

		if args[0] == "input.txt" {
			input, err := os.Open("input.txt")
			if err != nil {
				return err
			}
			data.inputFile = input
		}
		if args[1] == "output.txt" {
			output, err := os.OpenFile("output.txt", os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			data.outputFile = output
		}
	}

	if (*data.command).showUnique && (*data.command).showDuplicates || (*data.command).showDuplicates && (*data.command).countStrings || (*data.command).showUnique && (*data.command).countStrings {
		err = fmt.Errorf("You cannot use multiple mutually exclusive flags -d -c -u")
	}

	return err
}

func main() {
	command := Flags{false, false, false, false, 0, 0}

	data := Data{
		command: &command,

		inputFile:  os.Stdin,
		outputFile: os.Stdout,
	}

	if err := argsFlags(&data); err != nil {
		fmt.Println(err)
		return
	}
	if err := Uniq(&data); err != nil {
		fmt.Println(err)
		return
	}

}
