echo "\nDeploying Alpha Stocks...\n"

envsubst < namespaces.yaml | kubectl apply -f -

kubectl create configmap dex-config --from-file=configs/dex-config.yaml --namespace alpha-stocks-dev

kubectl apply -f volumes/

kubectl apply -f service-accounts/

kubectl apply -f services/

kubectl apply -f deployments/dex-deployment.yaml
envsubst < ./deployments/stocks-deployment.yaml | kubectl apply -f -
envsubst < ./deployments/order-deployment.yaml | kubectl apply -f -
envsubst < ./deployments/gateway-deployment.yaml | kubectl apply -f -

if [ "$ENABLE_TRATS" = "true" ]; then
    kubectl apply -f tokenetes/kubernetes
    kubectl apply -f tokenetes-resources
fi

./wait_for_services.sh