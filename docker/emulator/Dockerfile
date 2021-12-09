FROM node:lts-alpine

RUN apk --no-cache add openjdk11-jre bash     && \
    yarn global add firebase-tools@${VERSION} && \
    yarn cache clean                          && \
    firebase setup:emulators:firestore        && \
    firebase setup:emulators:ui               && \
    chown -R node:node $HOME

WORKDIR /home/node

ENV FIRESTORE_PROJECT_ID dummy-project-id
ENV PORT 8080
ENV PORT_UI 4000

EXPOSE "$PORT"
EXPOSE "$PORT_UI"

COPY entrypoint.sh .

ENTRYPOINT ["./entrypoint.sh"]
