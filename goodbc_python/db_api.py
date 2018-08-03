from .rpc_session import create_session


class CursorClosedError(Exception):
    pass


class NoRowsToFetchError(Exception):
    pass


class ConnectionClosedError(Exception):
    pass


class Cursor(object):
    def __init__(self, connection, close_callback):
        self._connection = connection
        self._close_callback = close_callback

        self._closed = False
        self._rows_id = None
        self._result_id = None
        self._rowcount = -1

    def __repr__(self):
        return '<{}(connection={}) at {}>'.format(
            self.__class__.__name__,
            repr(self._connection),
            hex(id(self))
        )

    def __enter__(self):
        return self

    def __exit__(self, type, value, traceback):
        self.close()

    @property
    def rowcount(self):
        if self._closed:
            raise CursorClosedError('attempt to access .rowcount on closed cursor')

        if self._result_id is not None:
            result_id = self._result_id
            self._result_id = None
            self._rowcount = self._connection.session.rowcount(result_id)
            self._result_id = None

        return self._rowcount

    def _query(self, query):
        if self._closed:
            raise CursorClosedError('attempt to call ._query() on closed cursor')

        self._rows_id = self._connection.session.query(query)
        self._rowcount = -1

    def fetchall(self):
        if self._closed:
            raise CursorClosedError('attempt to call .fetchall() on closed cursor')

        if self._rows_id is None:
            raise NoRowsToFetchError('attempt to call .fetchall() before ._query() called')
        elif self._rows_id is False:
            return []

        rows_id = self._rows_id

        self._rows_id = False

        rows = self._connection.session.fetchall(rows_id)

        return rows

    def _execute(self, query):
        if self._closed:
            raise CursorClosedError('attempt to call ._execute() on closed cursor')

        self._result_id = self._connection.session.execute(query)
        self._rowcount = -1

    def execute(self, query):
        if self._closed:
            raise CursorClosedError('attempt to call .execute() on closed cursor')

        is_select = False
        for word in query.split():
            if word.upper() == 'SELECT':
                is_select = True
                break

        if is_select:
            self._query(query)
        else:
            self._execute(query)

    def close(self):
        self._closed = True

        self._close_callback()


class Connection(object):
    def __init__(self, data_source_name):
        self._session = create_session(data_source_name=data_source_name)

        self._cursor = None

        self._session.connect()

    def __repr__(self):
        return '<{}(session={}) at {}>'.format(
            self.__class__.__name__,
            repr(self._session),
            hex(id(self))
        )

    def __enter__(self):
        return self

    def __exit__(self, type, value, traceback):
        self.close()

    @property
    def session(self):
        return self._session

    def _close_cursor(self):
        self._cursor = None

    def cursor(self):
        if self._session is None:
            raise ConnectionClosedError('attempt to call .cursor() on a closed connection')

        return self._cursor if self._cursor is not None else Cursor(
            connection=self,
            close_callback=self._close_cursor
        )

    def close(self):
        if self._session is None:
            raise ConnectionClosedError('attempt to call .close() on a closed connection')

        if self._cursor is not None:
            self._cursor.close()

        self._session.close()
        self._session = None
