#!/bin/bash

spruce --concourse merge --prune meta global.yml docker.yml deploy.yml > pipeline.yml
fly -t tutorial c play.spruce.cf -c pipeline.yml --vars-from credentials.yml
