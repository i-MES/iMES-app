import pytest

class Entity(object):
  def __init__(self, ip):
    self.ip = ip
    self.name = "iphone"
    self.status = "ready to test"

  def connect(self):
    # connect to entity from pc
    pass

@pytest.fixture(scope="session",params=None)
def entity() -> Entity:
  return create_entity("199.33.33.33")

def create_entity(ip):
  """创建被测实体

  :param ip: ip 地址
  :type ip: string
  :return: 被测实体实例化对象
  :rtype: Entity
  """
  e = Entity(ip)
  e.connect()
  return e

@pytest.fixture(scope="session",params=None)
def foo():
  return create_foo()
  
def create_foo():
  pass

def create_bar():
  pass