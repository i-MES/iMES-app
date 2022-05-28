import time
from random import randint
import pytest


@pytest.mark.UT
class Test_Network:
    def setup_class(self):
      pass

    def test_tcp(self):
      time.sleep(randint(2,5))
      
    def test_udp(self):
      time.sleep(randint(2,5))

    def test_socket(self):
      time.sleep(randint(2,5))

    def test_icmp(self):
      time.sleep(randint(2,5))
      
    def teardown_class(self):
      pass
