package main

import (
	"fmt"

	"github.com/cloudfoundry/gosigar"
	ui "github.com/gizak/termui"
)

func main() {
	strs := GetProcList()

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ls := ui.NewList()
	ls.Items = strs
	ls.ItemFgColor = ui.ColorYellow
	//ls.BorderLabel = "List"
	ls.BorderLabel = "  PID  PPID STIME     TIME    RSS S COMMAND"
	ls.Height = 10
	ls.Width = 250
	ls.Y = 0

	ui.Render(ls)
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()

}

func GetProcList() []string {
	pids := sigar.ProcList{}
	pids.Get()

	strs := []string{}
	// ps -eo pid,ppid,stime,time,rss,state,comm
	//strs = append(strs, fmt.Sprint("  PID  PPID STIME     TIME    RSS S COMMAND"))
	for _, pid := range pids.List {
		state := sigar.ProcState{}
		mem := sigar.ProcMem{}
		time := sigar.ProcTime{}

		if err := state.Get(pid); err != nil {
			continue
		}
		if err := mem.Get(pid); err != nil {
			continue
		}
		if err := time.Get(pid); err != nil {
			continue
		}

		strs = append(strs, fmt.Sprintf("%5d %5d %s %s %6d %c %s",
			pid, state.Ppid,
			time.FormatStartTime(), time.FormatTotal(),
			mem.Resident/1024, state.State, state.Name))
	}
	return strs
}
