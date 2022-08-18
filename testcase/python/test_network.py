import time
from random import randint
import pytest
from helper import delay


@pytest.mark.UT
class Test_Network:
    def setup_class(self):
      pass

    def test_tcp(self):
      assert delay()
      
    def test_udp(self):
      assert delay()

    def test_socket(self):
      assert delay()

    def test_icmp(self):
      assert delay()
      
    def test_4G(self, entity):
      print(entity)
      assert hasattr(entity, "ip") == True
      assert entity.ip.startswith("199") == True

    def teardown_class(self):
      pass
