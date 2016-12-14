#!/bin/bash

set -e
set -u

# = Vars

CLUSTER_NAME=k8s_patakube
CLUSTER_URL={{.ClusterUrl}}

CONTEXT_NAME=patakube
CONTEXT_NAMESPACE={{.Namespace}}

K8S_RELEASE_URL=https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin


# = Setup
#
# == Kubectl
echo '---> Downloading kubectl'
case "$(uname -s)" in

   Darwin)
     echo 'OS: Mac OS X'
     curl -LO ${K8S_RELEASE_URL}/darwin/amd64/kubectl
     ;;

   Linux)
     echo 'OS: Linux'
     curl -LO ${K8S_RELEASE_URL}/linux/amd64/kubectl
     ;;

   *)
     echo 'Unknown OS'
     exit 2;
     ;;
esac

chmod +x ./kubectl


# == Configure
echo '---> Configuring kubectl'

# create kubeconfig entry
./kubectl config set-cluster ${CLUSTER_NAME} \
    --server=${CLUSTER_URL}

# create context entry
./kubectl config set-context ${CONTEXT_NAME} \
    --cluster=${CLUSTER_NAME} \
    --namespace=${CONTEXT_NAMESPACE}

# use the context
./kubectl config use-context ${CONTEXT_NAME}


# = Run
echo '---> Testing kubectl config:'

echo '---> kubectl config view'
./kubectl config view

echo '---> kubectl cluster-info'
./kubectl cluster-info

echo '---> Done'
