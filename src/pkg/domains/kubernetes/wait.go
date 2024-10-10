package kube

import (
	"context"
	"fmt"
	"time"

	pkgkubernetes "github.com/defenseunicorns/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/cli-utils/pkg/object"
)

func EvaluateWait(ctx context.Context, cluster *Cluster, waitPayload Wait) error {
	if cluster == nil {
		return fmt.Errorf("cluster is nil")
	}

	// TODO: incorporate wait for multiple objects?
	obj, err := globalCluster.validateAndGetGVR(waitPayload.Group, waitPayload.Version, waitPayload.Resource)
	if err != nil {
		return fmt.Errorf("unable to validate GVR: %v", err)
	}
	objMeta := object.ObjMetadata{
		Name:      waitPayload.Name,
		Namespace: waitPayload.Namespace,
		GroupKind: schema.GroupKind{
			Group: obj.Group,
			Kind:  obj.Kind,
		},
	}

	// Set timeout
	timeoutString := waitPayload.Timeout
	if timeoutString == "" {
		timeoutString = "30s"
	}

	// Timeout control parameters
	duration, err := time.ParseDuration(timeoutString)
	if err != nil {
		return fmt.Errorf("invalid wait timeout: %s", timeoutString)
	}
	waitCtx, waitCancel := context.WithTimeout(ctx, duration)
	defer waitCancel()

	return pkgkubernetes.WaitForReady(waitCtx, cluster.watcher, []object.ObjMetadata{objMeta})
}
