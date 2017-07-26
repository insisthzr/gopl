package comma

import (
	"bytes"
	"strings"
)

const (
	jump = 3
)

func Comma(str string) string {
	if str == "" {
		return str
	}
	last := len(str) - 1
	sign := ""
	if str[last] == '+' || str[last] == '-' {
		sign = string(str[last])
		str = str[:last]
	}

	count := strings.Count(str, ".")
	if count > 1 {
		panic("nan")
	}

	left := str
	right := ""
	index := strings.Index(str, ".")
	if index > -1 {
		left = str[:index]
		if index < last {
			right = str[index+1:] //WATCH OUT index DON'T greater equal than slice capcity
		}
	}

	result := sign
	left = comma(left)
	right = comma(right)
	if right == "" {
		result += left
	} else {
		result += left + "." + right
	}
	return result
}

func comma(str string) string {
	buf := &bytes.Buffer{}
	last := len(str) - 1
	for i := last; i >= 0; i-- {
		diff := last - i
		if diff%jump == 0 && i != last {
			buf.WriteByte(',')
		}
		buf.WriteByte(str[i])
	}

	result := buf.Bytes()
	i := 0
	j := len(result) - 1
	for i < j {
		result[i], result[j] = result[j], result[i]
		i++
		j--
	}
	return string(result)
}
