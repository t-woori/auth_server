aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com
docker tag validate_token:latest 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com/validate_token:latest
docker push 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com/validate_token:latest
echo "CD to AWS Lambda"
aws lambda update-function-code \
    --function-name validate-token \
    --image-uri 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com/validate_token:latest \
    --no-cli-pager
echo "done"