FROM node:18.15-alpine3.16 as build-stage

WORKDIR /app

COPY . .

ARG PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true

RUN npx pnpm install \
  && npx prisma generate \
  && npx pnpm build

FROM node:18.15-alpine3.16

LABEL maintainer="Le Tran Phong <letranphong2k1@gmail.com>"

WORKDIR /app

RUN chown -R node:node /app

USER node

RUN mkdir dist

COPY --chown=node:node ./package.json ./pnpm-lock.yaml ./
COPY --chown=node:node ./prisma/schema.prisma ./prisma/schema.prisma
COPY --chown=node:node --from=build-stage /app/dist /app/dist

RUN npx pnpm install -P && npx prisma generate
# Remove cache
RUN rm -rf "$(npx pnpm store path) ~/.npm/"

CMD [ "node", "/app/dist/main.js" ]
