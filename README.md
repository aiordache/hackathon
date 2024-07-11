# Docker Compose Hub

## Objective
Create a solution that allows a Docker Compose file to be stored on Docker Hub and deployed using docker run. The docker run command should automatically detect whether the target is a regular container image or a Docker Compose file. If a Compose file is detected, it should be extracted and launched automatically using Docker Compose.


### Build
 ```
 $ make
 ```

### Push the compose image to Docker Hub
Create a docker-compose.yaml for an application (all services should specify the images they rely on). An example compose file can be found in the example directory.
```
$ ./bin/docker.exe v2 build -t <repo/name:tag> example/
```

### Run the docker compose app 
```
$ ./bin/docker.exe v2 run <repo/name:tag>
```
