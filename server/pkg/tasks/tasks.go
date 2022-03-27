package tasks

import (
	"math/rand"
	"time"
)

//TODO implemet to ordinary database
type Service interface {
	Get() *Task
}

type Tasks struct {
	tasks []*Task
}

type Task struct {
	Description string
	Variants    []string
	Answer      string
}

// New todo dynamic create tasks
func New() Service {
	return &Tasks{
		tasks: []*Task{
			{Description: "Search max number: ", Variants: []string{"1", "3", "0.5", "abc", "-1"}, Answer: "3"},
			{Description: "Search min number: ", Variants: []string{"1", "3", "0.5", "abc", "-1"}, Answer: "-1"},
			{Description: "Search max odd number: ", Variants: []string{"2", "3", "0.5", "abc", "-1"}, Answer: "3"},
			{Description: "Search max even number: ", Variants: []string{"2", "3", "0.5", "abc", "-1"}, Answer: "2"},
			{Description: "Search min odd number: ", Variants: []string{"2", "3", "0.5", "abc", "-1"}, Answer: "-1"},
			{Description: "Search min even number: ", Variants: []string{"2", "3", "0.5", "4", "-1"}, Answer: "2"},
			{Description: "Search area of a square with a side 1: ", Variants: []string{"1", "3", "0.5", "abc", "-1"}, Answer: "1"},
			{Description: "Search area of a square with a side 2: ", Variants: []string{"1", "3", "4", "abc", "-1"}, Answer: "4"},
			{Description: "Search area of a square with a side 3: ", Variants: []string{"1", "9", "0.5", "abc", "-1"}, Answer: "9"},
			{Description: "Search area of a square with a side 3: ", Variants: []string{"1", "9", "0.5", "abc", "-1"}, Answer: "9"},
		},
	}
}

func (t *Tasks) Get() *Task {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	return t.tasks[r.Intn(len(t.tasks))]
}
