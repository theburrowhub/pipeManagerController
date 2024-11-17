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
