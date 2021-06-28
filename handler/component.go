package handler

import (
	"context"
	"fmt"

	"github.com/enescakir/emoji"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// func getClusterID(C *kubernetes.Clientset) {
// 	clusterID,err := C.CoreV1().ComponentStatuses().List(context.TODO(),metav1.ListOptions{})
// }
func ComponentCheck(C *kubernetes.Clientset) {
	clusterData, err := C.CoreV1().ComponentStatuses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, n := range clusterData.Items {
		for _, condition := range n.Conditions {
			if condition.Type != "Healthy" {
				alertData := AlertSlack{
					Text:      fmt.Sprintf("*UnExpected Component state* : \n\t *%s* is in *%s* state", n.ObjectMeta.Name, condition.Type),
					Username:  "kubernetes Component",
					IconEmoji: ":warning:",
				}
				NotifyStatus(alertData)
			} else {
				fmt.Printf("\n > Component *%s* is *%s* \n", n.ObjectMeta.Name, condition.Type)
			}
		}
	}
}
func NodesCheck(C *kubernetes.Clientset) {
	nodes, err := C.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, n := range nodes.Items {
		fmt.Printf("\nNode : [ %s ] component status ", n.Name)
		for _, condition := range n.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status == "True" {
					fmt.Printf("\n\t %v Node component [ %s ] -> %s state", emoji.CheckMarkButton, condition.Type, condition.Status)
				} else {
					fmt.Printf("\n\t %v Node component [ %s ] -> %s state", emoji.CrossMark, condition.Type, condition.Status)
					alertData := AlertSlack{
						Text:      fmt.Sprintf(" *Unexpected Node state* : \n\t *%s* is in *%s* state", condition.Type, condition.Status),
						Username:  "kubernetes Node",
						IconEmoji: ":warning:",
					}
					NotifyStatus(alertData)
				}
			} else {
				if condition.Status != "False" {
					fmt.Printf("\n\t %v Node component [ %s ] -> %s state", emoji.CrossMark, condition.Type, condition.Status)
					alertData := AlertSlack{
						Text:      fmt.Sprintf(" *Unexpected Node state* : \n\t *%s* is in *%s* state", condition.Type, condition.Status),
						Username:  "kubernetes Node",
						IconEmoji: ":warning:",
					}
					NotifyStatus(alertData)
				} else {
					fmt.Printf("\n\t %v Node component [ %s ] -> %s state", emoji.CheckMarkButton, condition.Type, condition.Status)
				}

			}
		}
	}

}
func PodsCheck(C *kubernetes.Clientset) {
	getNS, err := C.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, ns := range getNS.Items {
		// fmt.Print("\n\n")
		// PrettyPrint(ns.Name)
		pods, err := C.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("\nThere are %d pods in the ' %s ' namespace", len(pods.Items), ns.Name)
		for _, po := range pods.Items {
			if po.Status.Phase == "Running" || po.Status.Phase == "Succeeded" {
				fmt.Printf(" \n\t %v pod: [ %s ] is in %s state", emoji.CheckMarkButton, po.Name, po.Status.Phase)
			} else {
				fmt.Printf(" \n\t %v pod: [ %s ] is in %s state", emoji.CrossMark, po.Name, po.Status.Phase)
				alertData := AlertSlack{
					Text:      fmt.Sprintf(" *Unexpected Pod state* : \n\t _*%s*_ is in *%s* state _%s_", po.Name, po.Status.Phase, po.Status.Message),
					Username:  "kubernetes Pods",
					IconEmoji: ":warning:",
				}
				NotifyStatus(alertData)
			}
		}
	}

}
