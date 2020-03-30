from abc import ABC, abstractmethod

import re
import unicodedata
import numpy as np


class Preprocessor(ABC):

    def __init__(self, domain_blacklist = None):

        self._domain_blacklist = domain_blacklist

    @abstractmethod
    def transform(self, input_string):
        raise NotImplementedError

    @staticmethod
    def strip_accents(text):
        """
        Strip accents from input String.

        :param text: The input string.
        :type text: String.

        :returns: The processed String.
        :rtype: String.
        """
        text = unicodedata.normalize('NFD', text)
        text = text.encode('ascii', 'ignore')
        text = text.decode("utf-8")

        return str(text)

    def text_to_id(self, text):
        """
        Convert input text to id.

        :param text: The input string.
        :type text: String.

        :returns: The processed String.
        :rtype: String.
        """
        text = self.strip_accents(text.lower())
        text = re.sub(r'[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+', "", text)
        text = re.sub(r'[0-9abcdef]{8}-[0-9abcdef]{4}-[0-9abcdef]{4}-[0-9abcdef]{4}-[0-9abcdef]{12}', "", text)
        text = re.sub(r'[0-9]+', '0', text)
        text = re.sub('[^a-zA-Z0@ ._-]', '', text)
        return text


class DistancePreprocessor(Preprocessor):

    def __init__(self, domain_blacklist=None, shortlog=None):
        super().__init__(domain_blacklist)
        self._domain_blacklist = domain_blacklist
        self.shortlog = shortlog
        if self.shortlog:
            self.update_blacklist()

        if domain_blacklist:
            self._blacklist_regex = re.compile("|".join(self._domain_blacklist))
        else:
            self._blacklist_regex = None
        return

    def transform(self, input_string):
        if self._blacklist_regex:
            res = self._blacklist_regex.sub("", input_string)
        else:
            res = input_string
        return self.text_to_id(res)

    @staticmethod
    def __extract_domain(email):
        try:
            return email.split(sep="@")[1]
        except IndexError:
            return ""
        except AttributeError:
            return ""

    @staticmethod
    def __calc_bound(domain_counts):
        unique_domain_count = domain_counts.shape[0]
        if unique_domain_count <= 5:
            return 1
        elif 5 < unique_domain_count <= 10:
            return 3
        else:
            return int(np.ceil(0.01*unique_domain_count))

    def update_blacklist(self):

        additional_domains = self.shortlog["email"].apply(self.__extract_domain)
        srs = additional_domains.value_counts()
        bound = self.__calc_bound(srs)
        additional_domains = list(srs[0:bound].index)
        if self._domain_blacklist is None:
            self._domain_blacklist = list(additional_domains)
            return

        for domain in additional_domains:
            if domain not in self._domain_blacklist:
                self._domain_blacklist.append(domain)
        return

