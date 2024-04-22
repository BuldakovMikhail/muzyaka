import json

from confluent_kafka import Consumer, KafkaException
from embeddings_creator import EmbeddingCreator
from database import DataBase
from model import Model
from minio import Minio
from flask import Flask, jsonify, request
from threading import Thread

#TODO
client = Minio()
emb_creator = EmbeddingCreator(client, 15)
bucket = "track-bucket"

db = DataBase("track.db")
model = Model(db, 10)

app = Flask(__name__)


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

                    message = json.loads(msg.value().decode('utf-8'))
                    id = int(message['track_id'])
                    event_id = message['event_id']
                    src = message['source']

                    if not db.is_event_happened(event_id):
                        if message['operation'] == "update":
                            embs = emb_creator.create_embeddings(src, bucket)
                            db.update(id, embs, event_id)
                            model.update_recs()
                        elif message['operation'] == "add":
                            embs = emb_creator.create_embeddings(src, bucket)
                            db.insert(id, embs, event_id)
                            model.update_recs()
                        elif message['operation'] == "delete":
                            db.delete(id, event_id)
                            model.update_recs()

                    consumer.commit()
        except KeyboardInterrupt:
            pass
        finally:
            consumer.close()  # не забываем закрыть соединение

    return basic_consume_loop


@app.route("/rec/<id>", methods=["GET"])
def get_recs(id):
    recs = model.get_recs(id, 10)
    return jsonify({"ids": recs})

def main():
    kafka_config = {
        'bootstrap.servers': 'localhost:29092',  # Список серверов Kafka
        'group.id': 'recsysgroup',  # Идентификатор группы потребителей
        'auto.offset.reset': 'earliest',  # Начальная точка чтения ('earliest' или 'latest')
        'enable.auto.commit': False
    }

    consumer_loop = create_consumer(kafka_config)
    thread = Thread(target=consumer_loop, args=(['sync-default-events'],))
    thread.daemon = True
    thread.start()

    app.run(host="0.0.0.0", port=12121)

if __name__ == '__main__':
    main()

