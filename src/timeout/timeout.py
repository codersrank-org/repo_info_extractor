"""
Simple function for timeouting a long running process. WARNING! As it uses the signal package, it won't work
under Windows.
"""
import signal
from contextlib import contextmanager
import shutil


@contextmanager
def timeout(time, repo_working_dir):
    # Register a function to raise a TimeoutError on the signal.
    signal.signal(signal.SIGALRM, raise_timeout)
    # Schedule the signal to be sent after ``time``.
    signal.alarm(time)

    try:
        yield
    except TimeoutError:
        try:
            shutil.rmtree(repo_working_dir)
            print("{} timeouted, deleted files successfully.".format(repo_working_dir))
        except (PermissionError, NotADirectoryError, Exception) as e:
            print("{} timeouted, deletion failed with {}".format(repo_working_dir, e))

    finally:
        # Unregister the signal so it won't be triggered
        # if the timeout is not reached.
        signal.signal(signal.SIGALRM, signal.SIG_IGN)


def raise_timeout(signum, frame):
    raise TimeoutError
