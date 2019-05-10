# grpctest

Example code for Istio on GKE with grpc services with grpc-json transcoding, external authorization and rate limiting.


#Create persistent disk

A gke persistent disk is needed for the grpc descriptor files used for transcoding.
 
 ## Create Disk
 
  ```
  export PB_DISK=<disk_name>
  export GCP_PROJECT=<project_id>
  export GCP_ZONE=<zone>
  export INSTANCE_NAME=<instance:
  gcloud compute disks create ${PB_DISK} --project=${GCP_PROJECT} --type=pd-standard --size=10GB --zone=${GCP_ZONE}    
  ```
  ## Create compute instance and attach disk
  ```
  gcloud compute --project=zing-dev-197522 instances create ${INSTANCE_NAME} --zone=${GCP_ZONE} --machine-type=f1-micro  --image-family=ubuntu-1804-lts --image-project=ubuntu-os-cloud --boot-disk-size=10GB --disk=name=${PB_DISK},device-name=${PB_DISK},mode=rw,boot=no  
 ```
 ## Initialize and mount disk
    Commands are run individually instead of in a script for transparency
 ```
 gcloud compute ssh ${INSTANCE_NAME} --command "sudo mkfs.ext4 -m 0 -F -E lazy_itable_init=0,lazy_journal_init=0,discard /dev/sdb &&  sudo mkdir -p /proto_mount && sudo mount -o discard,defaults /dev/sdb /proto_mount && sudo chmod a+w /proto_mount"
 ```
 ## Copy proto files to disk
 ```
  gcloud compute scp pb/*.pb ${INSTANCE_NAME}:/proto_mount/.
 ```
 ## Delete compute instance
 Need to delete so K8s can mount disk in read only mode
 ```
gcloud compute instances delete ${INSTANCE_NAME}

```

## modify istio sidecar injector yaml to have persistent disks

Download the current istio-sidecar-injector configmap
 
```
kubectl -n istio-system get configmap istio-sidecar-injector -o=jsonpath='{.data.config}' > inject-config.yaml
```
modify with gce volumes, find the correct places to add:
 ```
               - mountPath: /gce-disk2
                 name: gce-disk
                 readOnly: true
             volumes:
             - name: gce-disk
               gcePersistentDisk:
                 pdName: jpl-istio-proto
                 readOnly: true
                 fsType: ext4
``` 
   
create modified conffig map

```kubectl -n istio-system create configmap istio-jpl --from-file=config=inject-config.yaml```

   
Inject containers with modified sidecar yaml

```kubectl apply -f <(istioctl kube-inject --injectConfigMapName istio-jpl -f grpctest.yaml)```
 
 
#Build proto descriptorss ch 
 Create proto descriptor
 protoc --include_imports --include_source_info --proto_path=${GOOGLEAPIS_DIR} --proto_path=. --descriptor_set_out=grpc_test_1.pb grpc_test.proto
 
 
 using grpc transcoder
 https://www.envoyproxy.io/docs/envoy/latest/configuration/http_filters/grpc_json_transcoder_filter.html#how-to-generate-proto-descriptor-set
 
 
 
 
 curl -d '{"value":2}' -X POST -H "Content-Type: application/json" -kv https://jpl.zenoss.io/IanTestService/Square
 curl -d '{}' -X POST -H "Content-Type: application/json" -kv https://jpl.zenoss.io/IanTestService/Random
 curl -d '{}' -X POST -H "Content-Type: application/json" -kv https://jpl.zenoss.io/IanTestService/Randomfdjlksdf
 
 curl -H "Content-Type: application/json" -kv https://jpl.zenoss.io/math/random
 
 
 
 
 
 
 gcloud compute ssh  jpl-test  --command "sudo mkdir -p /demo-mount"
 gcloud compute ssh  jpl-test  --command "sudo mount -o discard,defaults /dev/sdb /demo-mount"
 gcloud compute scp *.pb jpl-test:/demo-mount/.
 
 
 gcloud compute instances delete jpl-test
