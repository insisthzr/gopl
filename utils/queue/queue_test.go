package queue

import (
	"reflect"
	"testing"
)

func TestMain(t *testing.T) {
	want := []int{}
	got := []int{}

	start := 1
	end := 10
	q := NewQueue()
	for i := start; i < end; i++ {
		q.Push(i)
		want = append(want, i)
	}

	for {
		value, err := q.Pop()
		if err == QueueEmpty {
			break
		}
		if err != nil {
			t.Fatal(err.Error())
		}
		got = append(got, value.(int))
	}

	for i := start; i < end; i++ {
		q.Push(i)
		want = append(want, i)
	}

	for {
		value, err := q.Pop()
		if err == QueueEmpty {
			break
		}
		if err != nil {
			t.Fatal(err.Error())
		}
		got = append(got, value.(int))
	}

	if ok := reflect.DeepEqual(want, got); !ok {
		t.Fatalf("want: %v go: %v\n", want, got)
	}
}
