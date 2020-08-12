#!/bin/bash
# Used to generate config.json from environment variables passed to docker container

echo "{ \"BlogName\":\"${BLOGNAME}\",\"ArticlesPerPage\":\"${ARTICLESPERPAGE}\",	\"ListenOn\":\"${LISTENON}\",\"Store\":\"${STORE}\",\"StoreHost\":\"${STORE_HOST}\",\"StoreDB\":\"${STORE_DB}\",\"StoreUser\":\"${STORE_USER}\",\"StorePassword\":\"${STORE_PASSWORD}\",\"CachingStore\":\"${CACHINGSTORE}\",\"HotSwapTemplates\": \"${HOTSWAPTEMPLATES}\"}" > config.json