package config

type Firebase struct {
	ProjectId  string `mapstructure:"projectId"`
	PrivateKey string `mapstructure:"privateKey"`
	Email      string `mapstructure:"email"`
}

var _ configLoader = (*Firebase)(nil)

func (Firebase) loadDefault() {
	//
}
