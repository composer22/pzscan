### [Dockerized] (http://www.docker.com) [pzscan](https://registry.hub.docker.com/u/composer22/pzscan/)


A docker image for pzscan. This is created as a single "static" executable using a lightweight image.

To make:

cd docker
./build.sh

Once it completes, you can run the service thus:
```
docker run --name <containername> --rm composer22/pzscan -H <hostname>
```
for example:
```
docker run --name scantest --rm composer22/pzscan -H "example.com"
```

#### Options

For additional unix tools in a small image use "FROM busybox" instead of "FROM scratch" in Dockerfile.final
