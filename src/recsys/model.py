from sklearn.preprocessing import StandardScaler
from sklearn.neighbors import NearestNeighbors
import pandas as pd


from database import IDataBase

class Model:
    def __init__(self, db:IDataBase, n_neighbors:int):
        self.nn = None
        self.scaled_data = None
        self.embedding_indexes = None
        self.df = None
        self.scaler = None
        self.db = db
        self.n_neighbors = n_neighbors
        self.update_recs()

    def get_recs(self, id:int, offset:int, limit:int):
        res = self.nn.kneighbors(
            self.scaled_data[self.df["id"] == id],
            min(limit + offset, len(self.scaled_data)),
            False
        )
        return self.df.iloc[res[0][offset:]]["id"]

    def update_recs(self):
        self.scaler = StandardScaler()
        self.df = self.db.convert_to_dataframe()

        if len(self.df['embedding']) == 0:
            return

        self.embedding_indexes = [f"{i}" for i in range(len(self.df['embedding'][0]))]

        self.df[self.embedding_indexes] = pd.DataFrame(
            self.df['embedding'].tolist(),
            index=self.df.index
        )

        self.scaled_data = self.scaler.fit_transform(self.df[self.embedding_indexes])
        self.nn = NearestNeighbors(n_neighbors=self.n_neighbors)
        self.nn.fit(self.scaled_data)

