package handlers

import (
	"bytes"
    "log"
    "os/exec"
    "strconv"
    "strings"
)

func CPU_load() float64 {
	var load float64
	cmd := exec.Command("ps", "aux")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil { log.Fatal(err) }

    for {
        line, err := out.ReadString('\n')
        if err!=nil { break }
        tokens := strings.Split(line, " ")
        ft := make([]string, 0)
        for _, t := range(tokens) {
            if t!="" && t!="\t" {
                ft = append(ft, t)
            }
        }
        cpu, err := strconv.ParseFloat(strings.Replace(ft[2], ",", ".", -1), 64)
        if err != nil { continue }
        load +=cpu
    }
    
    return load
}