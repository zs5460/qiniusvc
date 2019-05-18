#build
FROM golang AS build

WORKDIR /go/src/github.com/zs5460/qiniusvc

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build .

CMD ["./qiniusvc"]

#production
FROM scratch AS prod

COPY --from=build /go/src/github.com/zs5460/qiniusvc/public ./public

COPY --from=build /go/src/github.com/zs5460/qiniusvc/qiniusvc .

CMD ["./qiniusvc"]