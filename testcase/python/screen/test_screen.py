import time
from random import randint
import pytest
from helper import delay

# @pytest.mark.UT
class Test_Screen:
    def setup_class(self):
      pass

    def test_screen_on(self):
      assert delay()
      
    def test_screen_off(self):
      assert delay()

    def test_screen_flash(self):
      assert delay()

    def teardown_class(self):
      pass