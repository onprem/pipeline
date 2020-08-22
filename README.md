# Pipeline

Experimenting with pauseable, resumable tasks in a web service.

## Usage

### Using Docker

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

## API reference

#### `/upload` - Upload CSV file

| input  | description                         |
| ------ | ----------------------------------- |
| `file` | A CSV file which will get processed |

#### `/status` - Check status of a task

| input | description                                |
| ----- | ------------------------------------------ |
| `id`  | The task id of the task you want status of |

#### `/pause` - Pause a running task

| input | description                               |
| ----- | ----------------------------------------- |
| `id`  | The task id of the task you want to pause |

#### `/resume` - Resume a paused task

| input | description                                |
| ----- | ------------------------------------------ |
| `id`  | The task id of the task you want to resume |

#### `/terminate` - terminate a running/paused task

| input | description                                   |
| ----- | --------------------------------------------- |
| `id`  | The task id of the task you want to terminate |
