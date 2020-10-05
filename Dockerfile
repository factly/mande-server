FROM golang:1.14.2-alpine3.11

WORKDIR /app

COPY . .

RUN go mod download

ENV DSN $DSN
ENV MEILI_URL $MEILI_URL
ENV MEILI_KEY $MEILI_KEY
ENV RAZORPAY_KEY $RAZORPAY_KEY
ENV RAZORPAY_SECRET $RAZORPAY_SECRET

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -exclude-dir=.git -exclude-dir=docs --build="go build main.go" --command="./main -dsn=${DSN} -meili=${MEILI_URL} -meiliKey ${MEILI_KEY} -razorpayKey=${RAZORPAY_KEY} -razorpaySecret=${RAZORPAY_SECRET}"