package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
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

func get_events() []Event {

	arr := []Event{}

	// Открываем файл
	file, err := os.Open("log.log")
	if err != nil {
		log.Fatal("Ошибка при открытии файла:", err)
	}
	defer file.Close()

	// Создаем буферизированный считыватель
	scanner := bufio.NewScanner(file)

	const maxCapacity = 10 * 1024 * 1024 // 10 МБ
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	// Читаем файл построчно
	for scanner.Scan() {
		// Получаем строку
		line_json := scanner.Text()

		var event Event

		// Десериализуем JSON в структуру
		err := json.Unmarshal([]byte(line_json), &event)
		if err != nil {
			log.Fatal("Ошибка при десериализации JSON:", err)
		}

		arr = append(arr, event)
	}

	// Проверяем на ошибки при чтении
	if err := scanner.Err(); err != nil {
		log.Fatal("Ошибка при чтении файла:", err)
	}
	return arr
}
