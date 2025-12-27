package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mbordner/aoc2023/common/files"
)

var (
	reWorkflowLine      = regexp.MustCompile(`^(\w+)\{(.*)\}$`)
	reRuleWithCondition = regexp.MustCompile(`(x|m|a|s)\s*(>|<)\s*(\d+)\s*:\s*(\w+)`)
	reAcceptRejectRule  = regexp.MustCompile(`(A|R)`)
	reGotoRule          = regexp.MustCompile(`(\w+)`)
	rePartLine          = regexp.MustCompile(`^\{x=(\d+)\s*,\s*m=(\d+)\s*,\s*a=(\d+)\s*,\s*s=(\d+)\}$`)
)

func main() {
	workflows := getData("../data.txt")

	fmt.Println(len(workflows))
}

type Workflows map[string]*Workflow

type Rule string

type Rules []Rule

type Workflow struct {
	ID    string
	Rules Rules
}

func getData(filename string) Workflows {
	lines := files.MustGetLines(filename)

	workflows := make(Workflows)
	for l, line := range lines {
		if line == "" {
			lines = lines[l+1:]
			break
		} else {
			if reWorkflowLine.MatchString(line) {
				matches := reWorkflowLine.FindStringSubmatch(line)
				workflow := Workflow{ID: matches[1]}
				rules := strings.Split(matches[2], ",")
				workflow.Rules = make(Rules, len(rules))

				for r, rule := range rules {
					workflow.Rules[r] = Rule(rule)
				}

				workflows[workflow.ID] = &workflow
			} else {
				panic("invalid line: " + line)
			}
		}
	}

	return workflows
}
