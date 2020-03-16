"""
Simple function for timeouting a long running process.
"""

import _thread as thread
import sys


def timeout():
    # print to stderr, unbuffered in Python 2.
    sys.stderr.flush()  # Python 3 stderr is likely buffered.
    thread.interrupt_main()  # raises KeyboardInterrupt
