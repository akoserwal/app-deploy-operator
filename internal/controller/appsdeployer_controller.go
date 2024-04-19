/*
Copyright 2024.

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

package controller

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	applyappsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	applymetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/akoserwal/app-deploy-operator/api/v1alpha1"
)

// AppsdeployerReconciler reconciles a Appsdeployer object
type AppsdeployerReconciler struct {
	client.Client
	Kclient kubernetes.Interface
	Scheme  *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.akoserwal,resources=appsdeployers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.akoserwal,resources=appsdeployers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.akoserwal,resources=appsdeployers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Appsdeployer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *AppsdeployerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	apps := &appsv1alpha1.Appsdeployer{}
	err := r.Get(ctx, req.NamespacedName, apps)
	if err != nil {
		return ctrl.Result{}, err
	}
	logger.Info("request", "namespace:", req.Namespace)

	for _, appDev := range apps.Spec.Apps {
		fieldMgr := "appdeployed-manager"

		_, err = r.Kclient.CoreV1().Namespaces().Get(ctx, appDev.Namespace, metav1.GetOptions{})
		if err != nil {
			names := v1.Namespace{ObjectMeta: metav1.ObjectMeta{
				Name: appDev.Namespace,
			}}
			r.Kclient.CoreV1().Namespaces().Create(ctx, &names, metav1.CreateOptions{FieldManager: fieldMgr})
			if err != nil {
				logger.Error(err, "failed to create namespace")
			}
		}

		// Create a Deployment spec.

		appconfig := deployment(appDev.Name, appDev.Namespace, appDev.Image, appDev.Port)
		dep := r.Kclient.AppsV1().Deployments(appDev.Namespace)
		// extract

		resp, err := dep.Apply(ctx, appconfig, metav1.ApplyOptions{FieldManager: fieldMgr})
		if err != nil {
			log.Log.Error(err, "failed to apply dev")
		}

		logger.Info("resp post apply:", "resp", resp)

		depSvc := r.Kclient.CoreV1().Services(appDev.Namespace)
		svc := service(appDev.Name, appDev.Namespace, appDev.Port)
		respsvc, err := depSvc.Apply(ctx, svc, metav1.ApplyOptions{FieldManager: fieldMgr})
		if err != nil {
			log.Log.Error(err, "failed to apply svc")
		}
		logger.Info("service post apply:", "service", respsvc)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppsdeployerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.Appsdeployer{}).
		Complete(r)
}

func deployment(n string, namespace string, image string, port int32) *applyappsv1.DeploymentApplyConfiguration {
	name := n
	return applyappsv1.Deployment(name, namespace).
		WithLabels(map[string]string{"app": name}).
		WithSpec(applyappsv1.DeploymentSpec().
			WithReplicas(1).
			WithSelector(applymetav1.LabelSelector().WithMatchLabels(map[string]string{"app.kubernetes.io/instance": name})).
			WithTemplate(applycorev1.PodTemplateSpec().WithLabels(map[string]string{"app.kubernetes.io/instance": name}).
				WithSpec(applycorev1.PodSpec().WithContainers(
					applycorev1.Container().WithName(n).WithImage(image).
						WithPorts(applycorev1.ContainerPort().WithContainerPort(port).WithName(n)),
				))))
}

func service(name string, namespace string, port int32) *applycorev1.ServiceApplyConfiguration {
	s := applycorev1.Service(name, namespace)
	s.WithName(name).WithNamespace(namespace).WithSpec(applycorev1.ServiceSpec().WithPorts(applycorev1.ServicePort().WithName(name).WithPort(port)))
	return s
}
