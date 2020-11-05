
from abc import ABC, abstractmethod
import numpy as np

import pickle


class MyModel(ABC):
    """
    Abstract method for the distance calculation.
    """
    @abstractmethod
    def fit(self, *args):
        """
        Perform a function.

        Args:
            self: (todo): write your description
        """
        pass

    @abstractmethod
    def predict(self, input_pair, *args):
        """
        Predict the result.

        Args:
            self: (array): write your description
            input_pair: (array): write your description
        """
        pass


class DistanceModel(MyModel):
    """
    Model based on a TfIdf vectorizer and cosine similarity
    """

    def __init__(self,
                 vectorizer_path):
        """
        Initialize vectorizer.

        Args:
            self: (todo): write your description
            vectorizer_path: (str): write your description
        """

        with open(vectorizer_path, "rb") as f:
            self._vectorizer = pickle.load(f)

    def fit(self, *args):
        """
        Perform a function.

        Args:
            self: (todo): write your description
        """
        pass

    def predict(self, input_pair, *args):
        """
        Predict the vectorizer.

        Args:
            self: (array): write your description
            input_pair: (array): write your description
        """

        x_vec = self._vectorizer.transform(input_pair).toarray()

        return np.dot(*x_vec)


