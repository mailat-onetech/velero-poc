package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os/exec"
)

func main() {
	//CreateBackup("bn","default")
	//return
	config, err := clientcmd.BuildConfigFromFlags("", "./kubeconfig.yaml")
	if err != nil {
		return
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}
	nodeList := ListNode(c)
	fmt.Println(nodeList)

	//res:=CreateNamespace(c,"test-ns")
	//fmt.Println(res)
	CreateBackup("backup-test-ns", "test-ns")
}

func ListNode(c *kubernetes.Clientset) *corev1.NodeList {
	nodeList, err := c.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, n := range nodeList.Items {
		fmt.Println(n.Name)
	}
	return nodeList
}

func ListPod(c *kubernetes.Clientset, nameSpace string) *corev1.PodList {
	podList, err := c.CoreV1().Pods(nameSpace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, n := range podList.Items {
		fmt.Println(n.Name)
	}
	return podList
}

func CreatePod(c *kubernetes.Clientset, nameSpace string, podName string) *corev1.Pod {
	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "busybox", Image: "busybox:latest", Command: []string{"sleep", "100000"}},
			},
		},
	}

	pod, err := c.CoreV1().Pods(nameSpace).Create(context.Background(), newPod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(pod)

	return pod
}

func CreateNamespace(c *kubernetes.Clientset, nameSpace string) *corev1.Namespace {
	newNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: nameSpace,
		},
	}

	ns, err := c.CoreV1().Namespaces().Create(context.Background(), newNs, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(ns)

	return ns
}

func DeletePod(c *kubernetes.Clientset, nameSpace string, podName string) {

	err := c.CoreV1().Pods(nameSpace).Delete(context.Background(), podName, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
}

func CreateBackup(backupName string, nameSpace string) {
	backupCommand := "velero backup create " + backupName + " --include-namespaces " + nameSpace + " --kubeconfig ./kubeconfig.yaml"
	fmt.Println("==>", backupCommand)

	cmd := exec.Command(backupCommand)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}

func InitRestore(backupName string) {
	restoreCommand := "velero restore create --from-backup " + backupName
	fmt.Println("==>", restoreCommand)

	cmd := exec.Command(restoreCommand)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}
