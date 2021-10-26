DEPLOYMENT_DIR := .deployment
DEPLOYMENT_PATH := ${DEPLOYMENT_DIR}/docker-compose.yaml

# run do_deploy build locally to build and compress the docker images
do_deploy_build:
# clear deployment dir
	rm -r ${DEPLOYMENT_DIR}
	mkdir ${DEPLOYMENT_DIR}
# set up directories and copy files
	cp docker-compose.deploy.yaml ${DEPLOYMENT_DIR}/docker-compose.yaml
	cp ./migrations/postgres/20211019000000-init.sql ${DEPLOYMENT_DIR}/create_tables.sql
	cp makefile ${DEPLOYMENT_DIR}/makefile
	cp .env.deploy ${DEPLOYMENT_DIR}/.env
# build docker images
	docker-compose -f ${DEPLOYMENT_PATH} build  --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} bb3-ocpp-ws
	docker-compose -f ${DEPLOYMENT_PATH} build  pgweb
# save docker images
	docker save -o ${DEPLOYMENT_DIR}/image_bb3-ocpp-ws.tar.gz deployment_bb3-ocpp-ws
	docker save -o ${DEPLOYMENT_DIR}/image_pgweb.tar.gz deployment_pgweb

# run do_deploy_load to load the docker images after they have been transferred over
do_deploy_load:
# load docker images
	docker load -i image_bb3-ocpp-ws.tar.gz
	docker load -i image_pgweb.tar.gz

.PHONY: do_deploy_build
