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

    def ask_user_identity(self, repo):
        choices = []
        for key in repo.contributors:
            choices.append({
                'name': repo.contributors[key]['name'] + '-> ' + repo.contributors[key]['email'],
            })

        questions = [
            {
                'type': 'checkbox',
                'name': 'user_identity',
                'message': 'The following contributors were found in the repository. Select which ones you are. (With SPACE you can select more than one)',
                'choices': choices
            }
        ]

        return prompt(questions)