package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

type MetaData struct{}

type Metric struct {
	AccumulatorId int    `json:"accumulatorId"`
	MetricType    string `json:"metricType"`
	Name          string `json:"name"`
}

type SparkPlan struct {
	NodeName     string      `json:"nodeName"`
	SimpleString string      `json:"simpleString"`
	MetaData     MetaData    `json:"metadata"`
	Metrics      []Metric    `json:"metrics"`
	Children     []SparkPlan `json:"children"`
}

type Event struct {
	EventName               string    `json:"Event"`
	ExecutionId             int       `json:"executionId"`
	Description             string    `json:"description"`
	Details                 string    `json:"details"`
	PhysicalPlanDescription string    `json:"physicalPlanDescription"`
	SparkPlanInfo           SparkPlan `json:"sparkPlanInfo"`
	Time                    int64     `json:"time"`
}

func getEvents() []Event {

	arr := []Event{}

	file, err := os.Open("log.log")
	if err != nil {
		log.Fatal("Opening file error:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	const maxCapacity = 10 * 1024 * 1024 // 10 МБ
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		lineJson := scanner.Text()

		var event Event

		err := json.Unmarshal([]byte(lineJson), &event)
		if err != nil {
			log.Fatal("Deserialization error JSON:", err)
		}

		arr = append(arr, event)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Reading file error:", err)
	}
	return arr
}

func getSources() {
	for _, ev := range pageMain.Events {
		for _, child := range ev.SparkPlanInfo.Children {
			getSrc(child)
		}
	}
	pageSrc.TextArea.SetText(strings.Join(pagePlan.sSrc, "\n"), true)
}

func getSrc(sp SparkPlan) {
	if len(sp.Children) == 0 {
		checkSrc(sp)
	} else {
		for _, child := range sp.Children {
			getSrc(child)
		}
	}
}

func checkSrc(sp SparkPlan) {
	if strings.HasPrefix(sp.NodeName, "Scan") {
		pagePlan.sSrc = append(pagePlan.sSrc, sp.SimpleString)
	}
}
