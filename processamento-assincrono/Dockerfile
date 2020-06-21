FROM rabbitmq:management AS build

RUN apt-get update -y \
    && apt-get install wget -y


FROM build

WORKDIR /plugins
RUN wget https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/v3.8.0/rabbitmq_delayed_message_exchange-3.8.0.ez \
    && rabbitmq-plugins enable --offline rabbitmq_delayed_message_exchange rabbitmq_shovel_management
