# start target container
weave run 10.0.0.5/8 busybox sleep 999999999

# run build once
weave build -f Dockerfile 10.0.0.11/8

# run build second time, logs will show cached
weave build -f Dockerfile 10.0.0.11/8

# run build with extra params
# notice from logs no cache is used
weave build -f Dockerfile 10.0.0.11/8 --no-cache -t sleepy
# notice image was made with tags
docker images | grep sleepy

