DEPLOYMENT_DIR := .deployment
DEPLOYMENT_PATH := ${DEPLOYMENT_DIR}/docker-compose.yaml

# run do_deploy build locally to build and compress the docker images
do_deploy_build:
# set up directories and copy files
	mkdir -p ${DEPLOYMENT_DIR}/pg_data
	cp docker-compose.deploy.yaml ${DEPLOYMENT_DIR}/docker-compose.yaml
	cp create_tables.sql ${DEPLOYMENT_DIR}/create_tables.sql
	cp makefile ${DEPLOYMENT_DIR}/makefile
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
