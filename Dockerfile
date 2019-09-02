FROM python:2
WORKDIR /app
ADD ./install.sh install.sh
ADD ./requirements.txt requirements.txt
RUN ./install.sh
ENTRYPOINT ["./run.sh", "/repo"]
