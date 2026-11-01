// Runs the example API with a seeded MySQL database for development.
package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	Home             opctl.Dir    `opctl:"HOME,description='Home directory of the caller; mounts its Go module and build caches'"`
	GoBuildCachePath opctl.String `opctl:"goBuildCachePath,description='Go build cache path relative to HOME, for example .cache/go-build on Linux',default=Library/Caches/go-build"`
	User             opctl.String `opctl:"MYSQL_USER,description='Username of the MySQL user to create',default=testuser"`
	Password         opctl.String `opctl:"MYSQL_PASSWORD,description='Password of the MySQL user to create',default=testpassword"`
	Database         opctl.String `opctl:"MYSQL_DATABASE,description='Database to create',default=testapp"`
	Host             opctl.String `opctl:"MYSQL_HOST,description='Name of the MySQL container',default=run-a-go-service-mysql"`
}

func plan(in Inputs) opctl.Call {
	db := map[string]any{"MYSQL_USER": in.User, "MYSQL_PASSWORD": in.Password, "MYSQL_DATABASE": in.Database, "MYSQL_HOST": in.Host}
	api := opctl.Container("golang:1.25").
		// Other containers reach the API at this name through opctl DNS.
		Name("run-a-go-service-api").
		Dir("/src", opctl.OpDir("../..")).
		Dir("/go/pkg", in.Home.Dir("go/pkg")).
		Dir("/root/.cache/go-build", in.Home.Dir(in.GoBuildCachePath.Value())).
		EnvFrom("MYSQL_USER", in.User).
		EnvFrom("MYSQL_PASSWORD", in.Password).
		EnvFrom("MYSQL_DATABASE", in.Database).
		EnvFrom("MYSQL_HOST", in.Host).
		WorkDir("/src").
		Cmd("go", "run", "main.go")
	// The YAML mysql op runs the database while the API runs.
	return opctl.Parallel(opctl.CallOp("../mysql").Inputs(db), api)
}

func main() { opctl.Main(plan) }
