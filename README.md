# Weave hack to add networking to docker builds
*note: I just hacked put it in weave bash file instead of writing a wrapper*

How it works
------------
Step 1. modify Dockerfile to add a line to wait for weave before running
`until grep -q ethwe /proc/net/dev; do sleep .5; done;`
that line waits for weave interface to come up before running the command

Step 2. Run docker build and watch the output steam. when it is a non cached container
run weave attach with the CIDR ip provided

Use Case
--------
This can be helpful when you have some resource server that you host inside if your weave network

Demo
----
run `./demo.sh` for example

it first creates a target container. then runs docker build with a Dockerfile which just pings that container
it also shows normal build args like `-t` tag and `--no-cache` get passed to the build line

Limitations
-----------
right now you have to have `-f Dockerfile/path` as first argument