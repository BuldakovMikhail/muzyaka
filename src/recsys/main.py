import json

from confluent_kafka import Consumer, KafkaException
from embeddings_creator import EmbeddingCreator
from database import DataBase
from model import Model
import minio.credentials
from minio import Minio
from flask import Flask, jsonify, request
from threading import Thread


MIN_PAGE_SIZE = 10
MAX_PAGE_SIZE = 100

creds = minio.credentials.StaticProvider("minioadmin", "minioadmin", "")
client = Minio("localhost:9000", credentials=creds, secure=False)
bucket = "track-bucket"

emb_creator = EmbeddingCreator(client, 15)
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
                    print(msg.value().decode('utf-8'))
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


@app.route("/rec", methods=["GET"])
def get_recs():
    id = request.args.get('id', type=int)
    page = request.args.get('page', type=int)
    page_size = request.args.get('page_size', type=int)

    if page <= 0:
        page = 1

    if page_size < MIN_PAGE_SIZE:
        page_size = MIN_PAGE_SIZE
    elif page_size > MAX_PAGE_SIZE:
        page_size = MAX_PAGE_SIZE

    offset = (page - 1) * page_size

    recs = model.get_recs(id, offset, page_size)
    # print(list(recs))
    return jsonify({"ids": list(recs)})

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

