from __future__ import (absolute_import, division, print_function,
                        unicode_literals)

import datetime
import json
from builtins import *
from collections import namedtuple

from dateutil.parser import parse
from future import standard_library

standard_library.install_aliases()

MultiField = namedtuple('MultiField', [
    'Name',
    'Type',
    'IsNull',
    'BoolValue',
    'IntValue',
    'FloatValue',
    'StringValue',
    'DateTimeValue',
])


class GoRuntimeError(Exception):
    pass


def handle_exception(method, args, other=None):
    try:
        return method(*args)
    except RuntimeError as e:
        raise GoRuntimeError(
            '{0} raised on Go side while calling {1} with args {2} from {3}'.format(
                repr(e), repr(method), repr(args), repr(other)
            )
        )


def handle_records_json(records_json_string, session=None):
    try:
        records_json = json.loads(records_json_string)
    except ValueError as e:
        raise ValueError(
            '{0} raised {1} while parsing {2}'.format(
                session,
                e,
                repr(records_json_string)
            )
        )

    return [
        [
            MultiField(**column) for column in row
        ] for row in records_json
    ]


def handle_records(raw_records):
    records = []
    for raw_record in raw_records:
        record = []
        for column in raw_record:
            if column.Type == 'null':
                record += [None]
            elif column.Type == 'bool':
                record += [column.BoolValue]
            elif column.Type == 'int':
                record += [column.IntValue]
            elif column.Type == 'float':
                record += [column.FloatValue]
            elif column.Type == 'string':
                record += [column.StringValue]
            elif column.Type == 'datetime':
                try:
                    record += [
                        parse(column.DateTimeValue)
                    ]
                except ValueError as e:
                    raise ValueError(
                        'raised {0} while parsing {1}'.format(
                            e,
                            repr(column.DateTimeValue)
                        )
                    )

        records += [tuple(record)]

    return records
