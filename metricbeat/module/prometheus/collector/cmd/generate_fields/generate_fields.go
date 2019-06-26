package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"

	"github.com/elastic/beats/metricbeat/helper/prometheus"
)

type FamiliesByName []*dto.MetricFamily

func (f FamiliesByName) Len() int           { return len(f) }
func (f FamiliesByName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f FamiliesByName) Less(i, j int) bool { return f[i].GetName() < f[j].GetName() }

func filterIncomplete(families []*dto.MetricFamily) {
	k := 0
	for i, f := range families {
		if f.Name == nil {
			continue
		}
		families[k] = families[i]
		k++
	}
}

type FieldsEntry struct {
	Name        string
	Type        string
	Description string
}

func (e FieldsEntry) Print(w io.Writer, spaces int) {
	indent := strings.Repeat("  ", spaces)
	fmt.Fprintf(w, "%s- name: %s\n", indent, e.Name)
	fmt.Fprintf(w, "%s  type: %s\n", indent, e.Type)
	if e.Description != "" {
		fmt.Fprintf(w, "%s  description: >\n%s    %s\n", indent, indent, e.Description)
	}
}

func main() {
	flag.Parse()

	format := expfmt.FmtText
	families, err := prometheus.GetFamilies(os.Stdin, format)
	if err != nil {
		fmt.Printf("failed to get families from stdin: %v", err)
		os.Exit(1)
	}

	filterIncomplete(families)
	sort.Sort(FamiliesByName(families))

	for _, f := range families {
		switch *f.Type {
		case dto.MetricType_SUMMARY, dto.MetricType_HISTOGRAM:
			bucket := FieldsEntry{
				Name:        fmt.Sprintf("prometheus.metrics.%s_bucket", *f.Name),
				Type:        "double",
				Description: *f.Help,
			}
			sum := FieldsEntry{
				Name: fmt.Sprintf("prometheus.metrics.%s_sum", *f.Name),
				Type: "double",
			}
			count := FieldsEntry{
				Name: fmt.Sprintf("prometheus.metrics.%s_count", *f.Name),
				Type: "double",
			}
			if f.Help != nil && len(*f.Help) > 0 {
				sum.Description = "Sum of " + *f.Help
				count.Description = "Count of " + *f.Help
			}
			bucket.Print(os.Stdout, 1)
			sum.Print(os.Stdout, 1)
			count.Print(os.Stdout, 1)
		default:
			entry := FieldsEntry{
				Name:        fmt.Sprintf("prometheus.metrics.%s", *f.Name),
				Type:        "double",
				Description: *f.Help,
			}
			entry.Print(os.Stdout, 1)
		}
	}
}
