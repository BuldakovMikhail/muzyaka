from confluent_kafka import Consumer, KafkaException

def create_consumer(config):
    consumer = Consumer(config)
    def basic_consume_loop(topics):
        try:
            # подписываемся на топик
            consumer.subscribe(topics)
            while True:
                msg = consumer.poll(timeout=1.0)  # ожидание сообщения
                if msg is None:  # если сообщений нет
                    continue
                if msg.error():  # обработка ошибок
                    raise KafkaException(msg.error())
                else:
                    # действия с полученным сообщением
                    # TODO: Place model and db save here
                    print(f"Received message: {msg.value().decode('utf-8')}")
                    consumer.commit()
        except KeyboardInterrupt:
            pass
        finally:
            consumer.close()  # не забываем закрыть соединение

    return basic_consume_loop

def main():
    kafka_config = {
        'bootstrap.servers': 'localhost:29092',  # Список серверов Kafka
        'group.id': 'recsysgroup',  # Идентификатор группы потребителей
        'auto.offset.reset': 'earliest',  # Начальная точка чтения ('earliest' или 'latest')
        'enable.auto.commit': False
    }

    consumer_loop = create_consumer(kafka_config)
    consumer_loop(['sync-default-events'])


if __name__ == '__main__':
    main()
