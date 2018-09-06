
.PHONY: minikube-setup
minikube-setup: minikube-start set-namespace ##@minikube start minikube
	sudo minikube addons enable freshpod
	sudo minikube addons enable ingress
	sudo minikube addons enable heapster

.PHONY: minikube-start
minikube-start: ##@minikube start minikube
	sudo CHANGE_MINIKUBE_NONE_USER=true minikube start --vm-driver=none

.PHONY: set-namespace
set-namespace: ##@minikube set the namespace
	kubectl apply -f environment/namespace.yaml
	kubectl config set-context minikube --namespace=eventstore-example

.PHONY: minikube-stop
minikube-stop: ##@minikube stop minikube
	sudo minikube stop

.PHONY: minikube-delete
minikube-delete: ##@minikube remove minikube
	sudo minikube delete
