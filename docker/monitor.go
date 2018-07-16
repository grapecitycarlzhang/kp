package docker

import (
	"encoding/json"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func MonitorStats(cli *client.Client, ctx *context.Context) {
	filters := filters.NewArgs(filters.KeyValuePair{Key: "status", Value: "running"}, filters.KeyValuePair{Key: "label", Value: "alive=restart"})
	containers, err := cli.ContainerList(*ctx, types.ContainerListOptions{Filters: filters})
	if err != nil {
		log.Fatal(err)
	}
	if len(containers) == 0 {
		log.Printf("monitor done for containers len been 0")
	}
	for _, container := range containers {
		go process(cli, ctx, container.ID)
	}
}

func process(cli *client.Client, ctx *context.Context, id string) {
	response, _ := cli.ContainerStats(*ctx, id, false)
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)

	var (
		stats         *types.StatsJSON
		memPercent    float64
		mem, memLimit float64
	)
	dec.Decode(&stats)
	mem = calculateMemUsageUnixNoCache(stats.MemoryStats)
	memLimit = float64(stats.MemoryStats.Limit)
	memPercent = calculateMemPercentUnixNoCache(memLimit, mem)

	log.Println("ID: " + id + " Usage/Limit: " + memUsage(mem, memLimit) + " Per: " + memPerc(memPercent))
	sec, _ := time.ParseDuration("10s")
	log.Printf("restart container %v", id)
	err := cli.ContainerRestart(*ctx, id, &sec)
	if err != nil {
		log.Printf("restart container %v failed", id)
	}
}
