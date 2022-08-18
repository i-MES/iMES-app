import time
from random import randint

def delay():
  max = 8
  i = randint(1,max)
  time.sleep( 5 if i>5 else i)
  return i < max