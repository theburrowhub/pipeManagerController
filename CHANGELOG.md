## 0.2.0 (2025-01-13)

### Feat

- WIP admission webhook
- **internal/k8s/normalize.go**: Add function ToK8sNames to normalize slice of Kubernetes resource names
- add functions for kubernetes objects
- **controller**: config file flag

### Fix

- **internal/runpipeline/tekton/pipeline_run.go**: Add build tasks params function
- **internal/normalize/normalize.go**: Add fixRunAfterReferencesForBatchTasks function
- **internal/normalize/batch.go**: Implement runAfter reference adjustments for batch tasks
- **internal/normalize/normalize.go**: avoid proccess twice batch tasks
- **internal/normalize/launch_pipeline_task.go**: Update launch pipeline task with default volumes and volume mounts
- **internal/normalize/normalize.go**: Update pipeline data processing for consistency

### Refactor

- **internal/validate/validate.go, internal/webhook/v1alpha1/pipeline_webhook.go**: Introduce validation package for webhook admission controller
- **config/default/kustomization.yaml**: Enable and configure cert-manager with injected CA for webhooks
- **config/manager/kustomization.yaml, config/manager/manager.yaml**: Update image tag and add volume mounts for config file in deployment YAML
- **internal/runpipeline/tekton/pipeline_run.go**: Remove unnecessary imports and packages for encoding/json, fmt, os, k8s.io/api/core/v1, k8s.io/apimachinery/pkg/apis/meta/v1, k8s.io/apimachinery/pkg/selection and github.com/tektoncd/pipeline/pkg/apis/pipeline/v1
- **internal/controller/pipeline_controller.go**: Add Tekton pipelineruns RBAC and update Reconcile comments for future reference
- **config/crd/kustomization.yaml**: Update CRD kustomization configuration
- **config/rbac/role.yaml**: Add tekton.dev group permissions to manager-role
- **config/config_example.yaml**: Update S3 artifacts bucket URL and credentials configuration
- **internal/webhook/v1alpha1/pipeline_webhook.go**: Remove unused imports, reorder and format imports, update webhook configurations, and move unreachable code up to make the file more readable and maintainable.
- **internal/normalize/functions.go**: Import K8s API and organize code for better readability
- **internal/runpipeline/tekton/pipeline_run.go**: Update import and delete unused library
- **internal/runpipeline/tekton/pipeline_run.go**: Extract constants and variables for easier code maintenance
- **internal/runpipeline/run_normalized_pipeline.go**: Update pipeline deployment logic to use a map for pipeline types and functions
- **go.mod**: Update go.mod dependencies and version tags
- **internal/normalize/volumes.go**: Remove unnecessary variable declarations for volumes and volumeMounts
- **internal/runpipeline/tekton/pipeline_run.go**: retrieve workingDir from normalize package
- A temporal json serialize pipeline for debbuging
- **internal/normalize/launchpipelinetask.go, internal/normalize/launch_pipeline_task.go**: Rename and move LaunchPipelineTask file for consistency
- **Makefile**: Remove TODO comment and update build process
- **cmd/main.go**: Remove Kubernetes scaffold imports and setup logger configuration
- **config/config_example.yaml**: Remove unnecessary webhook configuration
- **go.mod**: Update dependencies and remove unused module
- **internal/normalize/artifacts.go**: Remove usage of envvars package and update defineDownloadArtifactsStep function
- **kubernetes**: Update pipeline launching process by appending clone repository step
- **go.sum**: Update dependencies to latest versions (sigs.k8s.io/controller-runtime, sigs.k8s.io/json, sigs.k8s.io/structured-merge-diff/v4, sigs.k8s.io/yaml)

## 0.1.0 (2024-11-17)

### Feat

- **config/crd/bases/pipemanager.sergiotejon.github.io_pipelines.yaml**: Initialize new Custom Resource Definition (CRD) file for PipelineManager resources
- **internal/controller/pipeline_controller.go**: Add Pipeline resource handling logic, refactor code structure
- **api/v1alpha1/pipeline_types.go**: Extend Pipeline structure with Active field

### Refactor

- **pipeline**: Improve Task object's DeepCopy method
- **api/v1alpha1/pipeline_types.go**: Update schema builder to include Pipeline resources
- **api/v1alpha1/groupversion_info.go**: Update GroupVersionInfo schema, including scheme initialization and SchemeBuilder adjustments
- **pipeline**: Implement deep copy function for pipeline components
- **io/k8s.apiextensions.v1**: Simplify Pipeline custom resource definition
- **config/rbac/role.yaml**: Update cluster role permissions for pipeline manager role
