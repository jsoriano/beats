FROM debian:9.6

RUN apt-get update && \
    apt-get install -y libvirt-daemon libvirt-clients dmidecode && \
    rm -rf /var/lib/apt/lists/*

RUN adduser --system --no-create-home --group libvirt-qemu
RUN addgroup libvirt
RUN addgroup libvirtd

COPY libvirtd.conf /etc/libvirtd/libvirtd.conf

CMD libvirtd --listen --config /etc/libvirtd/libvirtd.conf
