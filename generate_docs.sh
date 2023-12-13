#!/bin/bash

swag init -g netflix/netflix.go

cd docs

bootprint openapi swagger.json target
