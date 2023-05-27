# 2023 FMNX team.
# Use of this code is governed by GNU General Public License.
# Additional information can be found on official web page: https://fmnx.su/
# Contact email: help@fmnx.su

FROM docker.io/golang:latest as build

WORKDIR /src

COPY go.mod /src
COPY go.sum /src
RUN go mod download

COPY . /src/
RUN go build -o packbin .

FROM archlinux/archlinux:base-devel

LABEL maintainer="dancheg <dancheg@fmnx.su>"

COPY --from=build /src/packbin /usr/bin/pack

ENTRYPOINT ["pack"]
CMD ["-h"]