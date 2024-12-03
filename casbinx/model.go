package casbinx

var rbacmodel = `
[request_definition]
r = role, api, method

[policy_definition]
p = role, api, method

[matchers]
m = r.role == p.role && keyMatch2(r.api,p.api) && r.method==p.method

[policy_effect]
e = some(where (p.eft == allow))
`

type ApiPolice struct {
	Api    string
	Method string
}
