#!/bin/bash


mkdir -p ~/.aws
echo "[default]" >> ~/.aws/credentials
echo "aws_access_key_id = $AWS_S3_ACCESS_KEY" >> ~/.aws/credentials
echo "aws_secret_access_key = $AWS_S3_SECRET_KEY" >> ~/.aws/credentials

/go/bin/aem-s3-logsync