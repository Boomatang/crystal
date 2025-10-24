package applicaton

import (
	"github.com/boomatang/crystal/internal/workflow"
)

func NewApplictaionOne() *workflow.World {
	actionAlan := workflow.NewPoint(alan)
	actionTom := workflow.NewPoint(tom)
	actionJohn := workflow.NewPoint(john)
	actionMark := workflow.NewPoint(mark)
	worldOne := workflow.NewWorld("Temporay world")
	worldOne.PreCondition = actionTom
	worldOne.PostCondition = actionJohn
	worldOne.ErrorHandler = actionMark
	for range 10 {
		worldOne.AddAction(actionAlan)
	}

	worldMain := workflow.NewWorld("Main World")
	worldMain.PreCondition = actionAlan
	worldMain.AddAction(actionAlan)
	worldMain.AddAction(actionTom)
	worldMain.AddAction(actionMark)
	worldMain.AddAction(actionJohn)
	worldMain.AddAction(worldOne)
	worldMain.PostCondition = actionJohn

	return worldMain
}

func NewApplictaionTwo() *workflow.World {
	actionSleepy := workflow.NewPoint(sleepy)
	worldMain := workflow.NewWorld("Example two main world")
	worldMain.AddAction(actionSleepy)

	actionMark := workflow.NewPoint(mark)
	actionPostCondition := workflow.NewPoint(postCondition)
	worldMain.PreCondition = actionMark
	worldMain.PostCondition = actionPostCondition

	return worldMain
}
