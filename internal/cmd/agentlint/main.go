// Command agentlint provides custom linting utilities for the grafana/alloy
// repo.
package main

import (
	"github.com/grafana/alloy/internal/cmd/agentlint/internal/findcomponents"
	"github.com/grafana/alloy/internal/cmd/agentlint/internal/syntaxtags"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		findcomponents.Analyzer,
		syntaxtags.Analyzer,
	)
}
