from whaaaaat import style_from_dict, Token, prompt, print_json, default_style, Separator

class Questions:

    def ask_primary_remote_url(self, repo):
        '''
        Promots the user the possible remote ULRs and ask her to select the primary
        '''
        choices = []
        for remote in repo.original_remotes:
            choices.append(remote + ': ' + repo.original_remotes[remote])
        questions = [
            {
                'type': 'list',
                'name': 'remote_repo',
                'message': 'Cannot find remote origin. Select which one is the primary remote URL.',
                'choices': choices
            }
        ]

        return prompt(questions)