#!/bin/sh
docker build --push -t $1 --annotation "composehackv1=true" .