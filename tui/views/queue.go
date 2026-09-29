package views

import (
	"fmt"
	"strconv"
	"strings"
	"webcacher2/queue"
)

var Offset int = 15

func RenderQueue() string {
	le := queue.GQueue.Length()
	queue.GQueue.Mtx.Lock()
	defer queue.GQueue.Mtx.Unlock()
	rows := []string{
		TitleStyle.Render(fmt.Sprintf("Firsts [%d] elements in queue", Offset)),
		metric("Cola", strconv.Itoa(le)),
		"",
	}
	f := queue.GQueue.First
	i := 1
	for f != nil {
		rows = append(
			rows,
			fmt.Sprintf("[%s] - [ %s ]", ValueStyle.Render(strconv.Itoa(i)), f.Data.Url),
		)
		f = f.Next
		i += 1
		if i > Offset {
			break
		}
	}

	return PanelStyle.Render(strings.Join(rows, "\n"))
}
