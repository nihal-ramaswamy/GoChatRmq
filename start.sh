~/.docker/cli-plugins/docker-compose up -d go_db
~/.docker/cli-plugins/docker-compose up -d cache_db
~/.docker/cli-plugins/docker-compose up -d amqp
~/.docker/cli-plugins/docker-compose build
~/.docker/cli-plugins/docker-compose up go_chat
