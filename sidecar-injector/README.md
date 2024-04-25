# Kyverno Envoy Sidecar Injector 

Uses [MutatingAdmissionWebhook Controller](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook) in Kubernetes to inject kyverno-envoy-plugin sidecar into newly created pods. This injection occurs at pod creation time, targeting pods that have the label `kyverno-envoy-sidecar/injection=enabled`. By introducing this sidecar, we can enforce policies on all incoming HTTP requests and make external authorization decisions to the targeted pod without modifying the primary application code/containers.


## Prerequisites

Kubernetes 1.16.0 or above with the `admissionregistration.k8s.io/v1` API enabled.
Verify that by the following command:
```bash
 ~$ kubectl api-versions | grep admissionregistration.k8s.io/v1
```
The result should be:
```bash
admissionregistration.k8s.io/v1
```
## Installation  

#### Dedicated Namespace

Create a namespace `kyverno-envoy-sidecar-injector`, where you will deploy the Kyverno Envoy Sidecar Injector Webhook components.

```bash
 ~$ kubectl create namespace kyverno-envoy-sidecar-injector
```

#### Deploy Sidecar Injector

1. Create a signed cert/key pair and store it in a Kubernetes `secret` that will be consumed by sidecar injector deployment 
 
    Generate cert/key pair with openssl 
    ```bash
     ~$ openssl req -new -x509  \
        -subj "/CN=kyverno-envoy-sidecar.kyverno-envoy-sidecar-injector.svc" \
        -addext "subjectAltName = DNS:kyverno-envoy-sidecar.kyverno-envoy-sidecar-injector.svc" \
        -nodes -newkey rsa:4096 -keyout tls.key -out tls.crt 
    ```
    Now apply below command to create `secret`  
    ```bash
     ~$ kubectl create secret generic kyverno-envoy-sidecar-certs \
        --from-file tls.crt=tls.crt \
        --from-file tls.key=tls.key \
        --dry-run=client -n kyverno-envoy-sidecar-injector -oyaml > secret.yaml 
    ```  
    Apply the secret
    ```bash
    ~$ kubectl apply -f secret.yaml
    ```

2. Run the script to Patch the `Mutating Webhook Configuration` with the CA bundle extracted from the `secret` created in the previous step and apply the MutatingWebhookConfiguration changes: 
 
```bash
 ~$ ./manifests/create-mutating-webhook.sh
```
3. To Inject the Kyverno Envoy Sidecar, Create this configmap of name `kyverno-envoy-sidecar` in `kyverno-envoy--sidecar-injector` namespace. If their is requirement of multiple policy files, you can add more `--policy` flags and then add them in the `policy-files` configmap.

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: kyverno-envoy-sidecar
  namespace: kyverno-envoy-sidecar-injector
data:
  sidecars.yaml: |
    - name: kyverno-envoy-sidecar
      containers:
      - image: sanskardevops/plugin:0.0.25
        imagePullPolicy: IfNotPresent
        name: ext-authz
        ports:
        - containerPort: 8000
        - containerPort: 9000
        args:
        - "serve"
        - "--policy=/policies/policy.yaml"
        volumeMounts:
        - name: policy-files
          mountPath: /policies
      volumes:
      - name: policy-files
        configMap:
          name: policy-files      
EOF          
```

4. Deploy resources
```bash
 ~$ kubectl apply -f ./manifests/rbac.yaml
 ~$ kubectl apply -f ./manifests/deployment.yaml
 ~$ kubectl apply -f ./manifests/service.yaml
```

#### Verify Sidecar Injector Installation

1. The sidecar injector should be deployed in the `kyverno-envoy-sidecar-injector` namespace:

```bash
 ~$ kubectl -n kyverno-envoy-sidecar-injector get all
```
```bash
~$ kubectl -n  kyverno-envoy-sidecar-injector get all 
NAME                                        READY   STATUS    RESTARTS   AGE
pod/kyverno-envoy-sidecar-976c94445-2l66q   1/1     Running   0          46s

NAME                            TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
service/kyverno-envoy-sidecar   ClusterIP   10.96.137.93   <none>        443/TCP   3m11s

NAME                                    READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/kyverno-envoy-sidecar   1/1     1            1           46s

NAME                                              DESIRED   CURRENT   READY   AGE
replicaset.apps/kyverno-envoy-sidecar-976c94445   1         1         1       46s

```
2. Now create a pod with the label `kyverno-envoy-sidecar/injection=enabled` in any namespace other than `kyverno-envoy-sidecar-injector`. But before we have to apply configmap `policy-files` in the same namespace where we will create the pod.
```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: policy-files
  namespace: default
data:
  policy.yaml: |
    apiVersion: json.kyverno.io/v1alpha1
    kind: ValidatingPolicy
    metadata:
      name: check-dockerfile
    spec:
      rules:
        - name: deny-external-calls
          assert:
            all:
            - message: "HTTP calls are not allowed"
              check:
                request:
                    http:
                        method: GET
                        headers:
                            authorization:
                                (base64_decode(split(@, ' ')[1])): 
                                    (split(@, ':')[0]): alice
                        path: /foo    
EOF          
```

Now create a pod with the label `kyverno-envoy-sidecar/injection=enabled` in the default namespace:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  namespace: default
  labels:
    app.kubernetes.io/name: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: nginx
  template:
    metadata:
      labels:
        kyverno-envoy-sidecar/injection: enabled
        app.kubernetes.io/name: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.20.2
          ports:
            - containerPort: 80
```

Apply the following command to creat above example deployment:
```bash
kubectl apply -f ./example-manifest/exampledeploy.yaml
```

Now check the pods, you should see the Kyverno Envoy Sidecar injected:

```bash
 ~$ kubectl get pods 
```

Two pods are runing
```console 
~$ kubectl get pods 
NAME                     READY   STATUS    RESTARTS   AGE
nginx-77b746455b-xntjn   2/2     Running   0          3m46s
```

3. Check the logs of the Kyverno-envoy-sidecar container to verify the sidecar is running:

```bash
sanskar@sanskar-HP-Laptop-15s-du1xxx:~$ kubectl logs -n kyverno-envoy-sidecar-injector kyverno-envoy-sidecar-976c94445-nf777 -f 
time="2024-04-20T13:59:10Z" level=info msg="SimpleServer starting to listen in port 8443"
time="2024-04-20T14:03:32Z" level=info msg="AdmissionReview for Kind=/v1, Kind=Pod, Namespace=default Name= UID=a57c5c0b-96c0-4c9c-b903-6aa75f635c17 patchOperation=CREATE UserInfo={system:serviceaccount:kube-system:replicaset-controller 9e6576d2-f5c3-4b44-9b9d-952a20b70da7 [system:serviceaccounts system:serviceaccounts:kube-system system:authenticated] map[]}"
time="2024-04-20T14:03:32Z" level=info msg="sideCar injection for kyverno-envoy-sidecar-injector/nginx-77b746455b-: sidecars: kyverno-envoy-sidecar"

```