# Installing jupyter  notebook

Recommand to install on vm prevent mixing with any pre-owned docker containers

## jupyter docker
- docker run -p 8888:8888 jupyter/scipy-notebook

## native
- sudo apt install python3-pip
- pip3 install notebook
- jupyter notebook ( --ip=0.0.0.0 )
