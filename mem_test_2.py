"""
Manual mem test.
# See ./test_mem.sh

This one appears to leak! Things that have been tried:
- Using FreeTDS v1.1.17
"""
from goodbc_python import Connection
from goodbc_python.common import GoRuntimeError
import traceback
import time
from guppy import hpy


print(hpy().heap())

err_count = 0


while True:
    try:
        # This just causes a connection to fail over and over
        conn_str = """DRIVER={FreeTDS};
                      TDS_VERSION=8.0;
                      SERVER=10.10.10.200;
                      Port=1433;
                      DATABASE=something;
                      UID=aa;
                      PWD=bb;
                    """.replace('\n', '').replace(' ', '')
        connection = Connection(conn_str)
        cursor = connection.cursor()
        query = "select * from my_users;"
        cursor.execute(query)
        cursor.fetchall()
    except KeyboardInterrupt:
        break
    except (RuntimeError, GoRuntimeError):
        print(traceback.format_exc())
        err_count += 1
    finally:
        connection.close()

print(hpy().heap())

print('done')


