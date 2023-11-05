aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com
docker tag token_regenerater:latest 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com/token_regenerater:latest
docker push 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com/token_regenerater:latest
echo "CD to AWS Lambda"
aws lambda update-function-code \
    --function-name token_regenerater \
    --image-uri 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com/token_regenerater:latest \
    --no-cli-pager
echo "done"