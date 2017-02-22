package handlers

import (
	"strconv"

	"../datasource"
	"gopkg.in/mgo.v2/bson"
)

func BuildTrafficArray(deviceID bson.ObjectId) map[string][][]string {
	mb := datasource.MonitoringBase{}
	devstatusData := mb.LoadDeviceCheckDataBy(deviceID)
	var input [][]string
	var output [][]string

	for _, stat := range devstatusData {
		if stat.CheckType == "bytes-in" {
			tmp := []string{strconv.FormatInt(stat.Time.Unix(), 10), stat.Value}
			input = append(input, tmp)
		}
		if stat.CheckType == "bytes-out" {
			tmp := []string{strconv.FormatInt(stat.Time.Unix(), 10), stat.Value}
			output = append(output, tmp)
		}
	}

	result := map[string][][]string{
		"input":  input,
		"output": output,
	}
	return result
}
