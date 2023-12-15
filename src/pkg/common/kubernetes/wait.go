package kube

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/defenseunicorns/lula/src/types"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	watchtools "k8s.io/client-go/tools/watch"
	"k8s.io/kubectl/pkg/cmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/cmd/wait"
)

var errNoMatchingResources = errors.New("no matching resources found")

func EvaluateWait(waitPayload types.Wait) error {
	var forCondition string
	waitCmd := false
	if waitPayload.Condition != "" {
		forCondition = fmt.Sprintf("condition=%s", waitPayload.Condition)
		waitCmd = true
	}

	if waitPayload.Jsonpath != "" {
		if waitCmd {
			return fmt.Errorf("only one of waitFor.condition or waitFor.jsonpath can be specified")
		}
		forCondition = fmt.Sprintf("jsonpath=%s", waitPayload.Jsonpath)
		waitCmd = true
	}

	if waitCmd {
		var timeoutString string
		if waitPayload.Timeout != "" {
			timeoutString = fmt.Sprintf("%s", waitPayload.Timeout)
		}

		err := WaitForCondition(forCondition, waitPayload.Kind, waitPayload.Namespace, timeoutString)
		if err != nil {
			return err
		}
	}
	return nil
}

func WaitForCondition(condition string, kind string, namespace string, timeout string) (err error) {
	// Required parameters for kubectl
	o := cmd.KubectlOptions{}
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
	kubeConfigFlags.Namespace = &namespace
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)

	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	flags := wait.NewWaitFlags(f, o.IOStreams)

	flags.ForCondition = condition
	if timeout != "" {
		flags.Timeout, err = time.ParseDuration(timeout)
		if err != nil {
			return err
		}
	}
	args := []string{kind}
	opts, err := flags.ToOptions(args)
	if err != nil {
		return err
	}
	err = RunModifiedWait(opts)
	if err != nil {
		return err
	}
	return nil

}

type ConditionFunc func(ctx context.Context, info *resource.Info, o *wait.WaitOptions) (finalObject runtime.Object, done bool, err error)

// ResourceLocation holds the location of a resource
type ResourceLocation struct {
	GroupResource schema.GroupResource
	Namespace     string
	Name          string
}

// UIDMap maps ResourceLocation with UID
type UIDMap map[ResourceLocation]k8stypes.UID

func RunModifiedWait(o *wait.WaitOptions) error {
	ctx, cancel := watchtools.ContextWithOptionalTimeout(context.Background(), o.Timeout)
	defer cancel()

	visitCount := 0
	visitFunc := func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		visitCount++
		finalObject, success, err := o.ConditionFn(ctx, info, o)
		if success {
			return nil
		}
		if err == nil {
			return fmt.Errorf("%v unsatisified for unknown reason", finalObject)
		}
		return err
	}
	visitor := o.ResourceFinder.Do()
	isForDelete := strings.ToLower(o.ForCondition) == "delete"
	if visitor, ok := visitor.(*resource.Result); ok && isForDelete {
		visitor.IgnoreErrors(apierrors.IsNotFound)
	}

	err := visitor.Visit(visitFunc)
	if err != nil {
		return err
	}
	if visitCount == 0 && !isForDelete {
		return errNoMatchingResources
	}
	return err
}
