package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	substitution Substitution
	file         *os.File
	// fileName      string
	printSelected bool
	editInPlace   bool
	onlyRange     bool
}

type Substitution struct {
	command     string
	pattern     *regexp.Regexp
	replacement string
	flag        SubstFlag
	lineRange   []int
}

type SubstFlag struct {
	global                bool
	globalCaseInsensitive bool
	delete                bool
	print                 bool
}

func main() {
	cfg := loadConfig()

	substituteReader(cfg)

	// substitute(cfg)

}

func loadConfig() Config {
	var err error
	var cfg Config

	// output a range of lines from the file. specify a range, i.e.
	// for lines 2 to 4 we would use the command: cat -n ccsed -n '2,4p’ filename
	flag.BoolVar(&cfg.printSelected, "n", false, "only print selected")
	// flag.BoolVar(&cfg.doubleSpacing, "G", false, "double spacing a file")
	flag.BoolVar(&cfg.editInPlace, "i", false, "edit in place")

	flag.Parse()
	args := flag.Args()

	subst := ""

	if len(args) == 0 {
		fmt.Println("Please provide substitution and file as args")
		os.Exit(1)
	} else if len(args) == 1 {
		subst = args[0]
		cfg.file = os.Stdin
	} else {
		subst = args[0]
		cfg.file, err = os.Open(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if subst != "" {
		cfg.substitution, err = parseSubstitution(subst, cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if len(cfg.substitution.lineRange) != 0 && cfg.printSelected {
		cfg.onlyRange = true
	}

	return cfg
}

func parseSubstitution(subst string, cfg Config) (Substitution, error) {
	defaultSubst := "s///g"
	var res Substitution
	var err error
	var pattern bool
	var substList []string

	if subst == "G" {
		subst = "s/\n/\n\n/g"
	}

	for len(substList) < 3 {
		substList = strings.Split(subst, "/")

		if len(substList) == 3 {
			pattern = true
			break
		}
		if len(substList) == 4 {
			break
		}

		if cfg.printSelected {
			res.lineRange, err = parseRangeExpression(subst)
			if err == nil {
				res.flag.print = true
				subst = defaultSubst
			} else {
				err = fmt.Errorf("could not parse substitution: %v: %w", subst, err)
				return res, err
			}
		}
	}

	res.command = substList[0]
	res.pattern, err = regexp.Compile(substList[1])
	if err != nil {
		err = fmt.Errorf("not valid regex pattern: %v", substList[1])
		return res, err
	}
	if !pattern {
		res.replacement = substList[2]
		res.flag, err = parseSubstitutionFlag(substList[3])
		if err != nil {
			err = fmt.Errorf("invalid substitution flag: %v", substList[3])
			return res, err
		}
	} else {
		res.flag, err = parseSubstitutionFlag(substList[2])
		if err != nil {
			err = fmt.Errorf("invalid substitution flag: %v", substList[2])
			return res, err
		}

	}

	// if res.flag.delete {}

	// fmt.Printf("Parsed substitution:\n%v\n%v\n%v\n%v\n", res.command, res.pattern, res.replacement, res.flag.delete)

	return res, nil
}

func parseSubstitutionFlag(flag string) (SubstFlag, error) {
	var res SubstFlag
	switch flag {
	case "g":
		res.global = true
	case "gi":
		res.globalCaseInsensitive = true
	case "d":
		res.delete = true
	case "p":
		res.print = true
	}
	return res, nil
}

/*
func substitute(cfg Config) {
	numLines := 0
	re := cfg.substitution.pattern
	repl := cfg.substitution.replacement
	scanner := bufio.NewScanner(cfg.file)

	for scanner.Scan() {
		numLines++
		fmt.Println(re.ReplaceAllString(scanner.Text(), repl))
	}
}
*/

func substituteReader(cfg Config) {
	numLines := 0
	re := cfg.substitution.pattern
	repl := cfg.substitution.replacement
	reader := bufio.NewReader(cfg.file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading:", err)
			break
		}
		numLines++
		if cfg.onlyRange {
			start, end := cfg.substitution.lineRange[0], cfg.substitution.lineRange[1]
			if numLines >= start && numLines <= end {
				fmt.Printf("%d\t%s", numLines, re.ReplaceAllString(line, repl))
			} else {
				continue
			}
		} else if cfg.substitution.flag.print && cfg.printSelected {
			if re.MatchString(line) {
				fmt.Print(line)
			} else {
				continue
			}
		} else {
			fmt.Print(re.ReplaceAllString(line, repl))
		}
	}
}

func parseRangeExpression(subst string) ([]int, error) {
	var res []int
	splitOnComma := strings.Split(subst, ",")

	if len(splitOnComma) != 2 {
		return res, fmt.Errorf("invalid range expression: %s", subst)
	}

	startRange, err := strconv.Atoi(splitOnComma[0])
	if err != nil {
		return res, err
	}
	res = append(res, startRange)

	secondPart := strings.TrimSuffix(splitOnComma[1], "p")

	endRange, err := strconv.Atoi(secondPart)
	if err != nil {
		return res, err
	}
	res = append(res, endRange)

	return res, nil
}
