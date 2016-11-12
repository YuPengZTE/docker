package task

import (
	"fmt"
	"sort"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/command"
)

const (
	psTaskItemFmt = "%s\t%s\t%s\t%s\t%s %s ago\t%s\t%s\n"
	maxErrLength  = 30
)

type portStatus swarm.PortStatus

func (ps portStatus) String() string {
	if len(ps.Ports) == 0 {
		return ""
	}

	str := fmt.Sprintf("*:%d->%d/%s", ps.Ports[0].PublishedPort, ps.Ports[0].TargetPort, ps.Ports[0].Protocol)
	for _, pConfig := range ps.Ports[1:] {
		str += fmt.Sprintf(",*:%d->%d/%s", pConfig.PublishedPort, pConfig.TargetPort, pConfig.Protocol)
	}

	return str
}

type tasksBySlot []swarm.Task

func (t tasksBySlot) Len() int {
	return len(t)
}

func (t tasksBySlot) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t tasksBySlot) Less(i, j int) bool {
	// Sort by slot.
	if t[i].Slot != t[j].Slot {
		return t[i].Slot < t[j].Slot
	}

	// If same slot, sort by most recent.
	return t[j].Meta.CreatedAt.Before(t[i].CreatedAt)
}

// PrintQuiet shows task list in a quiet way.
func PrintQuiet(dockerCli *command.DockerCli, tasks []swarm.Task) error {
	sort.Stable(tasksBySlot(tasks))

	out := dockerCli.Out()

	for _, task := range tasks {
		fmt.Fprintln(out, task.ID)
	}

	return nil
}
