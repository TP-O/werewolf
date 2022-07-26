package factory

var roleFactoryInstance *roleFactory

func init() {
	roleFactoryInstance = &roleFactory{}
}

func GetRoleFactory() *roleFactory {
	return roleFactoryInstance
}
