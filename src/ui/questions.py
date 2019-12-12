from whaaaaat import style_from_dict, Token, prompt, print_json, default_style, Separator
import sys

def sort_by_checked_and_email(d):
    """ Sort by author match, then email. """
    email=d['name'].split(' -> ', 1)[1].lower()
    checked=d['checked']
    order = '1'+email
    if(checked==True):
        order = '0'+email
    return order

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

    def ask_user_identity(self, authors, err, default_email=''):
        choices = []
        for name, email in authors:
            checked = email == default_email
            choices.append({
                'name': name + ' -> ' + email,
                'checked': checked
            })
        choices.sort(key=sort_by_checked_and_email)
        message = 'The following contributors were found in the repository. \
            Select which ones you are. (With SPACE you can select more than one)'
        if err:
            message = "%s [ERROR] %s" % (message, err)

        questions = [
            {
                'type': 'checkbox',
                'name': 'user_identity',
                'message': message,
                'choices': choices
            }
        ]

        return prompt(questions)

    def ask_which_repos(self, repos):
        choices = []
        for repo in repos:
            choices.append({
                'name': repo
            })
            
        print("We found the following repos in the chosen path")

        sorted_choices = sorted(choices, key= lambda x: x['name'].lower())
        questions = [
            {
                'type': 'checkbox',
                'name': 'chosen_repos',
                'message': 'Select which ones you want to analyze (With SPACE you can select more than one)',
                'choices': sorted_choices
            }
        ]

        return prompt(questions)

    def query_yes_no(self, question, default="yes"):
        """Ask a yes/no question via raw_input() and return their answer.

        "question" is a string that is presented to the user.
        "default" is the presumed answer if the user just hits <Enter>.
            It must be "yes" (the default), "no" or None (meaning
            an answer is required of the user).

        The "answer" return value is True for "yes" or False for "no".
        """
        valid = {"yes": True, "y": True, "ye": True,
                 "no": False, "n": False}
        if default is None:
            prompt = " [y/n] "
        elif default == "yes":
            prompt = " [Y/n] "
        elif default == "no":
            prompt = " [y/N] "
        else:
            raise ValueError("invalid default answer: '%s'" % default)
        while True:
            print(question + prompt)
            choice = sys.stdin.readline().lower().strip()
            if default is not None and choice == '':
                return valid[default]
            elif choice in valid:
                return valid[choice]
            else:
                sys.stdout.write("Please respond with 'yes' or 'no' "
                                 "(or 'y' or 'n').\n")
