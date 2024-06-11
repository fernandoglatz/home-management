FROM alpine:3 AS builder
ARG PLUGIN_VERSION=3.13.0
RUN apk --no-cache add wget && \
    wget "https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/v${PLUGIN_VERSION}/rabbitmq_delayed_message_exchange-${PLUGIN_VERSION}.ez" -O /tmp/plugin.ez


FROM rabbitmq:3-management

COPY --from=builder --chown=rabbitmq:rabbitmq /tmp/plugin.ez $RABBITMQ_HOME/plugins/rabbitmq_delayed_message_exchange-${PLUGIN_VERSION}.ez

RUN rabbitmq-plugins enable rabbitmq_mqtt && \
    rabbitmq-plugins enable rabbitmq_delayed_message_exchange