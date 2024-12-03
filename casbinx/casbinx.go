/*
*
基于gorm适配器的rbac casbin
*/
package casbinx

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type CasbinxGorm struct {
	Casbin *casbin.Enforcer
}

func NewCasbinGorm(db *gorm.DB) (*CasbinxGorm, error) {
	m, err := model.NewModelFromString(rbacmodel)
	if err != nil {
		return nil, err
	}
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}
	e.LoadPolicy()
	return &CasbinxGorm{e}, nil
}

// 批量添加角色权限 旧角色会被删除
func (this *CasbinxGorm) AddPolicy(roleName string, policy []ApiPolice) (bool, error) {
	policies := [][]string{}
	for _, v := range policy {
		policies = append(policies, []string{roleName, v.Api, v.Method})
	}
	_, err := this.Casbin.RemoveFilteredPolicy(0, roleName)
	if err != nil {
		return false, err
	}
	b, err := this.Casbin.AddNamedPolicies("p", policies)
	if err != nil {

		return b, err
	}
	return b, err
}

// 删除角色权限
func (this *CasbinxGorm) RemoveRolePolicy(roleName string) (bool, error) {
	return this.Casbin.RemoveFilteredPolicy(0, roleName)
}

// 删除所有的角色
func (this *CasbinxGorm) Clear() {
	this.Casbin.ClearPolicy()
	this.Casbin.SavePolicy()
}

// 校验
func (this *CasbinxGorm) Enforce(role, api, method string) (bool, error) {
	return this.Casbin.Enforce(role, api, method)
}
