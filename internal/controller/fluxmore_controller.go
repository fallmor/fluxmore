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
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	fluxitv1alpha "github.com/fallmor/fluxmore/api/v1alpha"
	helm "github.com/fluxcd/helm-controller/api/v2"
)

// FluxMoreReconciler reconciles a FluxMore object
type FluxMoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// type PassVal struct {
// 	Checksuccess *bool
// }

var Checksuccess bool

// +kubebuilder:rbac:groups=fluxit.morbolt.dev,resources=fluxmores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fluxit.morbolt.dev,resources=fluxmores/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=fluxit.morbolt.dev,resources=fluxmores/finalizers,verbs=update
// +kubebuilder:rbac:groups=*,resources=secrets;configmaps;deployments;statefulsets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FluxMore object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *FluxMoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var fluxmore fluxitv1alpha.FluxMore
	var helmrelease helm.HelmRelease

	// ResourceName := Fluxmore.Spec.ResourcesCheck
	if err := r.Get(ctx, req.NamespacedName, &fluxmore); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)

	}
	if err := r.Get(ctx, types.NamespacedName{Name: fluxmore.Spec.HelmReleaseName, Namespace: fluxmore.Namespace}, &helmrelease); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)

	}
	switch kind := fluxmore.Spec.ResourcesKind; kind {
	case "configmap":

		var ConfigMap corev1.ConfigMap
		if helmrelease.Spec.Suspend {
			err := r.Get(ctx, types.NamespacedName{Name: fluxmore.Spec.ResourcesCheck, Namespace: fluxmore.Namespace}, &ConfigMap)
			// we use MergeForm instead of DeepCopy only to avoid conflict
			statusPatch := client.MergeFrom(fluxmore.DeepCopy())
			timeReconcile := metav1.NewTime(time.Now())

			if err != nil {
				if errors.IsNotFound(err) {
					// Configmap not found
					fluxmore.Status.ResourcesCheckFound = false
					l.Info("Configmap not found",
						"Namespace", fluxmore.Namespace,
						"Name", fluxmore.Spec.ResourcesCheck)
				} else {

					return ctrl.Result{}, err
				}
			} else {
				// No error it means Configmap exists
				fluxmore.Status.ResourcesCheckFound = true
				l.Info("Configmpae found",
					"Namespace", fluxmore.Namespace,
					"Name", fluxmore.Spec.ResourcesCheck)
			}
			fluxmore.Status.LastReconcileTime = &timeReconcile

			if err := r.Status().Patch(ctx, &fluxmore, statusPatch); err != nil {
				l.Error(err, "Unable to update Fluxmore status")
				return ctrl.Result{}, err
			}
		} else {
			l.Info("Nothing to do with the HelmRelease, already patched", "Name", fluxmore.Spec.HelmReleaseName)
		}

	case "secret":

		var Secret corev1.Secret
		if helmrelease.Spec.Suspend {
			err := r.Get(ctx, types.NamespacedName{Name: fluxmore.Spec.ResourcesCheck, Namespace: fluxmore.Namespace}, &Secret)
			statusPatch := client.MergeFrom(fluxmore.DeepCopy())
			timeReconcile := metav1.NewTime(time.Now())

			if err != nil {
				if errors.IsNotFound(err) {
					// Secret not found
					fluxmore.Status.ResourcesCheckFound = false
					l.Info("Secret not found",
						"Namespace", fluxmore.Namespace,
						"Name", fluxmore.Spec.ResourcesCheck)
				} else {

					return ctrl.Result{}, err
				}
			} else {
				// Secret exists
				fluxmore.Status.ResourcesCheckFound = true
				l.Info("Resource found",
					"Namespace", fluxmore.Namespace,
					"Name", fluxmore.Spec.ResourcesCheck)
			}
			fluxmore.Status.LastReconcileTime = &timeReconcile

			if err := r.Status().Patch(ctx, &fluxmore, statusPatch); err != nil {
				l.Error(err, "Unable to update Fluxmore status")
				return ctrl.Result{}, err
			}
		} else {
			l.Info("Nothing to do with the HelmRelease, already patched", "Name", fluxmore.Spec.HelmReleaseName)
		}

	case "deployment":
		var Deploy appsv1.Deployment

		if helmrelease.Spec.Suspend {
			err := r.Get(ctx, types.NamespacedName{Name: fluxmore.Spec.ResourcesCheck, Namespace: fluxmore.Namespace}, &Deploy)
			// we use MergeForm instead of DeepCopy only because we are doing a strategic merge
			statusPatch := client.MergeFrom(fluxmore.DeepCopy())
			timeReconcile := metav1.NewTime(time.Now())

			if err != nil {
				if errors.IsNotFound(err) {
					// Pod not found
					fluxmore.Status.ResourcesCheckFound = false
					l.Info("Pod not found",
						"Namespace", fluxmore.Namespace,
						"Name", fluxmore.Spec.ResourcesCheck)
				} else {

					return ctrl.Result{}, err
				}
			} else {
				fmt.Println(Deploy.Status.AvailableReplicas)
				if Deploy.Status.ReadyReplicas != *fluxmore.Spec.ReplicaNumber {
					fluxmore.Status.ResourcesCheckFound = false
					l.Info("Pods found in a running phase but the expected replicas does not match",
						"Namespace", fluxmore.Namespace,
						"Name", fluxmore.Spec.ResourcesCheck)

				} else {
					// No error it means Pod exists
					fluxmore.Status.ResourcesCheckFound = true
					l.Info("Deployment found  and all pods are ready",
						"Namespace", fluxmore.Namespace,
						"Name", fluxmore.Spec.ResourcesCheck)
				}

			}
			fluxmore.Status.LastReconcileTime = &timeReconcile

			if err := r.Status().Patch(ctx, &fluxmore, statusPatch); err != nil {
				l.Error(err, "Unable to update Fluxmore status")
				return ctrl.Result{}, err
			}
		} else {
			l.Info("Nothing to do with the HelmRelease, already patched", "Name", fluxmore.Spec.HelmReleaseName)
		}

	case "statefulset":
		var Stateful appsv1.StatefulSet

		if helmrelease.Spec.Suspend {
			err := r.Get(ctx, types.NamespacedName{Name: fluxmore.Spec.ResourcesCheck, Namespace: fluxmore.Namespace}, &Stateful)
			// we use MergeForm instead of DeepCopy only because we are doing a strategic merge
			statusPatch := client.MergeFrom(fluxmore.DeepCopy())
			timeReconcile := metav1.NewTime(time.Now())

			if err != nil {
				if errors.IsNotFound(err) {
					// Pod not found
					fluxmore.Status.ResourcesCheckFound = false
					l.Info("Pod not found",
						"Namespace", fluxmore.Namespace,
						"Name", fluxmore.Spec.ResourcesCheck)
				} else {

					return ctrl.Result{}, err
				}
			} else {
				fmt.Println(Stateful.Status.AvailableReplicas)
				if Stateful.Status.ReadyReplicas != *fluxmore.Spec.ReplicaNumber {
					fluxmore.Status.ResourcesCheckFound = false
					l.Info("Pods found in a running phase but the expected replicas does not match",
						"Namespace", fluxmore.Namespace,
						"Name", fluxmore.Spec.ResourcesCheck)

				} else {
					// No error it means Pod exists
					fluxmore.Status.ResourcesCheckFound = true
					l.Info("Statefulset found  and all pods are ready",
						"Namespace", fluxmore.Namespace,
						"Name", fluxmore.Spec.ResourcesCheck)
				}

			}
			fluxmore.Status.LastReconcileTime = &timeReconcile

			if err := r.Status().Patch(ctx, &fluxmore, statusPatch); err != nil {
				l.Error(err, "Unable to update Fluxmore status")
				return ctrl.Result{}, err
			}
		} else {
			l.Info("Nothing to do with the HelmRelease, already patched", "Name", fluxmore.Spec.HelmReleaseName)
		}

	}
	return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FluxMoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&fluxitv1alpha.FluxMore{}).
		Named("fluxmore").
		Owns(&fluxitv1alpha.FluxMore{}).
		Complete(r)
}
