@echo off
setlocal enabledelayedexpansion

set "DEBUG_CONTAINERS="
set "WARM_CONTAINERS="
set "REBUILD_IMAGES="
set "REMOVE_IMAGES=false"

:: Verifique os argumentos passados ao script
:parse_args
if "%~1"=="" goto end_args
  if "%~1"=="--debug" (
    set "DEBUG_CONTAINERS=--debug"
  ) else if "%~1"=="--warm" (
    set "WARM_CONTAINERS=--warm-containers eager"
  ) else if "%~1"=="--rebuild" (
    set "REBUILD_IMAGES=--no-cache"
    set "REMOVE_IMAGES=true"
  ) else (
    echo Argumento desconhecido: %1
    exit /b 1
  )
  shift
  goto parse_args

:end_args

:: Remova containers existentes
for /f %%i in ('docker ps -aq') do (
  docker rm -f %%i
)

:: Remova imagens, se necessário
if "%REMOVE_IMAGES%"=="true" (
  docker rmi -f authorizer.ecr.aws/lambda/go:1.x
  docker build -f .docker/docker-authorizer -t authorizer.ecr.aws/lambda/go:1.x . %REBUILD_IMAGES%

  docker rmi -f database.ecr.aws/dynamodb/java:latest
  docker build -f .docker/docker-dynamodb -t database.ecr.aws/dynamodb/java:latest . %REBUILD_IMAGES%

  docker rmi -f database.ecr.aws/lambda/go:1.x
  docker build -f .docker/docker-database -t database.ecr.aws/lambda/go:1.x . %REBUILD_IMAGES%

  docker rmi -f hello.ecr.aws/lambda/go:1.x
  docker build -f .docker/docker-hello -t hello.ecr.aws/lambda/go:1.x . %REBUILD_IMAGES%

  docker rmi -f token.ecr.aws/lambda/go:1.x
  docker build -f .docker/docker-token -t token.ecr.aws/lambda/go:1.x . %REBUILD_IMAGES%
)

:: Inicialize o container do DynamoDB Local
docker run -d -p 18000:8000 database.ecr.aws/dynamodb/java:latest

:: Execute o SAM Local com base nos argumentos
sam local start-api -t template.yaml --skip-pull-image %DEBUG_CONTAINERS% %WARM_CONTAINERS%
