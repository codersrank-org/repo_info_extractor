"""
Simple function for timeouting a long running process. WARNING! As it uses the signal package, it won't work
under Windows.
"""
import signal
from contextlib import contextmanager




@contextmanager
def timeout(time, repo_working_dir):
    # Register a function to raise a TimeoutError on the signal.
    signal.signal(signal.SIGALRM, raise_timeout)
    # Schedule the signal to be sent after ``time``.
    signal.alarm(time)

    try:
        yield
    except TimeoutError:
        print("{} timeouted.".format(repo_working_dir))
    finally:
        # Unregister the signal so it won't be triggered
        # if the timeout is not reached.
        signal.signal(signal.SIGALRM, signal.SIG_IGN)


def raise_timeout(signum, frame):
    raise TimeoutError
