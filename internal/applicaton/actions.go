package applicaton

import (
	"time"

	"github.com/boomatang/crystal/internal/api"
	"github.com/boomatang/crystal/internal/logger"
	"github.com/boomatang/crystal/internal/workflow"
)

func alan(nodes *workflow.NodeList) error {
	logger.Log.Info("This was a function call")
	deployments := nodes.GetNodes("Deployment")
	for _, node := range deployments {
		logger.Log.Info("This node was returned", "node", node)
		deployment := &api.Deployment{}
		err := node.Object(deployment)
		if err != nil {
			logger.Log.Info("Some thing bad happened trying to get the Object")
		}
		if deployment.Spec.Replicas != nil {
			logger.Log.Info("Number of replicas", "count", *deployment.Spec.Replicas)
		} else {
			logger.Log.Info("Number of replicas", "count", "nil")
		}
	}
	return nil
}

func tom(nodes *workflow.NodeList) error {
	logger.Log.Info("This was a function call")
	return nil
}

func john(nodes *workflow.NodeList) error {
	logger.Log.Info("This was a function call")
	return nil
}

func mark(nodes *workflow.NodeList) error {
	logger.Log.Info("This was a function call")
	return nil
}

func postCondition(nodes *workflow.NodeList) error {
	logger.Log.Info("This is the post condition, and the last action to be be called")
	return nil
}

func sleepy(nodes *workflow.NodeList) error {
	err := alan(nodes)
	if err != nil {
		logger.Log.Error("some unexpected error happend in alan call", "error", err)
		return err
	}

	logger.Log.Info("running a sleepy function")
	time.Sleep(5 * time.Second)
	logger.Log.Info("finished sleep")
	return nil
}
