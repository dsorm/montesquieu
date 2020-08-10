#!/bin/bash
# Used to generate config.json from environment variables passed to docker container

echo "{ \"BlogName\":\"${BLOGNAME}\",\"ArticlesPerPage\":\"${ARTICLESPERPAGE}\",	\"ListenOn\":\"${LISTENON}\",\"ArticleStore\":\"${ARTICLESTORE}\",\"ArticleStoreHost\":\"${ARTICLESTORE_HOST}\",\"ArticleStoreDB\":\"${ARTICLESTORE_DB}\",\"ArticleStoreUser\":\"${ARTICLESTORE_USER}\",\"ArticleStorePassword\":\"${ARTICLESTORE_PASSWORD}\",\"CachingEngine\":\"${CACHINGENGINE}\"}" > config.json