package applicaton

import (
	"github.com/boomatang/crystal/internal/workflow"
)

var DeploymentConfigMapLinker = workflow.Link{
	Parent: "Deployment",
	Child:  "ConfigMap",
	LinkFunc: func(p, c *workflow.Node) bool {
		return p.Name == c.Name
	},
}
