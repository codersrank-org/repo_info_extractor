FROM python:3.7-slim AS build
WORKDIR /app
# Make sure we use the virtualenv:
ENV VIRTUAL_ENV=/opt/venv
ENV PATH="$VIRTUAL_ENV/bin:$PATH"
COPY . /app
RUN apt-get update \
    && apt-get install -y --no-install-recommends build-essential gcc \
    && python -m venv $VIRTUAL_ENV \
    && pip install --upgrade pip \
    && pip install -r requirements.txt \
    && make test

FROM python:3.7-slim AS runtime
WORKDIR /app
# Make sure we use the virtualenv:
ENV VIRTUAL_ENV=/opt/venv
ENV PATH="$VIRTUAL_ENV/bin:$PATH"
COPY --from=build $VIRTUAL_ENV $VIRTUAL_ENV
RUN apt-get update \
    && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["./run.sh"]
