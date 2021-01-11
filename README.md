<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
        <li><a href="#application-concepts">Application Concepts</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#testing">Testing</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

IPGeoLocator is a simple API that will a receive a IP Address and return seveal facts about the location of that IP Address around the globe.

IPGeoLocator uses the free GeoLite2 database provided by MaxMind. When run in a Kubernetes cluster, the database will be easy to keep up-to-date.

### Built With

IPGeoLocator was built on the following technologies:
* [GO](https://golang.org/)
* [Docker](https://www.docker.com/)
* [Kubernetes](https://kubernetes.io/)
* [GeoLite2](https://www.maxmind.com/en/home)

<!-- APPLICATION CONCEPTS -->
### Application Concepts
The IPGeoLocator application is designed to be containerized and  deployed via Helm into a Kubernetes environment. The Helm chart contains all the approriate Kubernetes resrouce definitions to make the application and 1-line install command

Once started, Init containers will download the MaxMind GeoLite2-City database into the container filesystem. The IPGeoLocator app will use this database to return Geo location information pertaning to any IP Address that it recieves via its REST API.

Kubernetes cronjobs are deployed with a default once-per-week trigger to perform a rolling restart on the deployment. During this process, each pod is restarted which. When a pod restarts, it downloads the latest version of the MaxMind GeoLite2 database. Pods are restarted one-by-one. The next pod will not be restarted untill the previous pod is in a ready state. If you have only 1 replica, the service will be unavailable during this time. 

The prjoect utilizes Github Actions to build a new version of the container and push it to the public Docker repository whenever there is a push to the Main branch or when a PR is merged in.

<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

* GO >= 1.11 ( for running locally)
* In order to use the helm deployment, you will need a working kubernetes cluster
* A Maxmind license key - https://www.maxmind.com/en/home
* A tool for sending HTTP requests, like [Postman](https://www.postman.com/downloads/) or Curl

### (optional) Building the Container Locally
1. Clone the repo.
   ```sh
   git clone git@github.com:m-augustine/ipgeolocator.git
   ```
2. Buld the container.
   ```sh
   cd ipgeolocator
   docker build -t mmaugust/ipgeolocator:latest .
   ```
   If you build manually, you will need to make the images available to kubernetes cluster and alter the helm files to use the proper image repository.

## Installation

#### Helm
1. Get your MaxMind License Key at [the MaxMind website]https://www.maxmind.com/en/home)
2. Clone the repo.
   ```sh
   git clone git@github.com:m-augustine/ipgeolocator.git
   ```
3. Move in the 'helm' directry.
   ```JS
   cd helm
   ```liste
4. Update the settings in values.yaml to configure the application for your needs. You can set the GL2_LICENSE_KEY in the values.yaml file or in the helm command itself.
    Helm settings decriptions can be found in the Helm [README.md](https://github.com/m-augustine/ipgeolocator/blob/main/helm/ipgeolocator/README.md)
	
5. Run the helm installation command.
   ```JS
   helm install -n <your namespace> <helm name> <chart>
   or
   helm install -n <your namespace> --set app.secret.GL2_LICENSE_KEY="<your license key>" <name> <chart>
   
   i.e.
   helm install -n default --set app.secret.GL2_LICENSE_KEY="d99wjnkieje" ipgeolocator ./ipgeolocator
   ```

<!-- TESTING -->
## Testing

#### The HTTP Request
1. Create the POST request in your HTTP(s) request tool
2. The request body is as follows. Replace the IPV4 with the IP of your choice
```sh
{
	"address": "68.48.244.19"
}
   ```
 
#### Using the NodePort with Docker For Desktop

If you left the default for the Kubernetes service, then the service is listening on a nodeport of the worker machine that it is on.
1. The NodePort should be listening on your machines localhost (127.0.0.1). If it fails with this address, you will need to find the IP address of the node that the pod is running on.
2. Get the service NodePort.
```sh
$$ kubectl -n <pod namespace> get service
NAME           TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
ipgeolocator   NodePort   10.232.98.124   <none>        80:30452/TCP   4s

```
In the 'PORT(S)' column you will somthing similar to ```80:30452/TCP```. The nodeport in this example would be 30452.

3. Submit your HTTP request.
```sh
i.e.
curl -d '{"address":"68.48.244.192"}' -H 'Content-Type: application/json' http://127.0.0.1:30187/location

If successful, the response will look like:
{"Continent":"North America","Country":"United States","City":"Ann Arbor","Location":{"Timezone":"America/Detroit","Latitude":42.2324,"Longitude":-83.7017}}
```
 
   
#### Using Kubernetes Port-Forward
When the pod is running, you can port-forward to the application via the Kubernetes service or directly to the pod itself to test that it is working. 

#### Port-Forward to Service

Port forwarding to the Kubernetes service is very similar to port-forwarding to the pod. You simply target the service instead.

Start by getting the service name and service port by using.
```sh
kubectl -n <pod namespace> get service
   ```
   
 Then port-forward to the service using.
 ```sh
kubectl -n <pod namespace> port-forward svc/<service name> <local port>:<service port>
   ```
     
Variable | Description
------------ | -------------
```<pod namepsace>``` | The Kubernetes namespace that you deployed the pod to.
```<pod name>``` | The name of the Kubernetes pod.
```<local port>``` |  The port on your local machine that kubernetes will bind to and listen on in order to forward the requests to the <app port>. This can be any value within the normal range that isnt currently in use. 
```<service port>``` | The port that the service is listening on. This is configured in the helm values file.
  
 Once the service port-forward is running, you can send your HTTP(s) POST request to ```http://127.0.0.1:<local port>/<app url path>```
 
  ```<app url path>``` is the url path that the server is listening on. This is configured in the helm values file.
  i.e.   http://127.0.0.1:8000/location
  
  If you have enabled Metrics, you can send your HTTP(s) GET request to ```http://127.0.0.1:/<port>/<metrics url path>```
  
  ```<metrics url path>``` is the url path that the server is listening on. This is configured in the helm values file.
  i.e. http://127.00.1:80800/metrics
  
  

#### Port-Forward to Pod
```sh
kubectl -n <pod namespace>  port-forward <pod name> <local port>:<app port>
   ```
   
Variable | Description
------------ | -------------
```<pod namepsace>``` | The Kubernetes namespace that you deployed the pod to.
```<pod name>``` | The name of the Kubernetes pod.
```<local port>``` |  The port on your local machine that kubernetes will bind to and listen on in order to forward the requests to the <app port>. This can be any value within the normal range that isnt currently in use. 
```<app port>``` | The port that the application is listening on inside the pod. This is configured in the helm values file.
  

 Once the service port-forward is running, you can send your HTTP(s) POST request to ```http://127.0.0.1:<local port>/<app url path>```
 
  ```<app url path>``` is the url path that the server is listening on. This is configured in the helm values file.
  i.e.   http://127.0.0.1:8080/location
  
  If you have enabled Metrics, you can send your HTTP(s) GET request to ```http://127.0.0.1:/<port>/<metrics url path>```
  
  ```<metrics url path>``` is the url path that the server is listening on. This is configured in the helm values file.
  i.e. http://127.00.1:8080/metrics

#### Using Kubernetes Ingress

If you used a kubernetes ingress you will be able to target the Ingress IP address or configured URL  to make the HTTP(s) POST request.

#### Direct Access

If you have a machine on the Kubernetes service network, or are using another pod in the same cluster you can access the IPGeiLocator Kubernetes service directly.

1. Get the IPGeiLocator Kubernetes service address.
```sh
kubectl -n <pod namespace> get service
```
2. Note the value in the 'CLUSTER_IP' column. This will be the address that you use to reach the IPGeoLocator app.

3. Send your HTTP(s) request.


## Troubleshooting

#### The Pod is Constantly Restarting
You may need to increase the delay on the readiness and liveness probes. If the pod is taking longer that normal to download the database and start the app, the probe may be trying to reach the TCP socket before its ready and failing. After failing a set amount of times, the pod will be restarted

#### The Init Container is Failing
1. You may have entered a wrong or malformed Maxmind Licence key. This woudl cause the download to fail and subsequently, the init container.
2. Check that your pods are able to reach the internet
3. The download is failing because you have reached your download limit for your Maxmind account


#### The Application Crasher Randomly
Please submit and issues discoverd in [Issues Page](https://github.com/m-augustine/ipgeolocator/issues). 
Please provide as much detail as you can, including docker logs from the application and init containers.







