package services

import (
	"fmt"
	. "srv2/utils"
	"strings"
)

var RM = func(pq *PriorityQueue) func(i int, j int) bool {
	return func(i int, j int) bool {
		t1 := pq.Get(i).(*Task)
		t2 := pq.Get(j).(*Task)
		if t2.IsAperiodic && !t1.IsAperiodic {
			return true
		} else if !t2.IsAperiodic && t1.IsAperiodic {
			return false
		}

		if t1.Id() == t2.Id() {
			return t1.ExecTimeRemaining() < t2.ExecTimeRemaining()
		}
		return t1.Period() < t2.Period()
	}
}

var EDF = func(pq *PriorityQueue) func(i int, j int) bool {
	return func(i int, j int) bool {
		t1 := pq.Get(i).(*Task)
		t2 := pq.Get(j).(*Task)
		if t2.IsAperiodic && !t1.IsAperiodic {
			return true
		} else if !t2.IsAperiodic && t1.IsAperiodic {
			return false
		}
		return t1.Count*t1.Period() < t2.Count*t2.Period()
	}
}

type SchedullingService struct {
	tasksConfig []*Task
}

func NewSchedullingService(tasksConfig []*Task) *SchedullingService {
	return &SchedullingService{tasksConfig}
}

func (ss *SchedullingService) Run(method string) (string, []map[string]interface{}, error) {

	var pq *PriorityQueue = nil
	switch method {
	case "rm":
		pq = NewPriorityQueue(RM)
	case "edf":
		pq = NewPriorityQueue(EDF)
	default:
		return "", nil, fmt.Errorf("must be \"rm\" or \"edf\" as path param")
	}

	tasksOut := make([]Task, 0)
	maxIters := ss.hyperPeriod()
	for moment := uint64(0); moment < maxIters; moment++ {
		for _, taskConf := range ss.tasksConfig {
			if taskConf.CanSpawn(moment) {
				cloneAndSpawn(taskConf, pq)
			}
		}

		if t := pq.Peek(); t != nil {
			t.(*Task).Execute(moment, pq, func(task *Task) {
				tasksOut = append(tasksOut, *task)
			})
		}
	}

	for _, conf := range ss.tasksConfig {
		conf.Count = 0
	}

	title := fmt.Sprintf("Алгоритм %s. Суммарная загруженность: %.3f",
		strings.ToUpper(method), summaryLoad(ss.tasksConfig))
	return title, formTrace(tasksOut), nil
}

func (ss *SchedullingService) hyperPeriod() uint64 {
	var max uint64 = 0
	for _, conf := range ss.tasksConfig {
		if conf.Period() > max {
			max = conf.Period()
		}
	}
	return max
}

func summaryLoad(tasks []*Task) float64 {
	sum := .0
	for _, task := range tasks {
		sum += float64(task.ExecTime()) / float64(task.Period())
	}
	return sum
}

func cloneAndSpawn(taskConf *Task, pq *PriorityQueue) {
	inst := NewTask(taskConf.Id(), taskConf.Period(), taskConf.ExecTime(), taskConf.IsAperiodic)
	inst.SetName(taskConf.Name())
	taskConf.Count += 1
	inst.Count = taskConf.Count
	pq.Add(inst)
}

func formTrace(tasksOut []Task) []map[string]interface{} {
	traceData := make([]map[string]interface{}, 0)
	for _, task := range tasksOut {
		markers := make([]map[string]interface{}, 0)

		for _, marker := range task.Markers() {
			markerData := map[string]interface{}{
				"type":  "diamond",
				"value": marker.Time,
			}
			if !marker.IsStart {
				markerData["fill"] = "#FF0000"
			}
			markers = append(markers, markerData)
		}

		periods := make([]map[string]uint64, 0)
		for _, execMoment := range task.ExecMoments() {
			periodsLen := len(periods)
			if periodsLen > 0 {
				last := periods[periodsLen-1]
				if execMoment-last["end"] == 1 {
					last["end"] = execMoment
					continue
				}
			}
			periods = append(periods, map[string]uint64{
				"start": execMoment,
				"end":   execMoment,
			})
		}

		traceData = append(traceData, map[string]interface{}{
			"id":      task.Id(),
			"name":    task.Name(),
			"p":       float64(task.Period()) / 1000,
			"e":       float64(task.ExecTime()) / 1000,
			"markers": markers,
			"periods": periods,
		})
	}
	return traceData
}
