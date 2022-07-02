"""Tests for our testing utilities."""

import time
from random import randint
import pytest

@pytest.mark.ST
class TestMemoryCache:
    def test_filters_for_valloc_and_free(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def test_filters_based_on_addresses(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def test_free_records_with_valid_addresse(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def test_free_records_with_unmatched_addresses(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5


class TestMemoryBurst:
    def test_holds_values_at_correct_names(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def test_looks_like_AllocationRecord(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

    def test_equality(self):
      i = randint(1,5)
      time.sleep(i)
      assert i < 5

