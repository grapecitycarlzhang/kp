package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
)

func calculateMemUsageUnixNoCache(mem types.MemoryStats) float64 {
	return float64(mem.Usage - mem.Stats["cache"])
}

func calculateMemPercentUnixNoCache(limit float64, usedNoCache float64) float64 {
	if limit != 0 {
		return usedNoCache / limit * 100.0
	}
	return 0
}
func memUsage(mem float64, limit float64) string {
	return fmt.Sprintf("%s / %s", bytesSize(mem), bytesSize(limit))
}
func memPerc(per float64) string {
	return fmt.Sprintf("%.2f%%", per)
}
