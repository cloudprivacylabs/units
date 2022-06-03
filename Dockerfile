FROM golang:1.18-alpine AS build

WORKDIR /src/cmd
COPY *.go go.* /src/
COPY cmd/ /src/cmd/
RUN CGO_ENABLED=0 go build -o /bin/cmd

FROM node:17-alpine

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+)
COPY ucum/package*.json ./
COPY ucum/app.js ./
COPY runserver.sh ./

RUN npm install
# If you are building your code for production
# RUN npm ci --only=production

EXPOSE 8080

COPY --from=build /bin/cmd .

CMD [ "/bin/sh", "runserver.sh" ]