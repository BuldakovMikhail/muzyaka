import sqlite3
import pandas as pd
import pickle

from abc import ABC, abstractmethod

class IDataBase(ABC):
    @abstractmethod
    def insert(self, id, data, event_id):
        pass

    @abstractmethod
    def get(self, id):
        pass

    @abstractmethod
    def delete(self, id, event_id):
        pass

    @abstractmethod
    def update(self, id, data, event_id):
        pass

    @abstractmethod
    def is_event_happened(self, event_id):
        pass

    @abstractmethod
    def convert_to_dataframe(self):
        pass

class DataBase(IDataBase):
    def __init__(self, db_path):
        self.db_path = db_path
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        query = """
        CREATE TABLE IF NOT EXISTS 
        tracks(
            id INTEGER NOT NULL,
            embedding BLOB NOT NULL
        )
        """
        cursor.execute(query)

        query = """
                CREATE TABLE IF NOT EXISTS 
                events(
                    event_id STRING NOT NULL
                )
                """
        cursor.execute(query)

        cursor.close()
        conn.close()

    def insert(self, id, data, event_id):
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        query = """
                INSERT INTO tracks(id, embedding) VALUES (?, ?)
                    """
        cursor.execute(query, (id, pickle.dumps(data)))

        query = """
                        INSERT INTO events(event_id) VALUES (?)
                            """
        cursor.execute(query, (event_id,))

        conn.commit()

        cursor.close()
        conn.close()

    def get(self, id):
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        query = """
        SELECT embedding FROM tracks WHERE id=?
        """
        cursor.execute(query, (id,))
        data = cursor.fetchone()

        cursor.close()
        conn.close()

        return data

    def delete(self, id, event_id):
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        query = """
                DELETE FROM tracks WHERE id=?
                """
        cursor.execute(query, (id,))
        query = """
                                INSERT INTO events(event_id) VALUES (?)
                                    """
        cursor.execute(query, (event_id,))

        conn.commit()

        cursor.close()
        conn.close()

    def update(self, id, data, event_id):
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        query = """
                       UPDATE tracks SET embedding=? WHERE id=?
                        """
        cursor.execute(query, (pickle.dumps(data), id))
        query = """
                                INSERT INTO events(event_id) VALUES (?,)
                                    """
        cursor.execute(query, (event_id,))

        conn.commit()

        cursor.close()
        conn.close()

    def convert_to_dataframe(self):
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        df = pd.read_sql_query("SELECT * FROM tracks", conn)
        df['embedding'] = df['embedding'].apply(lambda x: pickle.loads(x))
        cursor.close()
        conn.close()

        return df

    def is_event_happened(self, event_id):
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        cursor.execute("SELECT * FROM events WHERE event_id=?", (event_id,))

        return cursor.fetchone()