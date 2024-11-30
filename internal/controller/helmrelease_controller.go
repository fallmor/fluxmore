package controller

import (
	"context"
	"time"

	fluxit "github.com/fallmor/fluxmore/api/v1alpha"
	helm "github.com/fluxcd/helm-controller/api/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type HelmReleaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=helm.toolkit.fluxcd.io,resources=helmreleases,verbs=get;list;watch;create;update;patch;delete

func (r *HelmReleaseReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	l := log.FromContext(ctx)

	// fluxmore := &fluxit.FluxMore{}
	// hrname := fluxmore.Spec.HelmReleaseName
	// hrnamespace := fluxmore.Namespace
	// ResourceName := fluxmore.Spec.ResourcesCheck
	// var CustomHr fluxitv1alpha.HelmReleaseList

	// helmRelease := fluxitv1alpha.HelmRelease{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      hrname,
	// 		Namespace: hrnamespace,
	// 	},
	// }

	var fluxMoreList fluxit.FluxMoreList
	if err := r.List(ctx, &fluxMoreList); err != nil {
		return ctrl.Result{}, err
	}

	for _, fluxMore := range fluxMoreList.Items {
		//if a secret/cm/pod is not found,exit the loop, go back and reconcile again
		if !fluxMore.Status.ResourcesCheckFound {
			continue
		}

		// Try to patch the associated HelmRelease
		var helmRelease helm.HelmRelease
		err := r.Get(ctx, types.NamespacedName{
			Name:      fluxMore.Spec.HelmReleaseName,
			Namespace: fluxMore.Namespace,
		}, &helmRelease)

		if err != nil {
			l.Error(err, "Unable to get HelmRelease",
				"Name", fluxMore.Spec.HelmReleaseName,
				"Namespace", fluxMore.Namespace)
			continue
		}
		if helmRelease.Spec.Suspend {
			// Prepare patch
			patch := client.MergeFrom(helmRelease.DeepCopy())

			helmRelease.Spec.Suspend = false

			// Apply patch
			if err := r.Patch(ctx, &helmRelease, patch); err != nil {
				l.Error(err, "Unable to patch HelmRelease",
					"Name", fluxMore.Spec.HelmReleaseName,
					"Namespace", fluxMore.Namespace)
				continue
			}

			// update the status of the Fluxmore
			statusPatch := client.MergeFrom(fluxMore.DeepCopy())
			fluxMore.Status.HelmReleasePatched = true
			if err := r.Status().Patch(ctx, &fluxMore, statusPatch); err != nil {
				l.Error(err, "Unable to update SecretMapper status")
			}

			l.Info("HelmRelease patched successfully by",
				"Resource", fluxMore.Spec.ResourcesKind,
				"Name", fluxMore.Spec.ResourcesCheck,
				"HelmRelease", fluxMore.Spec.HelmReleaseName)
		}
		l.Info("Nothing to do with the HelmRelease", "Name", fluxMore.Spec.HelmReleaseName)
	}

	return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil

}
func (r *HelmReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&helm.HelmRelease{}).
		Complete(r)
}
