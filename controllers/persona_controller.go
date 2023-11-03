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
	//"time" 

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	log "sigs.k8s.io/controller-runtime/pkg/log"

	compv2 "github.com/Aditya98Shukla/PersonaCRD/api/v2"
)

const AuthStateNeeded string = "AuthNeeded"
const OnlineState string = "Online"
const OfflineState string = "Offline"

// PersonaReconciler reconciles a Persona object
type PersonaReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=comp.genesis.xyz.com,resources=personas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=comp.genesis.xyz.com,resources=personas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=comp.genesis.xyz.com,resources=personas/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Persona object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *PersonaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Get the logger for logging within the reconciliation loop
	Log := log.FromContext(ctx)

	// Fetch the custom resource instance for the given request
	persona := &compv2.Persona{}
	err := r.Get(ctx, req.NamespacedName, persona)
	if err != nil {
		// Handle error, possibly return an error to indicate reconciliation failure
		Log.Error(err, "unable to fetch Persona")
		return ctrl.Result{Requeue: true}, client.IgnoreNotFound(err)
	}
	personCopy := persona.DeepCopy()

	// Add your custom logic here based on the state of the 'persona' object
	// For example, you can check the status, perform some actions, or update status.
	var personas = []compv2.PersonaSpec{
		{
			Name: "Michael Hudson",
			Age:  23,
		},
		{
			Name: "Wilston Roy",
			Age:  13,
		},
		{
			Name: "Lionel Messi",
			Age:  10,
		},
	}

	if persona.Status.State == "" {
		Log.Info("Persona is not active, performing some actions...")
		persona.Status.State = "AuthNeeded"
		persona.Status.Allowed = false
		persona.Status.ExpireDate = ""
	}

	if persona.Status.State != "" {
		found := false
		for _, obj := range personas {
			if persona.Spec.Name == obj.Name && persona.Spec.Age == obj.Age {
				found = true
				persona.Status.State = "Online"
				persona.Status.Allowed = true
				//persona.Status.ExpireDate = time.Now().AddDate(0, 0, 30).String()
				break
			}
		}

		// name and age does not match with table.
		if !found {
			persona.Status.Allowed = false
			persona.Status.State = "Offline"
			persona.Status.ExpireDate = ""
		}
	}

	if shouldUpdate(personCopy, persona) {
		Log.Info("Updating Persona Status")
		res, err := r.updateStatus(ctx, persona)
		if err != nil {
			return res, err
		}
	}
	// After completing your logic, return the result indicating whether the reconciliation is successful.
	// You can specify how long to wait before the next reconciliation by setting the RequeueAfter field.
	return ctrl.Result{}, nil
}

func shouldUpdate(copy *compv2.Persona, orig *compv2.Persona) (isModified bool) {
	isModified = false

	if orig.Status.State != copy.Status.State {
		isModified = true
	} else if orig.Status.Allowed != copy.Status.Allowed {
		isModified = true
	} else if orig.Status.ExpireDate != copy.Status.ExpireDate {
		isModified = true
	}
	return
}

func (r *PersonaReconciler) updateStatus(ctx context.Context, persona *compv2.Persona) (res ctrl.Result, err error) {
	err = r.Status().Update(ctx, persona)
	if err != nil {
		// Handle error, possibly return an error to indicate reconciliation failure
		log.FromContext(ctx).Error(err, "Unable to update Persona status")
		return ctrl.Result{Requeue: true}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PersonaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&compv2.Persona{}).
		Complete(r)
}
