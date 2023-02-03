# coding: utf-8

"""


    Generated by: https://openapi-generator.tech
"""

import unittest
from unittest.mock import patch

import urllib3

import Formance
from Formance.paths.api_payments_payments_payment_id_metadata import patch  # noqa: E501
from Formance import configuration, schemas, api_client

from .. import ApiTestMixin


class TestApiPaymentsPaymentsPaymentIdMetadata(ApiTestMixin, unittest.TestCase):
    """
    ApiPaymentsPaymentsPaymentIdMetadata unit test stubs
        Update metadata  # noqa: E501
    """
    _configuration = configuration.Configuration()

    def setUp(self):
        used_api_client = api_client.ApiClient(configuration=self._configuration)
        self.api = patch.ApiForpatch(api_client=used_api_client)  # noqa: E501

    def tearDown(self):
        pass

    response_status = 204
    response_body = ''




if __name__ == '__main__':
    unittest.main()