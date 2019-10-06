"""
Manual mem test.
1. Make a virtual env and install all deps
2. run `python setup.py install`
3. run `python mem_test.py`
4. Watch memory usage and make sure it doesn't get out of control
"""
from goodbc_python import Connection
from goodbc_python.common import GoRuntimeError
import traceback
import time
from guppy import hpy


print(hpy().heap())

err_count = 0


conn_str = """DRIVER={SQLite3};
                  Database=./testdb.sqlite;
            """.replace('\n', '').replace(' ', '')
connection = Connection(conn_str)
cursor = connection.cursor()
query = "select * from my_users;"

while True:
    try:
        # This just causes a connection to fail over and over
        cursor.execute(query)
        cursor.fetchall()
    except KeyboardInterrupt:
        break
    except (RuntimeError, GoRuntimeError):
        err_count += 1

print(hpy().heap())

print('done')


