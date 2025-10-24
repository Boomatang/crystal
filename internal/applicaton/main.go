package applicaton

import (
	"github.com/boomatang/crystal/internal/api"
	"github.com/boomatang/crystal/internal/workflow"
	"github.com/boomatang/crystal/pkg/logger"
)

func alan(nodes *workflow.NodeList) error {
	logger.Log.Debug("alan action started")
	deployments := nodes.GetNodes("Deployment")
	for _, node := range deployments {
		logger.Log.Debug("processing deployment node", "node", node.String())
		deployment := &api.Deployment{}
		err := node.Object(deployment)
		if err != nil {
			logger.Log.Error("failed to convert node to deployment",
				"error", err,
				"node_kind", node.Kind,
				"node_name", node.Name)
			continue
		}
		logger.Log.Info("deployment replicas",
			"replicas", *deployment.Spec.Replicas,
			"deployment", deployment.Metadata.Name)
	}
	return nil
}
func tom(nodes *workflow.NodeList) error {
	logger.Log.Debug("tom action started")
	node := nodes.Get("Deployment", "tree2")
	if node == nil {
		logger.Log.Debug("node not found yet",
			"node_kind", "Deployment",
			"node_name", "tree2")
		return nil
	}

	subGraph, err := nodes.GetSubGraph(node, workflow.NodeOptList{})
	if err != nil {
		logger.Log.Error("failed to extract subgraph",
			"error", err,
			"node_kind", node.Kind,
			"node_name", node.Name)
		return err
	}

	if subGraph == nil {
		logger.Log.Warn("subgraph is nil")
		return nil
	}

	logger.Log.Info("subgraph extracted", "node_count", subGraph.Len())
	return nil
}
func john(nodes *workflow.NodeList) error {
	logger.Log.Debug("john action started")
	return nil
}
func mark(nodes *workflow.NodeList) error {
	logger.Log.Debug("mark action started")
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
