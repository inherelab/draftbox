package rbac

type CheckFunc func(id string, permission string, params interface{}) bool
