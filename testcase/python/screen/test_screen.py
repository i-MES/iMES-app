import time
from random import randint
import pytest


@pytest.mark.UT
class Test_Screen:
    def setup_class(self):
      pass

    def test_screen_on(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5
      
    def test_screen_off(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def test_screen_flash(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def teardown_class(self):
      pass
