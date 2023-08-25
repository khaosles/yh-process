### sea-forecast

docker run -d \
--rm \
--name rmqbroker \
--privileged=true \
apache/rocketmq \
sh mqbroker
docker cp rmqbroker:/home/rocketmq/rocketmq-5.1.3/bin/runbroker.sh /broker/bin/runbroker.sh
