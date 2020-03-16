package utils

import "fmt"

var sequence uint64 = 0
var timelineSequence uint64 = 0

type Timeline struct {
	Id    uint64 `json:"id"`
	Start uint64 `json:"start"`
	End   uint64 `json:"end"`
}

type Timelines []Timeline

func (ts Timelines) last() *Timeline {
	if len(ts) == 0 {
		return nil
	}
	return &ts[len(ts)-1]
}

func (ts *Timelines) append(t Timeline) Timelines {
	*ts = append(*ts, t)
	return *ts
}

type Task struct {
	id                uint64
	name              string
	period            uint64
	execTime          uint64
	execTimeRemaining uint64
	timelines         Timelines
}

func NewTask(period, execTime uint64) *Task {
	sequence++
	return &Task{
		id:                sequence,
		name:              fmt.Sprintf("Task %d", sequence),
		period:            period,
		execTime:          execTime,
		execTimeRemaining: execTime,
		timelines:         make(Timelines, 0),
	}
}

func newTimeline(start, end uint64) Timeline {
	timelineSequence++
	return Timeline{
		Id:    timelineSequence,
		Start: start,
		End:   end,
	}
}

func (t *Task) Spawn(moment uint64, pq *PriorityQueue, onAdd func(moment uint64, t *Task)) {
	if moment%t.period == 0 {
		t.execTimeRemaining = t.execTime
		pq.Add(&(*t), func(i interface{}) {
			task := i.(*Task)

			//t.timelines.append(newTimeline(moment - t.execTime + t.execTimeRemaining, moment))

			if last := task.timelines.last(); last != nil {
				last.End = moment
				//task.timelines.append(newTimeline(moment, moment))
				//task.timelines.append(newTimeline(last.End, moment))
			}
			//task.timelines.append(newTimeline(moment, moment))
			//task.timelines.append(newTimeline(task.timelines.last().End, moment))

			onAdd(moment, task)
		})
	}
}

func (t *Task) Execute(moment uint64, pq *PriorityQueue, onExec func(moment uint64, task *Task)) {
	if t.execTime == t.execTimeRemaining {
		t.timelines.append(newTimeline(moment, moment))
	}

	t.execTimeRemaining--

	if t.execTimeRemaining == 0 {
		t.timelines.last().End = moment

		//t.timelines.append(newTimeline(moment - t.execTime + t.execTimeRemaining, moment))

		//if last := t.timelines.last(); last != nil {
		//	last.End = moment
		//}

		//t.timelines.append(newTimeline(t.timelines.last().End, moment))

		onExec(moment, t)
		pq.Pop()
	}
}

func (t *Task) Id() uint64                { return t.id }
func (t *Task) Name() string              { return t.name }
func (t *Task) Timelines() []Timeline     { return t.timelines }
func (t *Task) ExecTime() uint64          { return t.execTime }
func (t *Task) ExecTimeRemaining() uint64 { return t.execTimeRemaining }
func (t *Task) Period() uint64            { return t.period }
