#!/bin/bash

usage="Usage: ./save_db_data.sh [-p project_id] [path]"
project=dummy-project-id

while getopts ":p:" opt; do
  case $opt in
    p)
      project=$OPTARG
      ;;
    \?)
      echo $usage >&2
      exit 1
      ;;
    :)
      echo $usage >&2
      exit 1
      ;;
  esac
done

shift $((OPTIND - 1))

path=mount/db_data
if [ ! -z $1 ]
then
    path=$1
fi

container_id=`docker ps -aqf "name=uspy-backend-firestore_emulator"`
docker exec $container_id firebase emulators:export mount/db_data --force --project $project
