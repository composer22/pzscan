#!/bin/bash
docker build -t composer22/pzscan_build .
docker run -v /var/run/docker.sock:/var/run/docker.sock -v $(which docker):$(which docker) -ti --name pzscan_build composer22/pzscan_build
docker rm pzscan_build
docker rmi composer22/pzscan_build
