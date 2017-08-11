package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/insisthzr/gopl/ch7/eval"
)

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	fmt.Println("Expression: ")
	stdin.Scan()
	exprStr := stdin.Text()
	fmt.Println("Variables (<var>=<val>, eg: x=3): ")
	stdin.Scan()
	envStr := stdin.Text()
	err := stdin.Err()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	expr, err := eval.Parse(exprStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	env, err := parseEnv(envStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	result := expr.Eval(env)
	fmt.Println(result)
}

func parseEnv(envStr string) (eval.Env, error) {
	env := eval.Env{}
	assignments := strings.Fields(envStr)
	for _, assignment := range assignments {
		fields := strings.Split(assignment, "=")
		if len(fields) != 2 {
			return nil, fmt.Errorf("bad assignment: %s", assignment)
		}
		name := fields[0]
		value, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return nil, err
		}
		env[eval.Var(name)] = value
	}
	return env, nil
}
