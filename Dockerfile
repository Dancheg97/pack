# Copyright 2023 FMNX team.
# Use of this code is governed by GNU General Public License.
# Additional information can be found on official web page: https://fmnx.io/
# Contact email: help@fmnx.io

FROM archlinux/archlinux:base-devel

LABEL maintainer="dancheg97 <help@fmnx.io>"
LABEL source="https://fmnx.io/core/pack"

RUN pacman -Syu --needed --noconfirm git pacman-contrib wget go

RUN useradd --system --create-home pack
RUN echo "pack ALL=(ALL:ALL) NOPASSWD:ALL" > /etc/sudoers.d/pack
USER pack
WORKDIR /home/pack

COPY . /home/pack/pack
RUN sudo chmod a+rwx -R /home/pack
RUN cd pack && makepkg --noconfirm -sfri
RUN sudo mv /home/pack/pack/*.pkg.tar.zst /var/cache/pacman/pkg
RUN sudo rm -r /home/pack/pack
RUN sudo rm -r /home/pack/go
RUN sudo pacman --noconfirm -R wget go

ENTRYPOINT ["pack"]
CMD ["-h"]