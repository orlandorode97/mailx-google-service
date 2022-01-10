FROM golang:1.17

WORKDIR /app
COPY . .

RUN go build ./cmd/mailx-google-service

RUN ls
COPY ./dev.sh /
RUN chmod +x /dev.sh

CMD [ "/dev.sh"]