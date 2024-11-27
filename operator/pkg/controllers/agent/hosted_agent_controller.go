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

package agent

import (
	"context"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"open-cluster-management.io/api/addon/v1alpha1"
	addonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	globalhubv1alpha4 "github.com/stolostron/multicluster-global-hub/operator/api/operator/v1alpha4"
	"github.com/stolostron/multicluster-global-hub/operator/pkg/config"
	operatorconstants "github.com/stolostron/multicluster-global-hub/operator/pkg/constants"
	"github.com/stolostron/multicluster-global-hub/pkg/constants"
	"github.com/stolostron/multicluster-global-hub/pkg/utils"
)

// +kubebuilder:rbac:groups=operator.open-cluster-management.io,resources=multiclusterglobalhubs,verbs=get;list;watch;
// +kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons,verbs=create;delete;get;list;patch;update;watch
// +kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=addondeploymentconfigs,verbs=get;list;watch
// +kubebuilder:rbac:groups=addon.open-cluster-management.io,resources=clustermanagementaddons,verbs=get;list;watch

type HostedAgentController struct {
	c client.Client
}

var (
	hostedAddonController        *HostedAgentController
	isHostedAgentResourceRemoved = true
	hostedAgentController        *HostedAgentController
)

func StartHostedAgentController(initOption config.ControllerOption) (config.ControllerInterface, error) {
	if hostedAgentController != nil {
		return hostedAgentController, nil
	}
	if !config.GetImportClusterInHosted() {
		return nil, nil
	}
	if !ReadyToEnableAddonManager(initOption.MulticlusterGlobalHub) {
		return nil, nil
	}

	hostedAgentController = NewHostedAgentController(initOption.Manager)

	if err := hostedAgentController.SetupWithManager(initOption.Manager); err != nil {
		hostedAgentController = nil
		return nil, err
	}
	return hostedAgentController, nil
}

func (c *HostedAgentController) IsResourceRemoved() bool {
	log.Infof("HostedAgentController resource removed: %v", isHostedAgentResourceRemoved)
	return isHostedAgentResourceRemoved
}

func NewHostedAgentController(mgr ctrl.Manager) *HostedAgentController {
	return &HostedAgentController{c: mgr.GetClient()}
}

// SetupWithManager sets up the controller with the Manager.
func (r *HostedAgentController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).Named("AddonsController").
		For(&v1alpha1.ClusterManagementAddOn{},
			builder.WithPredicates(addonPred)).
		// requeue all cma when mgh annotation changed.
		Watches(&globalhubv1alpha4.MulticlusterGlobalHub{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				var requests []reconcile.Request
				for v := range config.HostedAddonList {
					request := reconcile.Request{
						NamespacedName: types.NamespacedName{
							Name: v,
						},
					}
					requests = append(requests, request)
				}
				return requests
			}), builder.WithPredicates(mghPred)).
		Complete(r)
}

var addonPred = predicate.Funcs{
	CreateFunc: func(e event.CreateEvent) bool {
		return config.HostedAddonList.Has(e.Object.GetName())
	},
	UpdateFunc: func(e event.UpdateEvent) bool {
		return config.HostedAddonList.Has(e.ObjectNew.GetName())
	},
	DeleteFunc: func(e event.DeleteEvent) bool {
		return false
	},
}

var mghPred = predicate.Funcs{
	CreateFunc: func(e event.CreateEvent) bool {
		return true
	},
	UpdateFunc: func(e event.UpdateEvent) bool {
		if reflect.DeepEqual(e.ObjectNew.GetAnnotations(), e.ObjectOld.GetAnnotations()) {
			return false
		}
		if e.ObjectNew.GetAnnotations()[operatorconstants.AnnotationImportClusterInHosted] !=
			e.ObjectOld.GetAnnotations()[operatorconstants.AnnotationImportClusterInHosted] {
			return true
		}
		return false
	},
	DeleteFunc: func(e event.DeleteEvent) bool {
		return false
	},
}

func (r *HostedAgentController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Debug("Reconcile ClusterManagementAddOn: %v", req.NamespacedName)
	mgh, err := config.GetMulticlusterGlobalHub(ctx, r.c)
	if err != nil {
		log.Error(err)
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}
	if mgh == nil || config.IsPaused(mgh) {
		return ctrl.Result{}, nil
	}
	if mgh.DeletionTimestamp != nil || !config.GetImportClusterInHosted() {
		if config.GetImportClusterInHosted() {
			hasmanagedHub, err := r.hasManagedHub(ctx)
			if err != nil {
				return ctrl.Result{}, err
			}
			if hasmanagedHub {
				log.Errorf("You need to detach all the managed hub clusters before uninstalling")
				return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
			}
		}
		err = r.revertClusterManagementAddon(ctx)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = r.pruneHostedResources(ctx)
		if err != nil {
			return ctrl.Result{}, err
		}
		isHostedAgentResourceRemoved = true
		return ctrl.Result{}, nil
	}

	isHostedAgentResourceRemoved = false
	cma := &v1alpha1.ClusterManagementAddOn{}
	err = r.c.Get(ctx, req.NamespacedName, cma)
	if err != nil {
		return ctrl.Result{}, err
	}

	needUpdate := addAddonConfig(cma)
	if !needUpdate {
		return ctrl.Result{}, nil
	}

	err = r.c.Update(ctx, cma)
	if err != nil {
		log.Errorf("Failed to update cma, err:%v", err)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *HostedAgentController) revertClusterManagementAddon(ctx context.Context) error {
	cmaList := &v1alpha1.ClusterManagementAddOnList{}
	err := r.c.List(ctx, cmaList)
	if err != nil {
		return err
	}
	for _, cma := range cmaList.Items {
		if !config.HostedAddonList.Has(cma.Name) {
			continue
		}
		err = r.removeGlobalhubConfig(ctx, cma)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *HostedAgentController) removeGlobalhubConfig(ctx context.Context, cma v1alpha1.ClusterManagementAddOn) error {
	if len(cma.Spec.InstallStrategy.Placements) == 0 {
		return nil
	}
	var desiredPlacementStrategy []addonv1alpha1.PlacementStrategy
	exist := false
	for _, pl := range cma.Spec.InstallStrategy.Placements {
		if reflect.DeepEqual(pl.PlacementRef, config.GlobalHubHostedAddonPlacementStrategy.PlacementRef) {
			exist = true
			continue
		}
		desiredPlacementStrategy = append(desiredPlacementStrategy, pl)
	}
	if !exist {
		return nil
	}
	cma.Spec.InstallStrategy.Placements = desiredPlacementStrategy
	return r.c.Update(ctx, &cma)
}

func (r *HostedAgentController) pruneHostedResources(ctx context.Context) error {
	addonDeployConfig := &addonv1alpha1.AddOnDeploymentConfig{}
	if err := r.c.Get(ctx, types.NamespacedName{
		Namespace: utils.GetDefaultNamespace(),
		Name:      "global-hub",
	}, addonDeployConfig); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return r.c.Delete(ctx, addonDeployConfig)
}

func (r *HostedAgentController) hasManagedHub(ctx context.Context) (bool, error) {
	mcaList := &v1alpha1.ManagedClusterAddOnList{}
	err := r.c.List(ctx, mcaList)
	if err != nil {
		return false, err
	}
	for _, mca := range mcaList.Items {
		if mca.Name == constants.GHManagedClusterAddonName {
			return true, nil
		}
	}

	return false, nil
}

// addAddonConfig add the config to cma, will return true if the cma updated
func addAddonConfig(cma *v1alpha1.ClusterManagementAddOn) bool {
	if len(cma.Spec.InstallStrategy.Placements) == 0 {
		cma.Spec.InstallStrategy.Placements = append(cma.Spec.InstallStrategy.Placements,
			config.GlobalHubHostedAddonPlacementStrategy)
		return true
	}
	for i, pl := range cma.Spec.InstallStrategy.Placements {
		if !reflect.DeepEqual(pl.PlacementRef, config.GlobalHubHostedAddonPlacementStrategy.PlacementRef) {
			continue
		}
		if reflect.DeepEqual(pl.Configs, config.GlobalHubHostedAddonPlacementStrategy.Configs) {
			return false
		}
		cma.Spec.InstallStrategy.Placements[i].Configs = append(pl.Configs,
			config.GlobalHubHostedAddonPlacementStrategy.Configs...)
		return true
	}
	cma.Spec.InstallStrategy.Placements = append(cma.Spec.InstallStrategy.Placements,
		config.GlobalHubHostedAddonPlacementStrategy)
	return true
}
