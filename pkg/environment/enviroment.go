package environment

type AppEnv string

const (
	AppEnvLocal      AppEnv = "local"
	AppEnvDevelop    AppEnv = "develop"
	AppEnvProduction AppEnv = "production"
	AppEnvTest       AppEnv = "test"
	AppEnvQA         AppEnv = "qa"
	AppEnvStage      AppEnv = "stage"
)

func (a AppEnv) IsLocal() bool {
	return a == AppEnvLocal
}

func (a AppEnv) IsDevelop() bool {
	return a == AppEnvDevelop
}

func (a AppEnv) IsProduction() bool {
	return a == AppEnvProduction
}

func (a AppEnv) IsTest() bool {
	return a == AppEnvTest
}

func (a AppEnv) IsQA() bool {
	return a == AppEnvQA
}

func (a AppEnv) IsStage() bool {
	return a == AppEnvStage
}
