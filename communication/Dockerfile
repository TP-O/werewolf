FROM node:20-alpine3.16 as build-stage

WORKDIR /app

ARG PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true

COPY . .

RUN npx pnpm install && npx prisma generate && npx pnpm build

FROM node:20-alpine3.16

LABEL maintainer="Le Tran Phong <letranphong2k1@gmail.com>"

WORKDIR /app

ENV NODE_ENV production

RUN mkdir dist
RUN chown -R node:node /app

COPY --chown=node:node ./package.json ./pnpm-lock.yaml ./
COPY --chown=node:node ./prisma/schema.prisma ./prisma/schema.prisma
COPY --chown=node:node --from=build-stage /app/dist /app/dist

USER node
RUN npx pnpm install -P && npx prisma generate

CMD [ "node", "/app/dist/main.js" ]
