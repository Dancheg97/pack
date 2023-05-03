FROM archlinux/archlinux:base-devel

LABEL maintainer="dancheg97 <dancheg97@fmnx.io>"
LABEL source="https://fmnx.io/core/pack"

RUN pacman -Syu --needed --noconfirm git pacman-contrib wget go

ARG user=pack
RUN useradd --system --create-home $user
RUN echo "$user ALL=(ALL:ALL) NOPASSWD:ALL" > /etc/sudoers.d/$user
USER $user
WORKDIR /home/$user

RUN git clone  https://fmnx.io/core/pack && cd pack && makepkg --noconfirm -sfri
RUN sudo mv /home/$user/pack/*.pkg.tar.zst /var/cache/pacman/pkg
RUN sudo rm -r /home/$user/pack
RUN sudo rm -r /home/$user/go
RUN sudo pacman --noconfirm -R wget go

ENTRYPOINT ["pack"]
CMD ["help"]