package util

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func GetDeployment(path string) (*appsv1.Deployment, error) {
	path = filepath.Clean(path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	deployment := &appsv1.Deployment{}
	err = yaml.Unmarshal(bytes, &deployment)
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func GetClusterRole(path string) (*rbacv1.ClusterRole, error) {
	path = filepath.Clean(path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	clusterRole := &rbacv1.ClusterRole{}
	err = yaml.Unmarshal(bytes, &clusterRole)
	if err != nil {
		return nil, err
	}
	return clusterRole, nil
}

func GetPod(path string) (*v1.Pod, error) {
	path = filepath.Clean(path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pod := &v1.Pod{}
	err = yaml.Unmarshal(bytes, &pod)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func GetConfigMap(path string) (*v1.ConfigMap, error) {
	path = filepath.Clean(path)
	bytes, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return nil, err
	}
	configMap := &v1.ConfigMap{}
	err = yaml.Unmarshal(bytes, &configMap)
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

func GetService(path string) (*v1.Service, error) {
	path = filepath.Clean(path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	service := &v1.Service{}
	err = yaml.Unmarshal(bytes, &service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func GetSecret(path string) (*v1.Secret, error) {
	path = filepath.Clean(path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	secret := &v1.Secret{}
	err = yaml.Unmarshal(bytes, &secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func GetIngress(path string) (*netv1.Ingress, error) {
	path = filepath.Clean(path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ingress := &netv1.Ingress{}
	err = yaml.Unmarshal(bytes, &ingress)
	if err != nil {
		return nil, err
	}
	return ingress, nil
}

func GetNamespace(name string) (*v1.Namespace, error) {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}, nil
}

func ExecuteCommand(cmd *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	_, output, err = ExecuteCommandC(cmd, args...)
	return cmd, output, err
}

func ExecuteCommandC(cmd *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	execErr := cmd.Execute()

	out, readErr := io.ReadAll(buf)
	if readErr != nil {
		return cmd, "", readErr
	}

	return cmd, string(out), execErr
}
