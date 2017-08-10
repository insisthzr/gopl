package eval

import (
	"fmt"
	"math"
)

func wrapParenthesis(str string, ok bool) string {
	if !ok {
		return str
	}
	return fmt.Sprintf("(%s)", str)
}

type Expr interface {
	fmt.Stringer
	Eval(Env) float64
	Check(vars map[Var]bool) error
	isAtom() bool
}

type Var string

func (v Var) Eval(env Env) float64 {
	return env[v]
}

func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

func (v Var) String() string {
	return string(v)
}

func (v Var) isAtom() bool {
	return true
}

type literal float64

func (l literal) Eval(env Env) float64 {
	return float64(l)
}

func (l literal) Check(vars map[Var]bool) error {
	return nil
}

func (l literal) String() string {
	return fmt.Sprintf("%f", l)
}

func (l literal) isAtom() bool {
	return true
}

type unary struct {
	op rune
	x  Expr
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

func (u unary) Check(vars map[Var]bool) error {
	switch u.op {
	case '+', '-':
	default:
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

func (u unary) String() string {
	if u.op == '+' {
		return fmt.Sprintf("%s", u.x)
	}
	x := wrapParenthesis(u.x.String(), !u.x.isAtom())
	return fmt.Sprintf("%c%s", u.op, x)
}

func (u unary) isAtom() bool {
	if u.op == '+' {
		return true
	}
	return false
}

type binary struct {
	op rune
	x  Expr
	y  Expr
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

func (b binary) Check(vars map[Var]bool) error {
	switch b.op {
	case '+', '-', '*', '/':
	default:
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (b binary) String() string {
	x := wrapParenthesis(b.x.String(), !b.x.isAtom())
	y := wrapParenthesis(b.y.String(), !b.y.isAtom())
	return fmt.Sprintf("%s%c%s", x, b.op, y)
}

func (b binary) isAtom() bool {
	return false
}

type call struct {
	fn   string
	args []Expr
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d",
			c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

func (c call) String() string {
	var str string
	switch c.fn {
	case "pow":
		str = fmt.Sprintf("pow(%s,%s)", c.args[0], c.args[1])
	case "sin":
		str = fmt.Sprintf("sin(%s)", c.args[0])
	case "sqrt":
		str = fmt.Sprintf("sqrt(%s)", c.args[0])
	default:
		panic("unkown function")

	}
	return str
}

func (c call) isAtom() bool {
	return true
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

type Env map[Var]float64
