//go:build linux

package main

import (
	"testing"

	"github.com/grafana/alloy/internal/cmd/integration-tests/common"
)

func TestBeylaMetrics(t *testing.T) {
	var beylaMetrics = []string{
		"???",
	}
	common.MimirMetricsTest(t, beylaMetrics, []string{}, "beyla")
}
