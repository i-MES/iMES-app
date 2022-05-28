import time
from random import randint
import pytest


@pytest.mark.UT
class Test_LED:
    def setup_class(self):
      pass

    def test_led_on(self):
      time.sleep(randint(2,5))
      
    def test_led_off(self):
      time.sleep(randint(2,5))

    def test_led_flash(self):
      time.sleep(randint(2,5))

    def teardown_class(self):
      pass
