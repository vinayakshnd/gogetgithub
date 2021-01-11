# gogetgithub

*gogetgithub* does the following things
1. Connects to Github 
2. Connects to a repo 
3. Creates a branch
4. Create a Pull request on the newly created branch for the repo by modifying some of the files

## Localy compile the code

```
make build
```
## Build docker container

```
make container
```
## Build and push docker container to docker repository

```
make deploy
```
## Run docker container

```
docker run  -e CLIENT_ID=<GITHUB_OAUTH_APP_CLIENT_ID> -e CLIENT_SECRET=<GITHUB_OAUTH_APP_CLIENT_SECRET> -p 20080:8080 vinayakinfrac/gogetgithub:latest
```