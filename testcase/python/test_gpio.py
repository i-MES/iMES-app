import time
from random import randint
import pytest

@pytest.mark.UT
class Test_GPIO:
    def setup_class(self):
      pass

    def test_gpio_buff(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5
      
    def test_gpio_flash(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def teardown_class(self):
      pass
