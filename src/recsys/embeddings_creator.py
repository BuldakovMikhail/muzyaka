import random


class EmbeddingCreator:
    def __init__(self, client, emb_len):
        self.client = client
        self.emb_len = emb_len

    def create_embeddings(self, source, bucket):
        try:
            response = self.client.get_object(bucket, source)
            # Read data from response.
        finally:
            response.close()
            response.release_conn()

        return [random.random() for i in range(self.emb_len)]