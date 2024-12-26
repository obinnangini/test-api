# Test API Starter
This is a basic golang api configured to be built into a docker image and deployed via Skaffold & Helm to Kubernetes(K8s)
This application can be deployed locally to Minikube, or to a GKE cluster or any other K8s cluster.
Currently, it is configured to deploy to Minikube and GKE

Doc loosely used: https://crypto-gopher.medium.com/the-complete-guide-to-deploying-a-golang-application-to-kubernetes-ecd85a46c565

## Next Steps
1. Add Unit Testing for endpoints, service layer and DAO
1. Implement pagination using a limit based on timestamp
1. Add DB + migrate job to K8s deployment/GCP
1. Add DB indexing on searched fields like `slug`
1. GRPC API
1. Add Redis for caching results
1. On creation of an order, send an event via messaging systems to process order
1. Handle order processing message
1. Validate order name on add

## Setup (For Windows)
1. Install docker desktop: https://docs.docker.com/desktop/setup/install/windows-install/
1. Install Chocolatey: https://chocolatey.org/install
1. Install Minikube: https://minikube.sigs.k8s.io/docs/start/?arch=%2Fwindows%2Fx86-64%2Fstable%2Fchocolatey
    `choco install minikube`
1. Install kubectl: https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/#install-nonstandard-package-tools
   `choco install kubernetes-cli`
1. Install helm: https://helm.sh/docs/intro/install/
    Install and add to PATH
1. Install skaffold: https://skaffold.dev/docs/install/
   Install and add to PATH
1. Install make
    `choco install make`
1. Install gcloud: https://cloud.google.com/sdk/docs/install
1. Login to gcloud and created project
   `gcloud auth login`
1. Get valid artifact registry locations: https://cloud.google.com/artifact-registry/docs/repositories/repo-locations
1. create registry via gcloud
    `gcloud artifacts repositories create <REPO_NAME> --location <LOCATION> --repository-format docker --mode standard-repository --project suite <PROJECT_ID>`
1. authorize pull and push to registry for your user
   ```
    gcloud artifacts repositories add-iam-policy-binding <REPO_NAME> --location <LOCATION> --member=user:<USER_EMAIL> --role=roles/artifactregistry.writer --project <PROJECT_ID>
    gcloud artifacts repositories add-iam-policy-binding <REPO_NAME> --location <LOCATION> --member=user:<USER_EMAIL> --role=roles/artifactregistry.reader --project <PROJECT_ID>

   ```
1. Add credHelper to docker local for the upload
    The repo url format is document here - https://cloud.google.com/artifact-registry/docs/docker/names
    For example 
   `gcloud auth configure-docker us-central1-docker.pkg.dev`
1. skaffold set registry URL
    `skaffold config set default-repo us-central1-docker.pkg.dev/<PROJECT_ID>/<REPO_NAME>`

## Execution
### Local 
#### Docker Compose
1. Start Docker Desktop.
1. Run the make target to start the app via docker compsose
    ```
   make compose-up
   ```
1. In a different window curl `http://localhost:8080` to confirm the api is available
    ```
   > curl localhost:8080
    root.
   ```

#### Skaffold
1. Start Docker Desktop. 
1. Start minikube
    ```
    minikube start
   kubectl config use-context minikube
   ```
1. switch kubectl to minkube context
    ```
    kubectl config use-context minikube
   ```
1. run local deploy
    ```
   make deploy-local
   ```
   The deploy deploys the API to a local minikube cluster, and port forwards port `8080` from the deployment to `localhost:8080`
   After the deploy is complete, in a different window you can switch to the namespace and confirm a pod is running
    ```
   kubectl config set-context --current --namespace=test-api
   kubectl get pods
   ```

1. In a different window curl `http://localhost:8080` to confirm the api is available
    ```
    curl http://localhost:8080
    StatusCode        : 200
    StatusDescription : OK
    Content           : welcome
    RawContent        : HTTP/1.1 200 OK
    Content-Length: 7
    Content-Type: text/plain; charset=utf-8
    Date: Tue, 19 Nov 2024 00:32:43 GMT
    
                        welcome
    Forms             : {}
    Headers           : {[Content-Length, 7], [Content-Type, text/plain; charset=utf-8], [Date, Tue, 19 Nov 2024 00:32:43
    GMT]}
    Images            : {}
    InputFields       : {}
    Links             : {}
    ParsedHtml        : mshtml.HTMLDocumentClass
    RawContentLength  : 7
   ```
1. Hit Ctrl+C in the first window to stop port forwarding. To cleanup resources, run the clean command:
     ```
   make deploy-local-clean
   ```
### GKE
1. Configure kubectl with GKE cluster credential via gcloud
    ```
    gcloud container clusters get-credentials <CLUSTER_NAME> --location <LOCATION>
    ```
    For example:
    ```
    gcloud container clusters get-credentials sample-cluster --location us-central1-a
    ```
1. Get all contexts:
    ```
    kubectl config get-contexts
   ```
1. Switch to GKE context
    ```
   kubectl config use-context <CLUSTER_NAME>
   ```
1. run GKE deploy
     ```
    make deploy-gcp
    ```
    The deploy deploys the API to the GKE cluster, and creates a K8s service with an external IP address that exposes the API
    After the deploy is complete, you can switch to the namespace and confirm a pod is running
      ```
     kubectl config set-context --current --namespace=test-api
     kubectl get pods
     ```
1. Get the external IP assigned to the service
    ```
    kubectl get services
    NAME             TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
    test1-test-api   LoadBalancer   34.118.233.234   35.188.223.166   8080:31784/TCP   92s
    
    ```
1. curl `http://<EXTERNALIP>:8080` to confirm the api is available
    ```
    curl http://35.188.223.166:8080
    StatusCode        : 200
    StatusDescription : OK
    Content           : welcome
    RawContent        : HTTP/1.1 200 OK
    Content-Length: 7
    Content-Type: text/plain; charset=utf-8
    Date: Tue, 19 Nov 2024 00:59:34 GMT
    
                        welcome
    Forms             : {}
    Headers           : {[Content-Length, 7], [Content-Type, text/plain; charset=utf-8], [Date, Tue, 19 Nov 2024 00:59:34 GMT]}
    Images            : {}
    InputFields       : {}
    Links             : {}
    ParsedHtml        : mshtml.HTMLDocumentClass
    RawContentLength  : 7
   ```
1. To cleanup resources, run the clean command:
     ```
   make deploy-gcp-clean
   ```