# Knockout

Knockout is a kubectl plugin to terminate resources that are stuck in the terminating phase.

## Installation

```bash
git clone https://github.com/bensch98/kube-knockout.git
cd kube-knockout
bash ./build.sh
```

## Usage

```bash
kubectl knockout -n <NAMESPACE>
```
