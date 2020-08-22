# Pipeline: Atlan - Backend Engineering Intern - Task

Experimenting with pauseable, resumable tasks in a web service.

## Usage

### Using Docker

Run `docker run -it -p 8080:8080 prmsrswt/pipeline`. The API will be accessible on `http://localhost:8080`.

#### Building image yourself

- Build the image using the provided Dockerfile.

  ```
  docker build -t pipeline .
  ```

- Now you can run it using Docker.

  ```
  docker run -it -p 8080:8080 pipeline
  ```

- You can now access the API on `http://localhost:8080`.

### Building from source

- Just run `go get github.com/prmsrswt/pipeline`.

- The binary will now be available inside your `$GOPATH/bin` (`~/go/bin` in most cases). You can simply run it using `~/go/bin/pipeline`.

### Kubernetes

You can find example manifests in the `manifests/` directory. Modify them according to your needs and deploy using `kubectl create -f ./manifests`.

## API reference

#### `/upload` - Upload CSV file

| input  | description                         |
| ------ | ----------------------------------- |
| `file` | A CSV file which will get processed |

```bash
$ curl -X POST -F "file=@path/to/test.csv" http://localhost:8080/upload

{
  "status": "success",
  "data": {
    "id": "be9367c3-c492-4ce7-a256-cf4f21aa7b34"
  }
}
```

#### `/status` - Check status of a task

| input | description                                |
| ----- | ------------------------------------------ |
| `id`  | The task id of the task you want status of |

```bash
$ curl -X POST -F "id=2c78e760-1c0d-414e-99a4-3ba27b76c0f0" http://localhost:8080/status

{
  "status": "success",
  "data": {
    "status": "running"
  }
}
```

#### `/pause` - Pause a running task

| input | description                               |
| ----- | ----------------------------------------- |
| `id`  | The task id of the task you want to pause |

```bash
$ curl -X POST -F "id=edba118b-03db-4bbf-a94c-70f1992ff4f1" http://localhost:8080/pause

{
  "status": "success",
  "data": {
    "message": "task paused"
  }
}
```

#### `/resume` - Resume a paused task

| input | description                                |
| ----- | ------------------------------------------ |
| `id`  | The task id of the task you want to resume |

```bash
$ curl -X POST -F "id=edba118b-03db-4bbf-a94c-70f1992ff4f1" http://localhost:8080/resume

{
  "status": "success",
  "data": {
    "message": "task resumed"
  }
}
```

#### `/terminate` - terminate a running/paused task

| input | description                                   |
| ----- | --------------------------------------------- |
| `id`  | The task id of the task you want to terminate |

```bash
$ curl -X POST -F "id=edba118b-03db-4bbf-a94c-70f1992ff4f1" http://localhost:8080/terminate

{
  "status": "success",
  "data": {
    "message": "task terminated"
  }
}
```
