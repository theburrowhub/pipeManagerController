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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/sergiotejon/pipeManagerController/internal/normalize"
	"github.com/sergiotejon/pipeManagerController/internal/runners"

	pipemanagerv1alpha1 "github.com/sergiotejon/pipeManagerController/api/v1alpha1"
	tektonpipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// PipelineReconciler reconciles a Pipeline object
type PipelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=pipemanager.sergiotejon.github.io,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=pipemanager.sergiotejon.github.io,resources=pipelines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=pipemanager.sergiotejon.github.io,resources=pipelines/finalizers,verbs=update

// Tekton pipelineruns
// +kubebuilder:rbac:groups=tekton.dev,resources=pipelineruns,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pipeline object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *PipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error

	logger := log.FromContext(ctx)

	// Fetch the Pipeline instance
	var pipeline pipemanagerv1alpha1.Pipeline
	if err = r.Get(ctx, req.NamespacedName, &pipeline); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Pipeline resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get pipeline.")
		return ctrl.Result{}, err
	}

	pipelineRunner := runners.GetRunner()

	// Find the runners for the pipeline
	err = r.List(ctx, pipelineRunner.List, client.InNamespace(pipeline.Namespace), client.MatchingLabels{
		"pipelineRef": string(pipeline.GetUID()),
	})
	if err != nil {
		logger.Error(err, "Failed to list runners objects")
		return ctrl.Result{}, err
	}
	// If the runners do not exist, create it
	if len(pipelineRunner.List.Items) == 0 {
		var normalizedPipelineSpec pipemanagerv1alpha1.PipelineSpec

		// Normalize the pipeline
		normalizedPipelineSpec, err = normalize.Normalize(logger, pipeline.Spec)
		if err != nil {
			logger.Error(err, "Error normalizing pipeline")
		}

		// Get the runners for the pipeline type
		err = pipelineRunner.BuildPipeline(string(pipeline.GetUID()), &normalizedPipelineSpec)
		if err != nil {
			logger.Error(err, "Error creating runners")
		}

		// Set the pipeline as the owner of the object created to run the pipeline.
		// Depending on the pipeline type, the object created can be any kind of Kubernetes object
		err = controllerutil.SetControllerReference(&pipeline, pipelineRunner.Object, r.Scheme)
		if err != nil {
			logger.Error(err, "Error setting controller reference")
		}
		// Set the pipeline as the owner of the object created to run the pipeline.
		logger.Info("Creating new runner", "UID", string(pipeline.GetUID()))
		if err = r.Create(ctx, pipelineRunner.Object); err != nil {
			logger.Error(err, "Failed to create new runner object")
			return ctrl.Result{}, err
		}
		// Requeue the request to ensure the Deployment is created
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// TODO: Update status of the pipeline if needed
	// pushMain.Status.SomeField = "SomeValue"
	// if err := r.Status().Update(ctx, &pushMain); err != nil {
	//     logger.Error(err, "Failed to update PushMain status")
	//     return ctrl.Result{}, err
	// }

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// TODO: When new pipeline types are added, the controller should be able to manage them
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipemanagerv1alpha1.Pipeline{}).
		Owns(&tektonpipelinev1.PipelineRun{}).
		Named("pipeline").
		Complete(r)
}
