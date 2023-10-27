# Kubernetes CRD and Custom Controllers

In Kubernetes, you can think of a "Custom Resource" as a way to create your own special objects to manage things that are unique to your application. These special objects are like the standard things Kubernetes knows how to handle, like Pods and Services, but they are tailored to your specific needs.

Imagine you're running a video game on Kubernetes, and you want to create a custom resource to represent a new type of in-game item. You can define the properties of this item, like its name, power, and special abilities, using something called a "Custom Resource Definition" (CRD).

A Custom Controller is a software component you create to watch and manage Custom Resources. It can automate tasks, like scaling applications or handling custom logic, based on the state of these custom objects.

# Current Environment

- Go (1.19)
- Kubectl (GitCommit : 1b4df30b3, Git Version: v1.27.0)
- KubeBuilder (3.5.0)
- Packages need to Install using apt such as make, build-essential
- Linux/AMD64

# Establish Current Environment
- Open Terminal.
- Install Packages such as make, build-essential and go 1.19.
```
controlplane ~ ➜ sudo apt update && sudo apt install make build-essential 
controlplane ~ ➜  export VERSION=1.19
controlplane ~ ➜  curl  -L https://golang.org/dl/go${VERSION}.linux-amd64.tar.gz -o go${VERSION} && tar -xzf go${VERSION} -C /usr/local
controlplane ~ ➜  export PATH=$PATH:/usr/local/go/bin
controlplane ~ ➜  go version
go version go1.19 linux/amd64
```
- Now Configure Your Git.
```
controlplane ~ ➜  git --version
git version 2.25.1
controlplane ~ ➜  git config --global user.email EMAILID
controlplane ~ ➜  git config --global user.name USERNAME
controlplane ~ ➜  git config --global credential.helper   cache
controlplane ~ ➜  
```
- Install KubeBuilder 3.5.0
```
controlplane ~ ➜  cat > kube_builder.sh
curl -L https://github.com/kubernetes-sigs/kubebuilder/releases/download/v3.5.0/kubebuilder_linux_amd64 -o kubebuilder_3.5.0_linux_amd64
chmod +x  kubebuilder_3.5.0_linux_amd64
sudo mv  kubebuilder_3.5.0_linux_amd64 /usr/local/bin/kubebuilder
^C
controlplane ~ ➜  bash kube_builder.sh
```

# Create an Empty Go Module
- Create CRD Directory
```
controlplane ~ ➜  mkdir Persona
controlplane ~ ➜  cd Persona
```
- Specify the location of all the project dependencies in GOMODCACHE.
```
controlplane ~/Persona ➜  go env -w GOMODCACHE=$PWD/.deps/pkg/mod
controlplane ~/Persona ➜  go env | grep GOMODCACHE
GOMODCACHE="/root/Persona/.deps/pkg/mod"
```
- Initialize this whole directory as Go Module and check the contents.
```
controlplane ~/Persona ➜  go mod init github.com/USERNAME/PersonaCRD
go: creating new go.mod: module github.com/USERNAME/PersonaCRD
controlplane ~/Persona via 🐹 v1.19 ✦ ➜  ls -la
total 20
drwxr-xr-x 3 root root 4096 Oct 27 07:26 .
drwx------ 1 root root 4096 Oct 27 07:25 ..
drwxr-xr-x 3 root root 4096 Oct 27 07:26 .deps
-rw-r--r-- 1 root root   53 Oct 27 07:26 go.mod
```
# Push This change to the Github Repo.
- Create Github Repo with name 'PersonaCRD' [https://github.com/new] with desired visibility.
- Go Back to Terminal 
- Initialized Persona as Git Repo
```
controlplane ~ ➜  cd Persona
controlplane ~/Persona via 🐹 v1.19 ➜  git init
Initialized empty Git repository in /root/Persona/.git/
```
- Add README.md file and Save it.
```
controlplane Persona on  master [?] via 🐹 v1.19 ➜  vi README.md
```
- View the Content  
```
controlplane Persona on  master [?] via 🐹 v1.19 ➜  git status
```
- Mention '.deps' in .gitignore. Whatever you mention in .gitignore, it will not be pushed to repository.
```
controlplane Persona on  master [?] via 🐹 v1.19 ➜  vi .gitignore
```
- Now if you try the below command, you will not see .deps and there is no need to worry becuase .deps can be re-generated via 'go mod download'.
```
controlplane Persona on  master [?] via 🐹 v1.19 ➜  git status
```
- Push the Content to Staging Area.
```
controlplane Persona on  master [?] via 🐹 v1.19 ➜  git add .
```
- View the Staged Changes.
```
controlplane Persona on  master [?] via 🐹 v1.19 ➜  git status
```
- Now Commit so that later it can be pushed to Github. The Output will be similar as shown below.
```
controlplane Persona on  master [+] via 🐹 v1.19 ➜  git commit -m "first commit"
[master (root-commit) 155261f] first commit
 2 files changed, 4 insertions(+)
 create mode 100644 .gitignore
 create mode 100644 go.mod
```
- Define Remote Repository URL. 
```
controlplane Persona on  master via 🐹 v1.19 ➜ git remote add origin https://github.com/Aditya98Shukla/PersonaCRD.git
controlplane Persona on  master via 🐹 v1.19 ➜ git branch -M main
```
- Push the Main branch to Remote Repo. Make Sure you have mention Gihub USERNAME and Personal Access Token as Password.
```
controlplane Persona on  master via 🐹 v1.19 ➜  git push -u origin main
```
## Initializing Kubernetes Project
- Clone the Repository
```
controlplane ~ ✦ ➜  git clone https://github.com/Aditya98Shukla/PersonaCRD.git
Cloning into 'PersonaCRD'...
remote: Enumerating objects: 5, done.
remote: Counting objects: 100% (5/5), done.
remote: Compressing objects: 100% (4/4), done.
remote: Total 5 (delta 0), reused 5 (delta 0), pack-reused 0
Unpacking objects: 100% (5/5), 2.20 KiB | 2.20 MiB/s, done.
controlplane ~ ✦ ➜  cd PersonaCRD
```
- Check If there are any changes. I hope not.
```
controlplane PersonaCRD on  main via 🐹 v1.19 ✦ ➜  git status
On branch main
Your branch is up to date with 'origin/main'.
nothing to commit, working tree clean
```
-  Initializes a Kubernetes project using the Kubebuilder framework and sets the domain for your project to "genesis.xyz.com." This domain typically represents the root domain for your Kubernetes Custom Resource Definitions (CRDs) and can help identify your custom resources in a cluster.
```
controlplane PersonaCRD on  main via 🐹 v1.19 ➜  kubebuilder init --domain genesis.xyz.com
Writing kustomize manifests for you to edit...
Writing scaffold for you to edit...
Get controller runtime:
$ go get sigs.k8s.io/controller-runtime@v0.12.1
......
go: downloading github.com/nxadm/tail v1.4.8
go: downloading gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
Next: define a resource with:
$ kubebuilder create api
```
- Note the KubeBuilder Version.
```
controlplane PersonaCRD on  main [!?] via 🐹 v1.19 ✦ ➜  kubebuilder version
Version: main.version{KubeBuilderVersion:"3.5.0", KubernetesVendor:"1.24.1", GitCommit:"26d12ab1134964dbbc3f68877ebe9cf6314e926a", BuildDate:"2022-06-24T12:17:52Z", GoOs:"linux", GoArch:"amd64"}
```
- You will notice the changes after applying 'git status'. Move these changes to staging area using 'git add .' and then commit. Later Push the committed changes.
