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

	firstInput := "0"
	firstOutput := "0"

	for _, stat := range devstatusData {
		if stat.CheckType == "bytes-in" {
			if firstInput == "0" {
				firstInput = stat.Value
			}
			tmp := []string{strconv.FormatInt(stat.Time.Unix(), 10), BitPerSecond(firstInput, stat.Value)}
			input = append(input, tmp)
			firstInput = stat.Value
		}
		if stat.CheckType == "bytes-out" {
			if firstOutput == "0" {
				firstOutput = stat.Value
			}
			tmp := []string{strconv.FormatInt(stat.Time.Unix(), 10), BitPerSecond(firstOutput, stat.Value)}
			output = append(output, tmp)
			firstOutput = stat.Value
		}
	}

	result := map[string][][]string{
		"input":  input,
		"output": output,
	}
	return result
}

func BitPerSecond(first string, next string) string {
	num1, _ := strconv.ParseInt(first, 10, 64)
	num2, _ := strconv.ParseInt(next, 10, 64)
	result := (num2 - num1)
	if result < 0 {
		result = 0
	}
	return strconv.FormatInt(result, 10)
}
