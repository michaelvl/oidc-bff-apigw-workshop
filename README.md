# Single-page Application with OIDC login using backend-for-frontend (BFF)

```
git clone https://github.com/michaelvl/oidc-oauth2-workshop.git
git clone https://github.com/michaelvl/oidc-oauth2-bff.git
git clone https://github.com/michaelvl/oidc-oauth2-workshop.git
```

```
# oidc-bff-apigw-workshop
make create-cluster deploy-metallb deploy-istio-base
```

```
# oidc-oauth2-workshop
kubectl apply -f kubernetes/identity-provider.yaml
kubectl apply -f kubernetes/identity-provider-gateway.yaml
```

```
# oidc-bff-apigw-workshop
kubectl apply -f kubernetes/gateway-httproutes.yaml
```

```
export SPA_GATEWAY_IP=$(kubectl get gateway spa -o jsonpath='{.status.addresses[0].value}')
export IDENTITY_PROVIDER_GATEWAY_IP=$(kubectl get gateway idp -o jsonpath='{.status.addresses[0].value}')
echo "SPA IP: $SPA_GATEWAY_IP"
echo "Identity provider IP: $IDENTITY_PROVIDER_GATEWAY_IP"
```


```
# oidc-bff-apigw-workshop
kubectl apply -f kubernetes/spa-cdn.yaml
kubectl apply -f kubernetes/spa-redis-session-store.yaml
cat kubernetes/spa-login-bff.yaml | envsubst | kubectl apply -f -
```
