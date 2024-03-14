package parser

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/rodaine/numwords"
)

type Parser interface {
	// Parse reads a transmission from a player and returns an intermediate representation (IR) of the transmitted request.
	// The IR's type must be determined by reflection. The boolean return value is true if the transmission was parsed into a valid IR
	// and false otherwise. If the return value is false, the IR must be nil.
	Parse(string) (any, bool)
}

type parser struct {
	// callsign of the GCI
	callsign string
}

func New() Parser {
	return &parser{}
}

const anyface = "anyface"

type requestWord string

const (
	alphaCheck requestWord = "alpha check"
	bogeyDope  requestWord = "bogey dope"
	declare    requestWord = "declare"
	picture    requestWord = "picture"
	radioCheck requestWord = "radio check"
	spiked     requestWord = "spiked"
	snaplock   requestWord = "snaplock"
)

func requestWords() []requestWord {
	return []requestWord{alphaCheck, bogeyDope, declare, picture, radioCheck, spiked, snaplock}
}

func (p *parser) parseWakeWord(scanner *bufio.Scanner) (string, bool) {
	ok := scanner.Scan()
	if !ok {
		return "", false
	}
	firstSegment := scanner.Text()
	if !(firstSegment == p.callsign || firstSegment == anyface) {
		return "", false
	}
	return firstSegment, true
}

// Parse implements Parser.Parse.
func (p *parser) Parse(tx string) (any, bool) {
	tx = p.sanitize(tx)

	scanner := bufio.NewScanner(strings.NewReader(tx))
	scanner.Split(bufio.ScanWords)

	// Check for a wake word (GCI callsign)
	_, ok := p.parseWakeWord(scanner)
	if !ok {
		return nil, false
	}

	// Scan until we find a request trigger word. Split the scanned tranmission into a callsign segment and a request word.
	var segment string
	callsign := ""
	var rWord requestWord
	for callsign == "" {
		ok := scanner.Scan()
		if !ok {
			return nil, false
		}

		segment = fmt.Sprintf("%s %s", segment, scanner.Text())
		for _, word := range requestWords() {
			if strings.HasSuffix(segment, string(word)) {
				rWord = word
				// Try to parse a callsign from the second segment.
				callsignSegment := strings.TrimSuffix(segment, string(word))
				callsign, ok = parseCallsign(callsignSegment)
				if !ok {
					// TODO send "say again" response?
					return nil, false
				}
				_ = scanner.Scan()

				break
			}
		}
	}

	// Try to parse a request from the remaining text in the scanner.
	switch rWord {
	case alphaCheck:
		// ALPHA CHECK, as implemented by this bot, is a simple request.
		return &alphaCheckRequest{callsign: callsign}, true
	case bogeyDope:
		return p.parseBogeyDope(callsign, scanner)
	case declare:
		return p.parseDeclare(callsign, scanner)
	case picture:
		return p.parsePicture(callsign, scanner)
	case radioCheck:
		// RADIO CHECK is a simple request.
		return &radioCheckRequest{callsign: callsign}, true
	case spiked:
		return p.parseSpiked(callsign, scanner)
	case snaplock:
		return p.parseSnaplock(callsign, scanner)
	}
	return nil, false
}

var sanitizerRex = regexp.MustCompile(`[\p{P}]+`)

// sanitize lowercases the input and replaces all punctuation with spaces.
func (p *parser) sanitize(s string) string {
	s = strings.ToLower(s)
	s = numwords.ParseString(s)
	s = sanitizerRex.ReplaceAllString(s, " ")
	return s
}

var numberWords = map[string]int{
	"0":     0,
	"zero":  0,
	"o":     0,
	"oh":    0,
	"1":     1,
	"one":   1,
	"wun":   1,
	"2":     2,
	"two":   2,
	"3":     3,
	"three": 3,
	"tree":  3,
	"4":     4,
	"four":  4,
	"fower": 4,
	"5":     5,
	"five":  5,
	"fife":  5,
	"6":     6,
	"six":   6,
	"7":     7,
	"seven": 7,
	"8":     8,
	"eight": 8,
	"ait":   8,
	"9":     9,
	"nine":  9,
	"niner": 9,
}

// parseCallsign attempts to parse a callsign in one of the following formats:
//
// - A single word, followed by a number consisting of any digits
//
// - A number consisting of any digits
//
// Garbage in between the digits is ignored. The result is normalize so that each digit is space-delimited.

func parseCallsign(tx string) (callsign string, isValid bool) {
	var scanner = bufio.NewScanner(strings.NewReader(tx))
	scanner.Split(bufio.ScanWords)

	ok := scanner.Scan()
	if !ok {
		return
	}
	firstToken := scanner.Text()
	if firstToken == "" {
		return
	}
	callsign, ok = appendNumber(callsign, firstToken)
	if !ok {
		callsign = firstToken
	} else {
		isValid = true
	}

	for scanner.Scan() {
		nextToken := scanner.Text()
		// Handle single digit
		callsign, ok = appendNumber(callsign, nextToken)
		if ok {
			isValid = true
		} else {
			// Handle case where multiple digits are not space-delimited
			for _, char := range nextToken {
				s, ok := appendNumber(callsign, string(char))
				if ok {
					callsign = s
					isValid = true
				} else {
					return
				}
			}
		}
	}
	return
}

func appendNumber(callsign string, number string) (string, bool) {
	if d, ok := numberWords[number]; ok {
		return fmt.Sprintf("%s %d", callsign, d), true
	}
	return callsign, false
}