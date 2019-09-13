import requests
import json
from ui.progress import progress


def uploadRepo(filename):
    url = "https://grpcgateway.codersrank.io/candidate/privaterepo/Upload"
    fin = open(filename, 'rb')
    files = {'file': fin}
    try:
        print('Uploading '+filename+" ...")
        r = requests.post(url, files=files)
        json_data = json.loads(r.text)
        r.raise_for_status()
        return json_data
    except requests.exceptions.RequestException as err:
        showError(err)
    except requests.exceptions.HTTPError as err:
        showError(err)
    finally:
        fin.close()


def showError(err):
    print(err)
    print('.')
    print('.')
    print("Unable to upload the file, please try uploading it manually.")
