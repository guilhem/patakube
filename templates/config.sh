#!/bin/bash

# = Kubectl
echo '---> Downloading kubctl'
case "$(uname -s)" in

   Darwin)
     echo 'OS: Mac OS X'
     curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/darwin/amd64/kubectl
     ;;

   Linux)
     echo 'OS: Linux'
     curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
     ;;

   *)
     echo 'Unknown OS'
     exit 2;
     ;;
esac

chmod +vx ./kubectl


# = Configure
export CLUSTER_NAME=k8s_patakube
export CLUSTER_URL=http://be59165d.eu.ngrok.io/
export KUBECONFIG=./kube_config
export CONTEXT_NAME=patakube

# create kubeconfig entry
./kubectl config set-cluster ${CLUSTER_NAME} \
    --server=${CLUSTER_URL}

# create context entry
./kubectl config set-context ${CONTEXT_NAME} \
    --cluster=${CLUSTER_NAME} \
    --namespace={{.Namespace}}

# use the context
./kubectl config use-context ${CONTEXT_NAME}

./kubectl config view
