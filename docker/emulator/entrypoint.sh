#!/bin/bash

if [ -z $IMPORT_DATA ]
then
    export IMPORT_DATA=mount/db_data
fi

echo "{
  \"emulators\": {
    \"firestore\": {
      \"port\": \"$PORT\",
      \"host\": \"0.0.0.0\"
    },
    \"ui\": {
      \"enabled\": true,
      \"host\": \"0.0.0.0\",
      \"port\": \"$PORT_UI\"
    }
  }
}" > firebase.json

if [ -d "$IMPORT_DATA" ]
then
    echo "Importing emulator data from '$IMPORT_DATA'"
    firebase emulators:start --project $FIRESTORE_PROJECT_ID --import $IMPORT_DATA
else
    echo "Starting a fresh database"
    firebase emulators:start --project $FIRESTORE_PROJECT_ID
fi
