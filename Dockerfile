FROM alpine:3.20 as build

RUN apk add --no-cache nodejs npm go make bash

WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download

COPY package.json package-lock.json /src/
RUN npm ci

COPY . /src

RUN make build

FROM alpine:3.20 as runtime

COPY --from=build /src/bin/sidequest /app/bin/sidequest
COPY --from=build /src/dist /app/dist

VOLUME /app/data

EXPOSE 8090

WORKDIR /app

CMD ["/app/bin/sidequest", "serve", "--http", "0.0.0.0:8090"]