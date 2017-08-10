package eval

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEval(t *testing.T) {
	assert := require.New(t)
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
		{"5 / 9 * (F - 32)", Env{"F": 32}, "0"},
		{"5 / 9 * (F - 32)", Env{"F": 212}, "100"},
	}

	for _, test := range tests {
		expr, err := Parse(test.expr)
		assert.NoError(err)
		fmt.Printf("\n%s\n", test.expr)
		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		assert.Equal(test.want, got)

		fmt.Println(expr)
		polishNatation := expr.String()
		expr, err = Parse(polishNatation)
		assert.NoError(err)
		fmt.Printf("\n%s\n", test.expr)
		got = fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		assert.Equal(test.want, got)
	}
}
