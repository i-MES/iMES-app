import time
from random import randint
import pytest


@pytest.mark.UT
class Test_Network:
    def setup_class(self):
      pass

    def test_tcp(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5
      
    def test_udp(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def test_socket(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def test_icmp(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5
      
    def teardown_class(self):
      pass
