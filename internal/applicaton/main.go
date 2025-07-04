package applicaton

import (
	"fmt"

	"github.com/boomatang/crystal/internal/workflow"
)

func alan() error {
	fmt.Println("This was a function call")
	return nil
}
func tom() error {
	fmt.Println("This was a function call")
	return nil
}
func john() error {
	fmt.Println("This was a function call")
	return nil
}
func mark() error {
	fmt.Println("This was a function call")
	return nil
}
func NewApplictaion() *workflow.World {
	actionAlan := workflow.NewPoint(alan)
	actionTom := workflow.NewPoint(tom)
	actionJohn := workflow.NewPoint(john)
	actionMark := workflow.NewPoint(mark)
	worldOne := workflow.NewWorld("Temporay world")
	worldOne.PreCondition = actionTom
	worldOne.PostCondition = actionJohn
	worldOne.ErrorHandler = actionMark
	for i := 0; i < 10; i++ {
		worldOne.AddAction(actionAlan)
	}

	worldMain := workflow.NewWorld("Main World")
	worldMain.PreCondition = worldOne
	worldMain.AddAction(actionAlan)
	worldMain.AddAction(actionTom)
	worldMain.AddAction(actionMark)
	worldMain.AddAction(actionJohn)
	worldMain.AddAction(worldOne)
	worldMain.PostCondition = worldOne

	return worldMain
}
