package main

import (
	"fmt"
	"regexp"
	"strconv"
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

// too high: 576266
func main() {
	workflows, parts := getData("../data.txt")

	ratingsSum := 0
	for _, p := range parts {
		if p.EvalAccepted(workflows) {
			ratingsSum += p.XMASRatingSum()
		}
	}

	fmt.Println(ratingsSum)
}

type Workflows map[string]*Workflow
type Parts []*Part

type Part struct {
	X int
	M int
	A int
	S int
}

func (p *Part) XMASRatingSum() int {
	return p.X + p.M + p.A + p.S
}

func (p *Part) EvalAccepted(workflows Workflows) bool {
	var next *string
	start := "in"
	next = &start
	for next != nil {
		wf := workflows[*next]
		var accepted bool
		next, accepted, _ = wf.Rules.EvalAccepted(p)
		if next == nil && accepted {
			return true
		}
	}
	return false
}

type Rule string

// EvalAccepted if *string is nil, there is no next, and use bool as accepted, otherwise use *string as next and ignore bool
func (r Rule) EvalAccepted(part *Part) (*string, bool, bool) {
	var next *string
	if !reRuleWithCondition.MatchString(string(r)) && (reAcceptRejectRule.MatchString(string(r)) || reGotoRule.MatchString(string(r))) {
		val := string(r)
		next = &val
	}

	terminated := true

	if reRuleWithCondition.MatchString(string(r)) {
		matches := reRuleWithCondition.FindStringSubmatch(string(r))
		var rating int
		switch matches[1] {
		case "x":
			rating = part.X
		case "m":
			rating = part.M
		case "a":
			rating = part.A
		case "s":
			rating = part.S
		}
		val, _ := strconv.Atoi(matches[3])
		if matches[2] == ">" {
			if rating > val {
				next = &matches[4]
			} else {
				terminated = false
			}
		} else {
			if rating < val {
				next = &matches[4]
			} else {
				terminated = false
			}
		}
	}

	if next != nil {
		if reAcceptRejectRule.MatchString(*next) {
			if *next == "A" {
				return nil, true, terminated
			}
		} else {
			return next, false, terminated
		}
	}

	return nil, false, terminated
}

type Rules []Rule

func (rs Rules) EvalAccepted(part *Part) (*string, bool, bool) {
	for _, r := range rs {
		next, accepted, terminated := r.EvalAccepted(part)
		if terminated {
			if next == nil {
				return nil, accepted, true
			}
			return next, false, true
		}
	}
	return nil, false, true
}

type Workflow struct {
	ID    string
	Rules Rules
}

func getData(filename string) (Workflows, Parts) {
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

	parts := make(Parts, 0, len(lines))
	for _, line := range lines {
		if rePartLine.MatchString(line) {
			matches := rePartLine.FindStringSubmatch(line)
			part := Part{}
			part.X, _ = strconv.Atoi(matches[1])
			part.M, _ = strconv.Atoi(matches[2])
			part.A, _ = strconv.Atoi(matches[3])
			part.S, _ = strconv.Atoi(matches[4])
			parts = append(parts, &part)
		} else {
			panic("invalid line: " + line)
		}
	}

	return workflows, parts
}
