# Lazy Panda User Service
Lazy Panda is a project aimed on managing employee time for consulting companies.

## What this is
- A service to create users for **Lazy Panda** project using Go, [GRPC](https://github.com/grpc/grpc-go)
- the service app has a GO / GRPC service for creating, updating, deleting, and reading users
- it's using 🔥 Firebase for saving data but can use any database that extends the [database package from utils repo](https://github.com/omaressameldin/lazy-panda-utils/tree/master/app/pkg/database)

## How to run
- make sure you have **docker version: 18.x+** installed
- create a firebase project and install json config [follow this tutorial to get firebase config](https://www.youtube.com/watch?v=9rN29jENirI)
- rename the config to `firebaseConfig.json` and add it to the root of the project
- run `docker-compose up --build` to launch service
- the service will be available at port `7500`

## Taking the service for a spin
**Note1:** Please, make sure that the service is running before testing any of the following snippets

**Note2:** This service is meant to run alongside google sign in that's why authid is not autogenerated but rather given

- you can check `./proto-gen/proto/v1/user.proto` file for available fields

- To connect to service:
```golang
import (
	v1 "github.com/omaressameldin/lazy-panda-user-service/pkg/api/v1"
	"google.golang.org/grpc"
)
func main() {
	client, err := grpc.Dial("localhost:7500", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err)
	}

	c := v1.NewUserServiceClient(client)
}
```

- To create a user:
```golang
import (
	"context"
	"log"

	v1 "github.com/omaressameldin/lazy-panda-user-service/pkg/api/v1"
	"google.golang.org/grpc"
)
func main() {
  // connecting to service
	client, err := grpc.Dial("localhost:7500", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err)
	}
  c := v1.NewUserServiceClient(client)

  // create request
	req := v1.CreateRequest{
		Api: "v1",
		User: &v1.User{
			AuthId:   "1",
			Email:    "newuser@gmail.com",
			Fullname: "New user",
			Nickname: "newuser",
		},
	}
  // create response
	res, err := c.Create(context.Background(), &req)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
  }
  // log response
	log.Printf("create: %v", res)
}
```

- To get a user:
```golang
import (
	"context"
	"log"

	v1 "github.com/omaressameldin/lazy-panda-user-service/pkg/api/v1"
	"google.golang.org/grpc"
)
func main() {
  // connecting to service
	client, err := grpc.Dial("localhost:7500", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err)
	}
  c := v1.NewUserServiceClient(client)

  // read request
	req := v1.ReadRequest{
		Api:    "v1",
		AuthId: "1",
  }
  // read response
	res, err := c.Read(context.Background(), &req)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
  }
  // log response
	log.Printf("read: %v", res)
}
```

- To update a user:
```golang
import (
	"context"
	"log"
	"github.com/golang/protobuf/ptypes/wrappers"

	v1 "github.com/omaressameldin/lazy-panda-user-service/pkg/api/v1"
	"google.golang.org/grpc"
)
func main() {
  // connecting to service
	client, err := grpc.Dial("localhost:7500", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err)
	}
  c := v1.NewUserServiceClient(client)

  // update request
	req := v1.UpdateRequest{
		Api:    "v1",
		AuthId: "1",
		User: &v1.UserUpdate{
			Email:    &wrappers.StringValue{Value: "updateduser@gmail.com"},
			Fullname: &wrappers.StringValue{Value: "Updated User"},
			Nickname: &wrappers.StringValue{Value: "updateduser"},
		},
	}
  //update response
	res, err := c.Update(context.Background(), &req)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
  }
  // log response
	log.Printf("update: %v", res)
}
```

- To read all users:
```golang
import (
	"context"
	"log"

	v1 "github.com/omaressameldin/lazy-panda-user-service/pkg/api/v1"
	"google.golang.org/grpc"
)
func main() {
  // connecting to service
	client, err := grpc.Dial("localhost:7500", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err)
	}
  c := v1.NewUserServiceClient(client)

  // read all request
	req := v1.ReadAllRequest{
		Api: "v1",
  }
  // read all response
	res, err := c.ReadAll(context.Background(), &req)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
  }
  // log response
	log.Printf("read all: %v", res)
}
```

- To delete a user:
```golang
import (
	"context"
	"log"

	v1 "github.com/omaressameldin/lazy-panda-user-service/pkg/api/v1"
	"google.golang.org/grpc"
)
func main() {
  // connecting to service
	client, err := grpc.Dial("localhost:7500", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err)
	}
  c := v1.NewUserServiceClient(client)

  // delete request
	req := v1.DeleteRequest{
		Api:    "v1",
		AuthId: "1",
  }
  // delete response
	res, err := c.Delete(context.Background(), &req)
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
  }
  // log response
	log.Printf("delete: %v", res)
}
```

## Technologies used
- Golang
- GRPC
- firebase
- Docker
- Docker-compose