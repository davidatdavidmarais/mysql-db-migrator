package parser

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
)

func IsValid(file string) bool {
	s := make([]*string, 0)
	lines := strings.Split(file, "\n")

	for _, l := range lines {
		cl := strings.TrimSpace(l)

		switch cl {
		case endMarker:
			p := peek(s)

			if p != startMarker {
				return false
			}

			s = s[:len(s)-1]
		case startMarker:
			m := startMarker
			s = append(s, &m)
		}
	}

	return len(s) == 0
}

func peek(s []*string) string {
	if len(s) == 0 {
		return ""
	}
	return *s[len(s)-1]
}

func popJoiner(s []*string) string {
	q := ""
	for i := len(s) - 1; i >= 0; i-- {
		t := s[i]
		if *t == startMarker {
			s = s[:i]
			return q
		}

		q = *s[i] + "\n" + q
	}

	return q
}

func Parse(file string) ([]Query, error) {
	v := IsValid(file)
	if !v {
		return nil, errors.New("invalid file")
	}

	s := make([]*string, 0)
	ll := strings.Split(file, "\n")
	ql := make([]Query, 0)

	for _, l := range ll {
		cl := strings.TrimSpace(l)

		switch cl {
		case endMarker:
			q := popJoiner(s)

			nq := strings.TrimSpace(q)

			sum := sha256.Sum256([]byte(nq))
			hash := fmt.Sprintf("%x", sum)

			ql = append(ql, Query{Query: nq, Hash: hash})
		default:
			s = append(s, &cl)
		}
	}

	return ql, nil
}
