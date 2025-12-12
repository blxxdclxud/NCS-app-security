FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    gcc \
    gdb \
    python3 \
    python3-pip \
    netcat \
    vim \
    && rm -rf /var/lib/apt/lists/*

RUN pip3 install pwntools

WORKDIR /app

COPY vulnerable-program.c /app/
COPY exploit.py /app/
COPY compile.sh /app/
COPY docker-entrypoint.sh /app/

RUN chmod +x /app/compile.sh /app/exploit.py /app/docker-entrypoint.sh

RUN /app/compile.sh

EXPOSE 4444

ENTRYPOINT ["/app/docker-entrypoint.sh"]

# Default: interactive shell (user runs exploit manually)
# To auto-run exploit: docker run --rm -it --privileged <image> python3 exploit.py local
CMD ["/bin/bash"]
