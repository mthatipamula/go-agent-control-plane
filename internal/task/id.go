package task

import (
	"fmt"
	"time"
)

func NewID() string {
	return fmt.Sprintf("task-%d", time.Now().UnixNano())
}
