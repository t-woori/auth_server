aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com
docker tag auth-kakao:latest 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com/auth-kakao:latest
docker push 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com/auth-kakao:latest
echo "CD to AWS Lambda"
aws lambda update-function-code \
    --function-name auth_kakao \
    --image-uri 020759504372.dkr.ecr.ap-northeast-2.amazonaws.com/auth-kakao:latest \
    --no-cli-pager
echo "done"