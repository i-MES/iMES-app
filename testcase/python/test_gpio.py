import time
import pytest
from helper import delay

@pytest.mark.UT
class Test_GPIO:
    def setup_class(self):
      pass

    def test_gpio_buff(self):
      """GPIO 测试"""
      assert delay()
      
    def test_gpio_flash(self):
      assert delay()

    def teardown_class(self):
      pass
