import time
from random import randint
import pytest

@pytest.mark.UT
class Test_GPIO:
    def setup_class(self):
      pass

    def test_gpio_buff(self):
      time.sleep(randint(2,5))
      
    def test_gpio_flash(self):
      time.sleep(randint(2,5))

    def teardown_class(self):
      pass
