
from abc import ABC, abstractmethod
import numpy as np

import pickle


class MyModel(ABC):
    """
    Abstract method for the distance calculation.
    """
    @abstractmethod
    def fit(self, *args):
        pass

    @abstractmethod
    def predict(self, input_pair, *args):
        pass


class DistanceModel(MyModel):
    """
    Model based on a TfIdf vectorizer and cosine similarity
    """

    def __init__(self,
                 vectorizer_path):

        with open(vectorizer_path, "rb") as f:
            self._vectorizer = pickle.load(f)

    def fit(self, *args):
        pass

    def predict(self, input_pair, *args):

        x_vec = self._vectorizer.transform(input_pair).toarray()

        return np.dot(*x_vec)


