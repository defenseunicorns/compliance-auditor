package kube

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/kubectl/pkg/cmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/cmd/wait"
)

// This is specific to Lula - Check if we need to execute any wait operations.
func EvaluateWait(waitPayload Wait) error {
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
		} else {
			timeoutString = "5m"
		}

		// Timeout control parameters
		duration, err := time.ParseDuration(timeoutString)
		expiration := time.Now().Add(duration)
		startTime := time.Now()

		// Wait for existence
		err = WaitForExistence(waitPayload.Kind, waitPayload.Namespace, duration)
		if err != nil {
			return err
		}

		// If just waiting for existence - return here
		switch waitPayload.Condition {
		case "", "exist", "exists", "Exist", "Exists":
			return nil
		}

		// Calculate time remaining to explicitly pass as a timeout
		timeoutRemaining := expiration.Sub(startTime)

		err = WaitForCondition(forCondition, waitPayload.Namespace, timeoutRemaining.String(), waitPayload.Kind)
		if err != nil {
			return err
		}
	}
	return nil
}

func WaitForExistence(kind string, namespace string, timeout time.Duration) (err error) {
	expired := time.After(timeout)
	name := strings.Split(kind, "/")[1]

	for {
		// Delay check for 2 seconds
		time.Sleep(time.Second * 2)

		select {
		case <-expired:
			return fmt.Errorf("Timeout Expired")
		default:
			gvr, err := getGroupVersionResource(kind)
			if err != nil {
				return err
			}

			resourceRule := ResourceRule{
				Group:      gvr.Group,
				Version:    gvr.Version,
				Resource:   gvr.Resource,
				Namespaces: []string{namespace},
				Name:       name,
			}

			resources, err := GetResourcesDynamically(context.TODO(), resourceRule)
			if err != nil {
				return err
			}

			if len(resources) > 0 {
				// success
				return nil
			}
		}
	}
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
