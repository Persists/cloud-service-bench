package timeline

import (
	"fmt"
	"time"
)

type TimeLine struct {
	warmUpDuration     time.Duration
	warmUpAction       func()
	experimentDuration time.Duration
	experiment         func()
	coolDownDuration   time.Duration
	coolDownAction     func()
}

// SetWarmUp sets the duration and action for the WarmUp phase of the timeline.
func (t *TimeLine) SetWarmUp(duration time.Duration, action func()) {
	t.warmUpDuration = duration
	t.warmUpAction = action
}

// SetExperiment sets the duration and action for the Experiment phase of the timeline.
func (t *TimeLine) SetExperiment(duration time.Duration, action func()) {
	t.experimentDuration = duration
	t.experiment = action
}

// SetCoolDown sets the duration and action for the CoolDown phase of the timeline.
func (t *TimeLine) SetCoolDown(duration time.Duration, action func()) {
	t.coolDownDuration = duration
	t.coolDownAction = action
}

// Run executes the timeline actions in three phases: WarmUp, Experiment, and CoolDown.
// Each phase starts at a specific time and runs for a predefined duration.
//
// Parameters:
//   - startAt: The time at which the timeline should start.
//
// The timeline consists of the following phases:
//   - WarmUp: Starts at 'startAt' and runs for 't.warmUpDuration'.
//   - Experiment: Starts after the WarmUp phase and runs for 't.experimentDuration'.
//   - CoolDown: Starts after the Experiment phase and runs for 't.coolDownDuration'.
//
// The function will block until all phases are completed.
func (t *TimeLine) Run(startAt time.Time) {
	WarmUpEnd := startAt.Add(t.warmUpDuration)
	ExperimentEnd := WarmUpEnd.Add(t.experimentDuration)
	CoolDownEnd := ExperimentEnd.Add(t.coolDownDuration)

	<-time.After(time.Until(startAt))
	fmt.Println("Starting the WarmUp phase at", time.Now())
	go t.warmUpAction()

	<-time.After(time.Until(WarmUpEnd))
	fmt.Println("Starting the Experiment phase at", time.Now())
	go t.experiment()

	<-time.After(time.Until(ExperimentEnd))
	fmt.Println("Starting the CoolDown phase at", time.Now())
	go t.coolDownAction()

	<-time.After(time.Until(CoolDownEnd))
	fmt.Println("Experiment finished at", time.Now())
}
