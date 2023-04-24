/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/util/intstr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/api/errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webserverappv1 "github.com/huberychang/webserver-operator/api/v1"
)

// WebserverReconciler reconciles a Webserver object
type WebserverReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webserverapp.my.domain,resources=webservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webserverapp.my.domain,resources=webservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webserverapp.my.domain,resources=webservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Webserver object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *WebserverReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here
	instance := &webserverappv1.Webserver{}

	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Webserver resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, err
		}

		log.Error(err, "Failed to get Webserver")
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	// 检查webserver deployment是否已经存在，不存在的话就创建一个新的
	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      instance.Name,
		Namespace: instance.Namespace,
	}, deployment)

	if err != nil && errors.IsNotFound(err) {
		dep, err := r.createDeployment(instance)
		log.Info("create a Webserver deployment")
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment")
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	if err != nil {
		log.Error(err, "Failed to get deployment")
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	// 确保 deployment 和 image 与 spec 相同
	replicas := int32(instance.Spec.Replicas)
	image := instance.Spec.Image

	var needUpd bool
	if *deployment.Spec.Replicas != replicas {
		log.Info("Deployment spec.replicas change")
		deployment.Spec.Replicas = &replicas
		needUpd = true
	}

	if (*deployment).Spec.Template.Spec.Containers[0].Image != image {
		log.Info("Deployment spec.template.spec.container[0].image change")
		deployment.Spec.Template.Spec.Containers[0].Image = image
		needUpd = true
	}

	if needUpd {
		err := r.Update(ctx, deployment)
		if err != nil {
			log.Error(err, "Failed to update deployment")
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// 检查 webserver service 是否已经存在，不存在的话就创建一个新的
	deploymentsvc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{
		Namespace: instance.Namespace,
		Name:      instance.Name + "-service",
	}, deploymentsvc)
	if err != nil && errors.IsNotFound(err) {
		svc, err := r.createService(instance)
		log.Info("create a new service")
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create service")
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	if err != nil {
		log.Error(err, "Failed to get service")
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	return ctrl.Result{RequeueAfter: time.Second * 10}, nil
}

func (r *WebserverReconciler) createDeployment(server *webserverappv1.Webserver) (*appsv1.Deployment, error) {
	replicas := int32(server.Spec.Replicas)
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      server.Name,
			Namespace: server.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"server": server.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"server": server.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx-server",
							Image: server.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
									Protocol:      "TCP",
								},
							},
						},
					},
				},
			},
		},
	}
	// 将Web服务器实例设置为所有者和控制器
	if err := ctrl.SetControllerReference(server, deploy, r.Scheme); err != nil {
		return deploy, err
	}
	return deploy, nil
}

func (r *WebserverReconciler) createService(server *webserverappv1.Webserver) (*corev1.Service, error) {
	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			// Name:      server.Name + "-service",
			Name:      server.Name,
			Namespace: server.Namespace,
			// Labels:    map[string]string{"service": server.Name},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					NodePort:   30010,
					Port:       80,
					TargetPort: intstr.FromString("http"),
				},
			},
			Selector: map[string]string{
				//"app":    "webserver",
				"server": server.Name,
			},
		},
	}
	if err := ctrl.SetControllerReference(server, svc, r.Scheme); err != nil {
		return svc, err
	}
	return svc, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebserverReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webserverappv1.Webserver{}).
		Complete(r)
}
