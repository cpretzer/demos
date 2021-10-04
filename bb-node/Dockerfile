FROM node:14 AS build-env
COPY . /app
WORKDIR /app

RUN yarn install --froze-lockfile --prod=true

FROM gcr.io/distroless/nodejs:14
COPY --from=build-env /app /app
WORKDIR /app
CMD ["index.js"]