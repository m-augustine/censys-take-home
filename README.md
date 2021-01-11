<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#Testing">Testing</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

IPGeoLocator is a simple API that will a receive a IP Address and return seveal facts about the location of that IP Address around the globe.

IPGeoLocator uses the free GeoLite2 database. When run in a Kubernetes cluster, the database will be easy to keep up-to-date

### Built With

IPGeoLocator was built on the following technologies:
* [GO](https://golang.org/)
* [Docker](https://www.docker.com/)
* [Kubernetes](https://kubernetes.io/)
* [GeoLite2](https://www.maxmind.com/en/home)



<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

This is an example of how to list things you need to use the software and how to install them.
* GO >= 1.11
* In order to use the helm deployment, you will need a working kubernetes cluster
* A Maxmind license key - https://www.maxmind.com/en/home
* A tool for sending HTTP requests, like [Postman](https://www.postman.com/downloads/)

### Helm  Installation

1. Get your Maxmind License Key at [https://www.maxmind.com/en/home]https://www.maxmind.com/en/home)
2. Clone the repo
   ```sh
   git clone git@github.com:m-augustine/censys-take-home.git
   ```
3. (Optional) If you decide to build the container locally, run the following command, otherwise skip to step 4
   ```sh
   docker build -t mmaugust/censys-take-home:latest .
   ```
   If you build manually, you will need to make the images available to kubernetes cluster and alter the helm files to use the proper image repository
4. Move in the 'helm' directry
   ```JS
   cd helm
   ```liste
5. Run the helm installation command
   ```JS
   helm install -n <your namespace> <helm name> <chart>
   
   i.e.
   helm install -n app ipgeolocator ./ipgeolocator
   ```
liste
<!-- USAGE EXAMPLES -->metrics
## Testing

#### The request
1. Create the POST request in your HTTP(s) request tool
2. The request body is as follows. Replace the IPV4 with the IP of your choice
```sh
{
	"address": "68.48.244.19"
}
   ```
 #### Using Kubernetes port-forward
When the pod is running, you can port-forward to the application to test that it is working. 

#### Port-forward to pod
```sh
kubectl -n <pod namespace>  port-forward <pod name> <local port>:<app port>
   ```
   
Variable | Description
------------ | -------------
```<pod namepsace>``` | The Kubernetes namespace that you deployed the pod to
```<pod name>``` | The name of the Kubernetes pod.
```<local port>``` |  The port on your local machine that kubernetes will bind to and listen on in order to forward the requests to the <app port>. This can be any value within the normal range that isnt currently in use. 
```<app port>``` | The port that the application is listening on inside the pod. This is configured in the helm values file.
  

 Once the service port-forward is running, you can send your HTTP(s) POST request to ```http://127.0.0.1:<local port>/<app url path>```
 
  ```<app url path>``` is the url path that the server is listening on. This is configured in the helm values file
  i.e.   http://127.0.0.1:8080/location
  
  If you have enabled Metrics, you can send your HTTP(s) GET request to ```http://127.0.0.1:/<port>/<metrics url path>```
  
  ```<metrics url path>``` is the url path that the server is listening on. This is configured in the helm values file
  i.e. http://127.00.1:8080/metrics


#### Port-forward to service

Port forwarding to the Kubernetes service is very similar to port-forwarding to the pod. You simply target the service instead

Start by getting the service name and service port by using
```sh
kubectl -n <pod namespace> get service
   ```
   
 Then port-forward to the service using
 ```sh
kubectl -n <pod namespace> port-forward svc/<service name> <local port>:<service port>
   ```
     
Variable | Description
------------ | -------------
```<pod namepsace>``` | The Kubernetes namespace that you deployed the pod to
```<pod name>``` | The name of the Kubernetes pod.
```<local port>``` |  The port on your local machine that kubernetes will bind to and listen on in order to forward the requests to the <app port>. This can be any value within the normal range that isnt currently in use. 
```<service port>``` | The port that the service is listening on. This is configured in the helm values file.
  
 Once the service port-forward is running, you can send your HTTP(s) POST request to ```http://127.0.0.1:<local port>/<app url path>```
 
  ```<app url path>``` is the url path that the server is listening on. This is configured in the helm values file
  i.e.   http://127.0.0.1:8000/location
  
  If you have enabled Metrics, you can send your HTTP(s) GET request to ```http://127.0.0.1:/<port>/<metrics url path>```
  
  ```<metrics url path>``` is the url path that the server is listening on. This is configured in the helm values file
  i.e. http://127.00.1:80800/metrics

#### Using Kubernetes ingress

If you used a kubernetes ingress you will be able to target the Ingress IP address or configured URL  to make the HTTP(s) POST request



