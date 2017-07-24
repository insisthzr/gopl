package comma

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type commaStruct struct {
	input string
	want  string
}

func TestComma(t *testing.T) {
	assert := require.New(t)
	tests := []*commaStruct{
		{"1", "1"},
		{"12", "12"},
		{"123", "123"},
		{"1234", "1,234"},
		{"123456", "123,456"},
		{"123456789", "123,456,789"},
		{"+1", "+1"},
		{"-12", "-12"},
		{"123", "123"},
		{"1234", "1,234"},
		{"123456", "123,456"},
		{"1234.56789", "1,234.56,789"},
		{".1", ".1"},
		{"1.", "1"},
		{"1", "1"},
	}

	for _, test := range tests {
		got := Comma(test.input)
		assert.Equal(test.want, got)
	}
}
