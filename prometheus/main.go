package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func main() {
	client, err := api.NewClient(api.Config{
		Address: "http://lt4.vscsv.com:64080/prometheus",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, `up{job="node-exporter"}`, time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	switch result.Type() {
	case model.ValVector:
		fmt.Println("Vector Type")
		v, _ := result.(model.Vector)
		displayVector(v)
	default:
		fmt.Printf("Unknow Type")
	}
}

func displayVector(v model.Vector) {
	for _, i := range v {
		fmt.Printf("%s %s\n", praseMetric(i.Metric.String()), i.Value.String())
	}
}

func praseMetric(s string) string {
	res := strings.Split(s, ":")
	res = strings.Split(res[0], `"`)
	return res[1]
}
