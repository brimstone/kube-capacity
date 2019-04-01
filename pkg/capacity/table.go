// Package capacity - text.go contains all the messy details for the text printer implementation
package capacity

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/fatih/color"
)

type tablePrinter struct {
	cm       *clusterMetric
	showPods bool
	showUtil bool
	w        *tabwriter.Writer
}

const (
	formatBase    = "%-25s %25s %25s %25s %25s\n"
	formatPod     = "%-25s %-25s %-50s %25s %25s %25s %25s\n"
	formatUtil    = "%-25s %25s %25s %25s %25s %25s %25s\n"
	formatPodUtil = "%-25s %-25s %-50s %25s %25s %25s %25s %25s %25s\n"
)

func (tp tablePrinter) Print() {
	tp.w.Init(os.Stdout, 0, 8, 2, ' ', 0)
	names := make([]string, len(tp.cm.nodeMetrics))

	i := 0
	for name := range tp.cm.nodeMetrics {
		names[i] = name
		i++
	}
	sort.Strings(names)

	tp.printHeaders()

	for _, name := range names {
		tp.printNode(name, tp.cm.nodeMetrics[name])
	}

	tp.w.Flush()
}

func (tp *tablePrinter) printHeaders() {
	if tp.showPods && tp.showUtil {
		fmt.Fprintf(tp.w, formatPodUtil, "NODE",
			"NAMESPACE",
			"POD",
			color.WhiteString("CPU REQUESTS"),
			color.WhiteString("CPU LIMITS"),
			color.WhiteString("CPU UTIL"),
			color.WhiteString("MEMORY REQUESTS"),
			color.WhiteString("MEMORY LIMITS"),
			color.WhiteString("MEMORY UTIL"))

		if len(tp.cm.nodeMetrics) > 1 {
			fmt.Fprintf(tp.w, formatPodUtil, "*", "*", "*",
				tp.cm.cpu.requestString(),
				tp.cm.cpu.limitString(),
				tp.cm.cpu.utilString(),
				tp.cm.memory.requestString(),
				tp.cm.memory.limitString(),
				tp.cm.memory.utilString())

			fmt.Fprintln(tp.w, "\t\t\t\t\t\t\t\t")
		}
	} else if tp.showPods {
		fmt.Fprintf(tp.w, formatPod, "NODE", "NAMESPACE", "POD",
			color.WhiteString("CPU REQUESTS"),
			color.WhiteString("CPU LIMITS"),
			color.WhiteString("MEMORY REQUESTS"),
			color.WhiteString("MEMORY LIMITS"))

		fmt.Fprintf(tp.w, formatPod, "*", "*", "*",
			tp.cm.cpu.requestString(),
			tp.cm.cpu.limitString(),
			tp.cm.memory.requestString(),
			tp.cm.memory.limitString())

		fmt.Fprintln(tp.w, "\t\t\t\t\t\t")

	} else if tp.showUtil {
		fmt.Fprintf(tp.w, formatUtil, "NODE",
			color.WhiteString("CPU REQUESTS"),
			color.WhiteString("CPU LIMITS"),
			color.WhiteString("CPU UTIL"),
			color.WhiteString("MEMORY REQUESTS"),
			color.WhiteString("MEMORY LIMITS"),
			color.WhiteString("MEMORY UTIL"))

		fmt.Fprintf(tp.w, formatUtil, "*",
			tp.cm.cpu.requestString(),
			tp.cm.cpu.limitString(),
			tp.cm.cpu.utilString(),
			tp.cm.memory.requestString(),
			tp.cm.memory.limitString(),
			tp.cm.memory.utilString())

	} else {
		fmt.Fprintf(tp.w, formatBase, "NODE",
			color.WhiteString("CPU REQUESTS"),
			color.WhiteString("CPU LIMITS"),
			color.WhiteString("MEMORY REQUESTS"),
			color.WhiteString("MEMORY LIMITS"))

		if len(tp.cm.nodeMetrics) > 1 {
			fmt.Fprintf(tp.w, formatBase, "*",
				tp.cm.cpu.requestString(), tp.cm.cpu.limitString(),
				tp.cm.memory.requestString(), tp.cm.memory.limitString())
		}
	}
}

func (tp *tablePrinter) printNode(name string, nm *nodeMetric) {
	podNames := make([]string, len(nm.podMetrics))

	i := 0
	for name := range nm.podMetrics {
		podNames[i] = name
		i++
	}
	sort.Strings(podNames)

	if tp.showPods && tp.showUtil {
		fmt.Fprintf(tp.w, formatPodUtil,
			name, "*", "*",
			nm.cpu.requestString(),
			nm.cpu.limitString(),
			nm.cpu.utilString(),
			nm.memory.requestString(),
			nm.memory.limitString(),
			nm.memory.utilString())

		for _, podName := range podNames {
			pm := nm.podMetrics[podName]
			fmt.Fprintf(tp.w, formatPodUtil,
				name,
				pm.namespace,
				pm.name,
				pm.cpu.requestString(),
				pm.cpu.limitString(),
				pm.cpu.utilString(),
				pm.memory.requestString(),
				pm.memory.limitString(),
				pm.memory.utilString())
		}

		fmt.Fprintln(tp.w, "\t\t\t\t\t\t\t\t")

	} else if tp.showPods {
		fmt.Fprintf(tp.w, formatPod,
			name, "*", "*",
			nm.cpu.requestString(),
			nm.cpu.limitString(),
			nm.memory.requestString(),
			nm.memory.limitString())

		for _, podName := range podNames {
			pm := nm.podMetrics[podName]
			fmt.Fprintf(tp.w, formatPod,
				name,
				pm.namespace,
				pm.name,
				pm.cpu.requestString(),
				pm.cpu.limitString(),
				pm.memory.requestString(),
				pm.memory.limitString())
		}

		fmt.Fprintln(tp.w, "\t\t\t\t\t\t")

	} else if tp.showUtil {
		fmt.Fprintf(tp.w, formatUtil,
			name,
			nm.cpu.requestString(),
			nm.cpu.limitString(),
			nm.cpu.utilString(),
			nm.memory.requestString(),
			nm.memory.limitString(),
			nm.memory.utilString())

	} else {
		fmt.Fprintf(tp.w, formatBase, name,
			nm.cpu.requestString(), nm.cpu.limitString(),
			nm.memory.requestString(), nm.memory.limitString())
	}
}
