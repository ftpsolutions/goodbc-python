from hamcrest import assert_that, equal_to

from goodbc_python import Connection


def test_odbc_connection():
    conn_str = """DRIVER={SQLite3};
                  Database=/workspace/testdb.sqlite;
            """.replace(
        "\n", ""
    ).replace(
        " ", ""
    )

    connection = Connection(conn_str)
    cursor = connection.cursor()

    query = "select * from my_users;"

    cursor.execute(query)

    records = cursor.fetchall()
    assert_that(len(records), equal_to(1))

    cursor.close()
    connection.close()
