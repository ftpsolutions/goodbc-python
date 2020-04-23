import datetime
import os
import sys

from guppy import hpy

from goodbc_python.common import GoRuntimeError
from goodbc_python.db_api import Connection

if __name__ == "__main__":
    connection_string = "driver={{FreeTDS}}; server={}; port={}; database={}; uid={}; pwd={};".format(
        "host.docker.internal",
        11433,
        "master",
        "sa",
        os.getenv("MSSQL_PASSWORD", "unset"),
    )

    print(hpy().heap())

    print("\nNOTE: a . is a successful query, a ! is a failed query\n")

    query_count = 0
    error_count = 0
    connection = None

    started = datetime.datetime.now()
    while True:
        try:
            connection = Connection(connection_string)
            cursor = connection.cursor()
            query = "SELECT TOP 2000 * FROM locations"
            cursor.execute(query)
            cursor.fetchall()
            sys.stdout.write(".")
            sys.stdout.flush()
            query_count += 1
        except KeyboardInterrupt:
            break
        except (RuntimeError, GoRuntimeError) as e:
            sys.stdout.write("!")
            sys.stdout.flush()
            error_count += 1
        finally:
            if connection is not None:
                connection.close()
    stopped = datetime.datetime.now()
    duration = stopped - started
    duration_seconds = duration.total_seconds()

    print("\n\n{}\n".format(hpy().heap()))

    for name, value in [("queries", query_count), ("errors", error_count)]:
        print(
            "{} {} in {}; {} per second".format(
                value,
                name,
                duration,
                (value / duration_seconds) if duration_seconds > 0 else -1,
            )
        )
