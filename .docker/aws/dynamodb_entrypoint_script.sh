#!/bin/sh

exec /dynamodb/dynamodb-start.sh

if [ -z "${AWS_LAMBDA_RUNTIME_API}" ]; then
  exec /usr/local/bin/aws-lambda-rie ${LAMBDA_TASK_ROOT}/main $@
else
  exec ${LAMBDA_TASK_ROOT}/main $@
fi