kubectl run player --image=patakube/player --expose --port=8080 --env=[]
kubectl scale --replicas=3 deployment/player

kubectl apply -f deployment.yaml
