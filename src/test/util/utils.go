package util

import (
	"bytes"
	"io"
	"os"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func GetDeployment(deploymentFilePath string) (*appsv1.Deployment, error) {
	bytes, err := os.ReadFile(deploymentFilePath)
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

func GetClusterRole(clusterRoleFilePath string) (*rbacv1.ClusterRole, error) {
	bytes, err := os.ReadFile(clusterRoleFilePath)
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

func GetPod(podFilePath string) (*v1.Pod, error) {
	bytes, err := os.ReadFile(podFilePath)
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

func GetConfigMap(configMapFilePath string) (*v1.ConfigMap, error) {
	bytes, err := os.ReadFile(configMapFilePath)
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

func GetService(serviceFilePath string) (*v1.Service, error) {
	bytes, err := os.ReadFile(serviceFilePath)
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

func GetSecret(secretFilePath string) (*v1.Secret, error) {
	bytes, err := os.ReadFile(secretFilePath)
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

func GetIngress(ingressFilePath string) (*netv1.Ingress, error) {
	bytes, err := os.ReadFile(ingressFilePath)
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

func ExecuteCommand(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	_, output, err = ExecuteCommandC(root, args...)
	return root, output, err
}

func ExecuteCommandC(cmd *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	cmd.Execute()

	out, err := io.ReadAll(buf)

	return cmd, string(out), err
}
