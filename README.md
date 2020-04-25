# k8s-toleration-injection-webhook

A k8s mutation webhook responsible for injection a toleration to all pods created in namespaces with a specific configurable label.

The toleration can be configured depending on the annotation.

## install

TODO: see how we can configure ssl properly, currently the image is packaging a certificate and can't really be reused.

`kubectl apply -f deploy/`

## use

For each type of toleration, define a webhook config like this:
```yaml
---
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  name: toleration-injection-for-application
  labels:
    app: toleration-injection-for-application
webhooks:
  - name: toleration-injection.kube-system.svc.cluster.local
    clientConfig:
      caBundle: $(cat ssl/toleration-injection.pem | base64)
      service:
        name: toleration-injection
        namespace: kube-system
        path: "/mutate/application"
    rules:
      - operations: ["CREATE"]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["pods"]
    sideEffects: None
    timeoutSeconds: 5
    failurePolicy: Fail
    namespaceSelector:
      matchLabels:
        type: application
```

The type need to be configured both in path `/mutate/{type}` and in the desired namespace selector.

Then define a configmap that defines the toleration for each type:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: toleration-config
  namespace: kube-system
data:
  TOLERATION_KEY_APPLICATION: dedicated
  TOLERATION_VALUE_APPLICATION: application
  TOLERATION_EFFECT_APPLICATION: NoSchedule
```

repeat `TOLERATION_(KEY|VALUE|EFFECT)_{upper(type)}` for each type.

## build 

```
make
```

## test

```
make test
```

## ssl/tls

the `ssl/` dir contains a script to create a self-signed certificate, not sure this will even work when running in k8s but that's part of figuring this out I guess

_NOTE: the app expects the cert/key to be in `ssl/` dir relative to where the app is running/started and currently is hardcoded to `mutateme.{key,pem}`_

```
cd ssl/ 
make 
```

## docker

to create a docker image .. 

```
make docker
```

it'll be tagged with the current git commit (short `ref`) and `:latest`

don't forget to update `IMAGE_PREFIX` in the Makefile or set it when running `make`


## kudos

- this is largely inspired by [alex-leonhardt/k8s-mutate-webhook](https://github.com/alex-leonhardt/k8s-mutate-webhook)
- other source of inspiration: [dotJobs/toleration-injection-webhook](https://github.com/dotJobs/toleration-injection-webhook)
- blog post [Writing a very basic kubernetes mutating admission webhook](https://medium.com/ovni/writing-a-very-basic-kubernetes-mutating-admission-webhook-398dbbcb63ec)  
- tutorial [https://github.com/morvencao/kube-mutating-webhook-tutorial](https://github.com/morvencao/kube-mutating-webhook-tutorial)
