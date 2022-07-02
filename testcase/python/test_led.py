import time
from random import randint
import pytest


@pytest.mark.UT
class Test_LED:
    def setup_class(self):
      pass

    def test_led_on(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5
      
    def test_led_off(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def test_led_flash(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def teardown_class(self):
      pass
