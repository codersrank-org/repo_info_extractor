import numpy as np
import copy


class IdentityMatcher:

    def __init__(self,
                 preprocessor,
                 model,
                 threshold=0.75,
                 debug=False):

        self._preprocessor = preprocessor
        self._model = model
        self._threshold = threshold
        self._debug = debug

        return

    def get_emails(self, start_seed, repo_commits_data):
        """
        Find all emails of the user based on string similarity.
        """

        seed = copy.deepcopy(start_seed)
        output_seed = copy.deepcopy(start_seed)

        while True:
            for commit_line in repo_commits_data:
                score = self._get_max_sim_score(seed, commit_line)
                if score > self._threshold:
                    name = commit_line[1]
                    email = commit_line[2]
                    if name not in output_seed["names"]:
                        output_seed["names"].append(name)
                    if email not in output_seed["emails"]:
                        output_seed["emails"].append(email)
            if self._debug:
                print("Run done.")
                print(output_seed)
                print(seed)
                print("----------------------------------------------")
            if output_seed == seed:
                return output_seed
            else:
                seed = copy.deepcopy(output_seed)

    def _get_sim_score_lists(self, seed, commit_stats):
        """
        Calculate the similarity between the names and emails found in the commit and the seed
        """

        full_name = self._preprocessor.transform(commit_stats[1])
        email = self._preprocessor.transform(commit_stats[2])

        user_name = seed.get("user_name")
        names = seed.get("names")
        emails = seed.get("emails")

        name_sim_scores = {s: self._calc_similarity(s, full_name) for s in names + [user_name]}
        email_sim_scores = {s: self._calc_similarity(s, email) for s in emails + [user_name]}

        return name_sim_scores, email_sim_scores

    def _get_max_sim_score(self, seed, commit_stats):
        """
        Get the maximum of the similarity scores.
        """

        n_scores, email_scores = self._get_sim_score_lists(seed, commit_stats)
        n_scores_list = list(n_scores.values())
        email_scores_list = list(email_scores.values())

        return np.max(n_scores_list + email_scores_list)

    def _calc_similarity(self, seed_string, input_vector):
        seed_vector = self._preprocessor.transform(seed_string)
        score = self._model.predict((seed_vector, input_vector))
        if isinstance(score, float):
            return score
        else:
            return score[0][0]