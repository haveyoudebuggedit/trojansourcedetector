package trojansourcedetector

import (
	"fmt"
	"path/filepath"
	"regexp"
)

// match reports whether a name matches a shell file name pattern. The pattern syntax is:
//
//	pattern:
//		{ term }
//	term:
//      '**'        matches any sequence of characters between two path separators, including other separators
//		'*'         matches any sequence of non-Separator characters
//		'?'         matches any single non-Separator character
//		'[' [ '^' ] { character-range } ']'
//		            character class (must be non-empty)
//		c           matches character c (c != '*', '?', '\\', '[')
//		'\\' c      matches character c
//
//	character-range:
//		c           matches character c (c != '\\', '-', ']')
//		'\\' c      matches character c
//		lo '-' hi   matches character c for lo <= c <= hi
func compile(input string) (compiled pattern, err error) {
	output := "^"
	state := stateStart
	for {
		if len(input) == 0 {
			break
		}
		found := false
		for _, rule := range tokenizerRules {
			if rule.inputState != state {
				continue
			}
			matches := rule.pattern.FindSubmatch([]byte(input))
			if len(matches) != 0 {
				state = rule.outputState
				input = input[len(string(matches[0])):]
				output = output + rule.transformer(matches)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("failed to match remaining string '%s' from state '%s'", input, state)
		}
	}
	output = output + "$"
	re, err := regexp.Compile(output)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regexp '%s' from pattern '%s'", output, input)
	}
	return patternImpl{
		re,
	}, nil
}

type tokenizerState string

const (
	stateStart                tokenizerState = "start"
	stateEscape               tokenizerState = "escape"
	stateCharacterClass       tokenizerState = "character_class"
	stateCharacterClassEscape tokenizerState = "character_class_escape"
	stateCharacterClassEnd    tokenizerState = "character_class_end"
)

type tokenizerRule struct {
	inputState  tokenizerState
	pattern     *regexp.Regexp
	outputState tokenizerState
	transformer func(matches [][]byte) string
}

var tokenizerRules = []tokenizerRule{
	{
		stateStart,
		regexp.MustCompile(`^\\`),
		stateEscape,
		func(matches [][]byte) string {
			return ""
		},
	},
	{
		stateEscape,
		regexp.MustCompile(`^.`),
		stateStart,
		func(matches [][]byte) string {
			return regexp.QuoteMeta(string(matches[0]))
		},
	},
	{
		stateStart,
		regexp.MustCompile(`^\*\*(/|$)`),
		stateStart,
		func(matches [][]byte) string {
			return "([^/]+(/|$))*"
		},
	},
	{
		stateStart,
		regexp.MustCompile(`^\*`),
		stateStart,
		func(matches [][]byte) string {
			return "[^/]*"
		},
	},
	{
		stateStart,
		regexp.MustCompile(`^\?`),
		stateStart,
		func(matches [][]byte) string {
			return "[^/]"
		},
	},
	{
		stateStart,
		regexp.MustCompile(`^\[`),
		stateCharacterClass,
		func(matches [][]byte) string {
			return "["
		},
	},
	{
		stateCharacterClass,
		regexp.MustCompile(`^\\`),
		stateCharacterClassEscape,
		func(matches [][]byte) string {
			return "\\"
		},
	},
	{
		stateCharacterClassEscape,
		regexp.MustCompile(`^.`),
		stateCharacterClass,
		func(matches [][]byte) string {
			return string(matches[0])
		},
	},
	{
		stateCharacterClass,
		regexp.MustCompile(`^]`),
		stateCharacterClassEnd,
		func(matches [][]byte) string {
			return "]"
		},
	},
	{
		stateCharacterClass,
		regexp.MustCompile(`^!`),
		stateCharacterClass,
		func(matches [][]byte) string {
			return "^"
		},
	},
	{
		stateCharacterClass,
		regexp.MustCompile(`^.`),
		stateCharacterClass,
		func(matches [][]byte) string {
			return string(matches[0])
		},
	},
	{
		stateCharacterClassEnd,
		regexp.MustCompile(`^\*`),
		stateStart,
		func(matches [][]byte) string {
			return "*"
		},
	},
	{
		stateCharacterClassEnd,
		regexp.MustCompile(`^[^/]+`),
		stateStart,
		func(matches [][]byte) string {
			return regexp.QuoteMeta(string(matches[0]))
		},
	},
	{
		stateCharacterClassEnd,
		regexp.MustCompile(`^/`),
		stateStart,
		func(matches [][]byte) string {
			return regexp.QuoteMeta(string(matches[0]))
		},
	},
	{
		stateStart,
		regexp.MustCompile(`^[^/*?!^]+`),
		stateStart,
		func(matches [][]byte) string {
			return regexp.QuoteMeta(string(matches[0]))
		},
	},
	{
		stateStart,
		regexp.MustCompile(`^/`),
		stateStart,
		func(matches [][]byte) string {
			return regexp.QuoteMeta(string(matches[0]))
		},
	},
}

type patternImpl struct {
	pattern *regexp.Regexp
}

func (p patternImpl) string() string {
	return p.pattern.String()
}

func (p patternImpl) match(name string) (matched bool) {
	return p.pattern.Match([]byte(filepath.ToSlash(name)))
}

type pattern interface {
	match(name string) (matched bool)
	string() string
}
