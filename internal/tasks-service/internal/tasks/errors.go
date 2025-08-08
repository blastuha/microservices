package tasks

import "fmt"

var ErrTaskNotFound = fmt.Errorf("task not found")
var ErrInvalidInput = fmt.Errorf("task has no title")
