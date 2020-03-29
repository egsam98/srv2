package utils

import "fmt"

type Marker struct {
	Time    uint64
	IsStart bool
}

type Task struct {
	id                uint64
	name              string
	period            uint64
	execTime          uint64
	execTimeRemaining uint64
	execMoments       []uint64
	markers           []Marker
	Count             uint64
	IsAperiodic       bool
}

func NewTask(id, period, execTime uint64, isAperiodic bool) *Task {
	return &Task{
		id: id,
		name: fmt.Sprintf("Задача №%d (p: %.01fs, e: %.01fs",
			id, float64(period)/1000, float64(execTime)/1000),
		period:            period,
		execTime:          execTime,
		execTimeRemaining: execTime,
		execMoments:       make([]uint64, 0),
		IsAperiodic:       isAperiodic,
	}
}

func (t *Task) CanSpawn(moment uint64) bool {
	return moment%t.period == 0
}

func (t *Task) Execute(moment uint64, pq *PriorityQueue, onPop func(*Task)) {
	if t.execTime == t.execTimeRemaining {
		t.markers = append(t.markers, Marker{
			Time:    moment,
			IsStart: true,
		})
	}

	t.execTimeRemaining--

	t.execMoments = append(t.execMoments, moment)

	if t.execTimeRemaining == 0 {
		pq.Pop()
		t.markers = append(t.markers, Marker{
			Time:    moment,
			IsStart: false,
		})
		onPop(t)
	}
}

func (t *Task) Id() uint64                { return t.id }
func (t *Task) SetName(name string)       { t.name = name }
func (t *Task) Name() string              { return t.name }
func (t *Task) ExecTime() uint64          { return t.execTime }
func (t *Task) ExecTimeRemaining() uint64 { return t.execTimeRemaining }
func (t *Task) Period() uint64            { return t.period }
func (t *Task) ExecMoments() []uint64     { return t.execMoments }
func (t *Task) Markers() []Marker         { return t.markers }
