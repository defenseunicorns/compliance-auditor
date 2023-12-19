package kube

import (
	"fmt"
	"os"
	"time"

	"github.com/defenseunicorns/lula/src/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/kubectl/pkg/cmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/cmd/wait"
)

// This is specific to Lula - Check if we need to execute any wait operations.
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

		err := WaitForCondition(forCondition, waitPayload.Namespace, timeoutString, waitPayload.Kind)
		if err != nil {
			return err
		}
	}
	return nil
}

// This is required bootstrapping for use of RunWait()
func WaitForCondition(condition string, namespace string, timeout string, args ...string) (err error) {
	// Required for printer - investigate exposing this as needed for modification
	ioStreams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	o := cmd.KubectlOptions{
		IOStreams: ioStreams,
	}
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
	// Namespace is attributed here
	kubeConfigFlags.Namespace = &namespace
	// Setup factory and flags
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	flags := wait.NewWaitFlags(f, o.IOStreams)
	// Add condition
	flags.ForCondition = condition
	if timeout != "" {
		flags.Timeout, err = time.ParseDuration(timeout)
		if err != nil {
			return err
		}
	}
	opts, err := flags.ToOptions(args)
	if err != nil {
		return err
	}
	err = opts.RunWait()
	if err != nil {
		return err
	}
	return nil
}
