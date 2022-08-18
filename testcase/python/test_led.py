import time
from random import randint
import pytest
from helper import delay


@pytest.mark.UT
class Test_LED:
    def setup_class(self):
      pass

    def test_led_on(self):
      assert delay()
      
    def test_led_off(self):
      assert delay()

    def test_led_flash(self):
      assert delay()

    def teardown_class(self):
      pass
