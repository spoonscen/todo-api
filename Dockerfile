FROM golang:onbuild
RUN go get github.com/codegangsta/gin
CMD ["gin", "-a", "8070", "-i"]
EXPOSE 3000