package applicaton

import (
	"strings"

	"github.com/boomatang/crystal/internal/workflow"
)

var DeploymentConfigMapLinker = workflow.Link{
	Parent: "Deployment",
	Child:  "ConfigMap",
	LinkFunc: func(p, c *workflow.Node) bool {
		return p.Name == c.Name
	},
}

var ConfigSecretMapLinker = workflow.Link{
	Parent: "ConfigMap",
	Child:  "Secret",
	LinkFunc: func(p, c *workflow.Node) bool {
		return strings.HasPrefix(p.Name, c.Name)
	},
}
