# Single-page Application with OIDC login using backend-for-frontend (BFF)

```
git clone https://github.com/michaelvl/oidc-oauth2-workshop.git
git clone https://github.com/michaelvl/oidc-oauth2-bff.git
git clone https://github.com/michaelvl/oidc-oauth2-workshop.git
```

```
make -C oidc-bff-apigw-workshop create-cluster deploy-metallb deploy-istio-base
```

```
kubectl apply -f oidc-oauth2-workshop/kubernetes/identity-provider-gateway.yaml

export IDENTITY_PROVIDER_GATEWAY_IP=$(kubectl get gateway idp -o jsonpath='{.status.addresses[0].value}')
echo "Identity provider IP: $IDENTITY_PROVIDER_GATEWAY_IP"

cat oidc-oauth2-workshop/kubernetes/identity-provider.yaml | envsubst | kubectl apply -f -
```

```
kubectl apply -f oidc-bff-apigw-workshop/kubernetes/gateway-httproutes.yaml
```

```
export SPA_GATEWAY_IP=$(kubectl get gateway spa -o jsonpath='{.status.addresses[0].value}')
echo "SPA IP: $SPA_GATEWAY_IP"
```


```
kubectl apply -f oidc-bff-apigw-workshop/kubernetes/spa-cdn.yaml
kubectl apply -f oidc-bff-apigw-workshop/kubernetes/spa-redis-session-store.yaml
cat oidc-bff-apigw-workshop/kubernetes/spa-login-bff.yaml | envsubst | kubectl apply -f -
```

```
stern -l app=spa-login-bff
```

> ![Initial SPA page](images/spa-pre-login.png)
> ![Click login brings us to the identity provider](images/idp-login.png)
> ![The identity provider ask us to authorize the clients use of our data](images/idp-authorize.png)
> ![Back to the SPA before loading userdata](images/spa-logged-in-initial.png)
> ![SPA with loaded userdata](images/spa-logged-in-with-userdata.png)
