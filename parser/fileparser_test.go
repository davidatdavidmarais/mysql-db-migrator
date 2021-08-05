package parser

import (
	"testing"

	"gotest.tools/assert"
)

func TestIsValid(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output bool
	}{
		{
			name:   "Success",
			input:  "-- start\nasdfljk\nasdflkjasdlfj\nasdflkjasdf\n-- end",
			output: true,
		},
		{
			name:  "Failed 1",
			input: "-- start\nasdfljk\nasdflkjasdlfj\nasdflkjasdf\n-- end\n-- start\nasdflkjasdf",
		},
		{
			name:  "Failed 2",
			input: "\nasdfljk\nasdflkjasdlfj\nasdflkjasdf\n-- end\n-- start\nasdflkjasdf\n-- end",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := IsValid(tc.input)
			assert.Equal(t, tc.output, resp)
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output []Query
	}{
		{
			name:  "Success 1",
			input: "-- start\nasdfljk\nasdflkjasdlfj\nasdflkjasdf\n-- end",
			output: []Query{
				{
					Query: "asdfljk\nasdflkjasdlfj\nasdflkjasdf",
					Hash:  "148a752161cd5ae2fd93c0bf0da3952dc7a13c81ab51fdda3c5a1e592db53753",
				},
			},
		},
		{
			name:  "Success 2",
			input: "-- start\nasdfljk\nasdflkjasdlfj\nasdflkjasdf\n-- end\n-- start\nasdfljk\nasdflkjasdlfj\nasdflkjasdf\n-- end",
			output: []Query{
				{
					Query: "asdfljk\nasdflkjasdlfj\nasdflkjasdf",
					Hash:  "148a752161cd5ae2fd93c0bf0da3952dc7a13c81ab51fdda3c5a1e592db53753",
				},
				{
					Query: "asdfljk\nasdflkjasdlfj\nasdflkjasdf",
					Hash:  "148a752161cd5ae2fd93c0bf0da3952dc7a13c81ab51fdda3c5a1e592db53753",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := Parse(tc.input)
			assert.NilError(t, err)
			assert.Equal(t, len(tc.output), len(resp))
			for i, v := range resp {
				assert.Equal(t, tc.output[i], v)
			}
		})
	}
}
