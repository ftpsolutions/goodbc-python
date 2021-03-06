from __future__ import unicode_literals
from __future__ import print_function
from __future__ import division
from __future__ import absolute_import
from future import standard_library
standard_library.install_aliases()
from builtins import *
import unittest
import gc

from hamcrest import assert_that, equal_to, greater_than_or_equal_to
from mock import call, patch

from .rpc_session import RPCSession, create_session


class SessionTest(unittest.TestCase):
    def setUp(self):
        self._subject = RPCSession(
            session_id=0,
        )

    @patch('goodbc_python.rpc_session.RPCConnect')
    def test_connect(self, rpc_call):
        rpc_call.return_value = None

        assert_that(
            self._subject.connect(),
            equal_to(None)
        )

        assert_that(
            rpc_call.mock_calls,
            equal_to([
                call(0)
            ])
        )

    @patch('goodbc_python.rpc_session.RPCClose')
    def test_close(self, rpc_call):
        rpc_call.return_value = None

        assert_that(
            self._subject.close(),
            equal_to(None)
        )

        assert_that(
            rpc_call.mock_calls,
            equal_to([
                call(0)
            ])
        )

    @patch('goodbc_python.rpc_session.RPCClose')
    def test_del(self, rpc_call):
        rpc_call.return_value = None

        del(self._subject)

        gc.collect()

        assert_that(
            len(rpc_call.mock_calls),
            greater_than_or_equal_to(1)
        )


class ConstructorsTest(unittest.TestCase):
    @patch('goodbc_python.rpc_session.RPCSession')
    @patch('goodbc_python.rpc_session._new_session')
    def test_create_session(self, go_session_constructor, py_session_constructor):
        go_session_constructor.return_value = -1

        subject = create_session(
            data_source_name=u'some_data_source_name'
        )

        assert_that(
            go_session_constructor.mock_calls,
            equal_to([
                call('some_data_source_name')
            ])
        )

        assert_that(
            py_session_constructor.mock_calls,
            equal_to([
                call(data_source_name=u'some_data_source_name', session_id=-1)
            ])
        )

        assert_that(
            subject,
            equal_to(
                py_session_constructor()
            )
        )
