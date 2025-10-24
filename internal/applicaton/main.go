package applicaton

import (
	"fmt"

	"github.com/boomatang/crystal/internal/api"
	"github.com/boomatang/crystal/internal/workflow"
)

func alan(nodes *workflow.NodeList) error {
	fmt.Println("This was a function call")
	deployments := nodes.GetNodes("Deployment")
	for _, node := range deployments {
		fmt.Printf("This node was returned, %s\n", node)
		deployment := &api.Deployment{}
		err := node.Object(deployment)
		if err != nil {
			fmt.Println("Some thing bad happened trying to get the Object")
		}
		fmt.Printf("Number of replicas: %v\n", deployment.Spec.Replicas)
	}
	return nil
}
func tom(nodes *workflow.NodeList) error {
	fmt.Println("This was a function call")
	return nil
}
func john(nodes *workflow.NodeList) error {
	fmt.Println("This was a function call")
	return nil
}
func mark(nodes *workflow.NodeList) error {
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
	for range 10 {
		worldOne.AddAction(actionAlan)
	}

	worldMain := workflow.NewWorld("Main World")
	worldMain.PreCondition = actionAlan
	worldMain.AddAction(actionAlan)
	worldMain.AddAction(actionTom)
	worldMain.AddAction(actionMark)
	worldMain.AddAction(actionJohn)
	// worldMain.AddAction(worldOne)
	worldMain.PostCondition = actionJohn

	return worldMain
}
