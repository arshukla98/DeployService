init:
	mkdir -p ${HOME}/go/bin ${HOME}/go/pkg ${HOME}/go/src/github.com/DeployService && \
	cp -r . ${HOME}/go/src/github.com/DeployService && \
	echo "Run -> cd ${HOME}/go/src/github.com/DeployService"

env:
	export GOPATH=${HOME}/go && \
	export GO111MODULE=off

deps:
	go get -v ./...

create_crd:
	kubectl create -f artifacts/crd.yaml

build:
	go build -o main .

run:
	./main --kubeconfig=${HOME}/.kube/config

get_cr_yaml:
	cat artifacts/example.yaml

get_all:
	echo 'Getting DeployServices\n' && \
	kubectl get deployservices && \
	echo '\nGetting Deployments, Services\n' && \
	kubectl get deploy,svc

create_cr:
	kubectl create -f artifacts/example.yaml

desc:
	watch -n 0.1 'kubectl describe deployservices example-kube1'

cleanup:
	kubectl delete deployservices example-kube1
