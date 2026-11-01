---
title: Run a go service
sidebar_label: run a go service
---

We'll now look at an op to run a sample Go application. Please see the code at [https://github.com/opctl/opctl/tree/main/examples/run-a-go-service](https://github.com/opctl/opctl/tree/main/examples/run-a-go-service)
The sample application we have requires a mysql database, therefore to run it locally we need to:
1. Start a MySQL database with some DB credentials
2. Seed the DB with sample data
3. Start the application and provide the MySQL DB credentials as inputs - our application will read those from environment variables

The ops to make that happen are explained below.

#### mysql
```yaml
name: mysql
description: runs a mysql container, seeding it with sample data
inputs:
  MYSQL_USER:
    string:
      default: testuser
      description: username for MySQL user to create
  MYSQL_PASSWORD:
    string:
      default: testpassword
      description: password for MySQL user to create
  MYSQL_DATABASE:
    string:
      default: testapp
      description: name of mysql database to create
  MYSQL_HOST:
    string:
      default: run-a-go-service-mysql
run:
  container:
    dirs:
      # mount our data seed scripts
      /docker-entrypoint-initdb.d:
    image:
      ref: 'mysql:8'
    # this sets an opctl DNS A record so the containers available via this name.
    name: run-a-go-service-mysql
    envVars:
      MYSQL_USER:
      MYSQL_PASSWORD:
      MYSQL_DATABASE:
      MYSQL_RANDOM_ROOT_PASSWORD: "yes"
```
This YAML op runs a MySQL database in a container and seeds it with a single table holding a single row.

Note that instead of writing the `cmd` inline, we could have also put the run script in a file in the op directory and referenced it.
e.g.
```yaml
cmd: ["/run.sh"]
```
This becomes useful for readability if the script we're running in a container is large.

Note that while we could have also built a custom docker image that contains that script and ran it, that's not a recommended practice. While the op remains portable with that approach, it is less transparent and harder to maintain, since we would have the extra step of building and pushing that image whenever we make a change to the script. Instead, we usually opt for the leanest Docker image that fulfills our use case, and include any scripts or binaries we need in the op directory.

#### dev
The `dev` op is a [Go op](training/authoring-ops-in-go.md): its directory holds `op.go` and a `go.mod` instead of an op file.

```go
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
```

This op does the following, in parallel:
1. Calls the YAML `mysql` op by its path, passing the database inputs through.
2. Runs our service in a `golang` container. `OpDir("../..")` mounts the service's source, and `in.Home.Dir(...)` mounts the caller's Go module and build caches.

The doc comment becomes the op's description, and each field of `Inputs` declares an input. `MYSQL_HOST` defaults to the name of the `mysql` op's container, so the service can reach the database through opctl DNS.

Now to run the service locally, we only run `opctl run dev`

Notice how `Parallel` runs calls concurrently. Our app runs parallel to the MySQL DB container (via the `mysql` op), because we need mysql running while our app runs.
