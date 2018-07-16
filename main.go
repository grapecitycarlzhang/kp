package main

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	// "strconv"
	"keep/bridge"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/labstack/echo"
	"golang.org/x/net/context"
)

func main() {

	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		if runtime.GOOS != "windows" {
			os.Setenv("DOCKER_HOST", "unix:///tmp/docker.sock")
		} else {
			os.Setenv("DOCKER_HOST", "npipe:////./pipe/docker_engine")
		}
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	b := bridge.New(cli, &ctx)

	b.StartMonitor()

	// e := echo.New()
	// e.GET("/", handlerindex)
	// log.Println("starting echo")
	// err := e.Start(":8080")
	// if err != nil {
	// 	log.Fatal("echo", err)
	// }
}
func handlerindex(c echo.Context) error {
	log.Println("hello world handlerindex")
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		if runtime.GOOS != "windows" {
			os.Setenv("DOCKER_HOST", "unix:///tmp/docker.sock")
		} else {
			os.Setenv("DOCKER_HOST", "npipe:////./pipe/docker_engine")
		}
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	var s string
	for _, container := range containers {
		fmt.Println("==============================")
		fmt.Println(container.Names)
		response, _ := cli.ContainerStats(context.Background(), container.ID, false)
		defer response.Body.Close()
		dec := json.NewDecoder(response.Body)

		var (
			v             *types.StatsJSON
			memPercent    float64
			mem, memLimit float64
		)
		dec.Decode(&v)

		mem = calculateMemUsageUnixNoCache(v.MemoryStats)
		memLimit = float64(v.MemoryStats.Limit)
		memPercent = calculateMemPercentUnixNoCache(memLimit, mem)

		// body, _ := ioutil.ReadAll(response.Body)
		// var stt types.Stats
		// json.Unmarshal(body, &stt)
		// st := stt.MemoryStats
		// fmt.Println("Usage:-----------------")
		// fmt.Println(st.Usage)
		// fmt.Println(float64(st.Usage) / (1024.00 * 1024.00))
		// fmt.Println("MaxUsage:-----------------")
		// fmt.Println(st.MaxUsage)
		// fmt.Println(float64(st.MaxUsage) / (1024.00 * 1024.00))
		// fmt.Println("Limit:-----------------")
		// fmt.Println(st.Limit)
		// fmt.Println(float64(st.Limit) / (1024.00 * 1024.00))
		// fmt.Println("Per:-----------------")
		// fmt.Println(float64(st.Usage) / float64(st.Limit) * 100.00)

		// u := strconv.FormatFloat(float64(st.Usage)/(1024.00*1024.00), 'E', -1, 64)
		// m := strconv.FormatFloat(float64(st.MaxUsage)/(1024.00*1024.00), 'E', -1, 64)
		// l := strconv.FormatFloat(float64(st.Limit)/(1024.00*1024.00), 'E', -1, 64)
		// p := strconv.FormatFloat(float64(st.Usage)/float64(st.Limit)*100.00, 'E', -1, 64)
		mus := MemUsage(mem, memLimit)
		p := MemPerc(memPercent)
		s += " Usage/Limit: " + mus + " Per: " + p
	}

	return c.JSON(http.StatusOK, `{"hello":`+s+`}`)
}
func calculateMemPercent(stat *types.StatsJSON) float64 {
	var memPercent = 0.0
	if stat.MemoryStats.Limit > 0 {
		memPercent = float64(stat.MemoryStats.Usage) / float64(stat.MemoryStats.Limit) * 100.0
	}
	return memPercent
}
func calculateMemUsageUnixNoCache(mem types.MemoryStats) float64 {
	return float64(mem.Usage - mem.Stats["cache"])
}

func calculateMemPercentUnixNoCache(limit float64, usedNoCache float64) float64 {
	// MemoryStats.Limit will never be 0 unless the container is not running and we haven't
	// got any data from cgroup
	if limit != 0 {
		return usedNoCache / limit * 100.0
	}
	return 0
}

var binaryAbbrs = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}

func getSizeAndUnit(size float64, base float64, _map []string) (float64, string) {
	i := 0
	unitsLimit := len(_map) - 1
	for size >= base && i < unitsLimit {
		size = size / base
		i++
	}
	return size, _map[i]
}

// CustomSize returns a human-readable approximation of a size
// using custom format.
func CustomSize(format string, size float64, base float64, _map []string) string {
	size, unit := getSizeAndUnit(size, base, _map)
	return fmt.Sprintf(format, size, unit)
}
func BytesSize(size float64) string {
	return CustomSize("%.4g%s", size, 1024.0, binaryAbbrs)
}
func MemUsage(mem float64, limit float64) string {
	return fmt.Sprintf("%s / %s", BytesSize(mem), BytesSize(limit))
}
func MemPerc(per float64) string {
	return fmt.Sprintf("%.2f%%", per)
}

const shortLen = 12

func TruncateID(id string) string {
	if i := strings.IndexRune(id, ':'); i >= 0 {
		id = id[i+1:]
	}
	if len(id) > shortLen {
		id = id[:shortLen]
	}
	return id
}
