# Patakube for admins

## Requirements

* Kubernetes cluster (even if only 1 node)
* Room with wifi able to reach cluster

You can use a [GKE cluster](https://cloud.google.com/container-engine/) to ease you.

## Before event

- [ ] Ask people to come with a working computer
- [ ] Create number of namespaces  
There will be 1 per user who want to play.  
`kubectl create namespace player-{NAME I WANT}`
* [ ] spawn _target_, _trap_ and _arbiter_ agent into `patakube` namespace  
```
kubectl create namespace patakube
kubectl create --namespace=patakube --filename=config/admin.yaml
```

## During event

TODO
