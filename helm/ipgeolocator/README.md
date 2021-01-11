# IPGeoLocator

##App settings

app.config. | Type | format | Description | Required or Default
------------ | ------------- | ------------- | ------------- | -------------
LISTEN_ADDRESS_PORT | string | "<ip address>:<pord>"| This the address and port that the IPGeoLocator app will bind to in the container. | Default: 0.0.0.0:8080
APP_URL_PATH | string | "/path" | This is the URL path that the app service will listen on for HTTP requests | Default: "/location"
METRICS_URL_PATH | string | "/metrics" | This is the URL path that the metrics service will listen on for HTTP requests | Default: "/metrics"
GL2_DATABASE_PATH | string | "/tmp" | This is the path to the directory that the GeoLite2 database file is stored in. | Required: No Default
GL2_DATABASE_FILENAME | string | "filename.mmdb" | This is the name of the GeoLite2 database file | Required: No Default
GL2_LOCALE | string | "en" | This is the locale used when reading the database | Default: "en"
METRICS | bool | true | Set as 'true' to enable the metrics server | Default: false
DEBUG | bool | true | set as 'true' to enable debug mode with logging | Default: false
  
  
secret.config. | Type | format | Description | Required or Default
------------ | ------------- | ------------- | ------------- | -------------
GL2_LICENSE_KEY | string | "du93ij39djd" | This is your license key to access the GeoLite2 databases for download. This is stored B64 encrypted in Kubernetes. This can also be set in the helm command | Required:  No Default


cronjob. | Type | format | Description 
------------ | ------------- | ------------- | -------------
schedule | string | ```'*/5 * * * *''``` | The cronjob schedule. When and how often the cronjob should trigger
backoffLimit | decimal | ```2``` | How many times the cronjob and fail and retry for failing permanently 
activeDeadlineSeconds | decimal | ```600``` | How long the cronjob can run for before triggering a failure


## Install command
```sh
helm install -n <namespace> <helm name> ipgeolocator/
or
helm install -n <namespace> --set app.secret.GL2_LICENSE_KEY="<your license key>" <helm name> ipgeolocator/
