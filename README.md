## kubectl-cns

kubectl-cns is a kubectl plugin that clean the namespace quickly. 

This plugin has been tested to work with following delete types:

- Delete Active namespace
- Delete multi Active namespace
- Delete Terminating namespace
- Delete multi Terminating namespace

## Build

if you have go env. build yourself. or [download](https://github.com/cuisongliu/kubectl-cns/releases)


```
git clone git@github.com:cuisongliu/kubectl-cns.git

go build .

cp kubectl-cns $YOUR_PATH

```

## Usage

Binary Use

```
 
# delete an tx namespace
kubectl-cns tx

# delete an tx namespace by force
kubectl-cns tx --force

# delete some namespaces like tx, staging
kubectl-cns tx staging

# delete some namespaces like tx, staging , qa by force
kubect-cns tx staging qa --force
	

```

As Plugin Use

```
# delete an tx namespace
kubectl cns tx

# delete an tx namespace by force
kubectl cns tx --force

# delete some namespaces like tx, staging
kubectl cns tx staging

# delete some namespaces like tx, staging , qa by force
kubectl cns tx staging qa --force

```
