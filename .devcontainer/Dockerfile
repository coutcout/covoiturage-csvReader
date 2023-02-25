FROM ubuntu:23.04

COPY setup.sh .
ENV PATH "$PATH:/usr/local/go/bin:/home/user/go/bin:/home/user/.cargo/bin"
RUN bash setup.sh
COPY postCreate.sh /home/user
USER user
