## goodbc-python

The purpose of this module is to provide a Python interface to the Golang [goodbc](https://github.com/alexbrainman/odbc) module.

It was made very easy with the help of the Golang [gopy](https://github.com/go-python/gopy) module.

#### Limitations

* Python command needs to be prefixed with GODEBUG=cgocheck=0 (or have that in the environment)

#### Prerequisites

* Go 1.13
* Python 2.7+
* pip
* virtualenvwrapper
* pkgconfig/pkg-config
* unixodbc 
    * Linux: ```apt-get install unixodbc unixodbc-dev freetds-bin freetds-dev```
    * OSX: ```brew install freetds --with-unixodbc```

#### Installation (for prod)
* ```python setup.py install``` 

#### Making a python wheel install file (for distribution)
* ```python setup.py bdist_wheel``` 

#### Setup (for dev)
Ensure pkg-config is installed
Ensure unixodbc is installed

* ```mkvirtualenvwrapper -p (/path/to/pypy) goodbc-python``` 
* ```pip install -r requirements-dev.txt```
* ```./build.sh```
* ```GODEBUG=cgocheck=0 py.test -v```

#### What's worth knowing if I want to further the development?

* gopy doesn't like Go interfaces; so make sure you don't have any public (exported) interfaces
    * this includes a struct with a public property that may eventually lead to an interface


#### Example Python usage

To create an ODBC session in Python do the following:

```
from goodbc_python import Connection

ip = "127.0.0.1"
port = 5432
database = "test"
username = "test"
password = "test"

conn_str = """
            DRIVER={FreeTDS};
            TDS_VERSION=8.0;
            SERVER=%s;
            Port=%i;
            DATABASE=%s;
            UID=%s;
            PWD=%s;
        """ % (
    ip, port, database,
    username, password
).replace('\n', '').replace(' ', '')

connection = Connection(conn_str)
cursor = connection.cursor()

query = "SELECT NOW()"

cursor.execute(query)

records = cursor.fetchall()

print("Records:")
print(records)

cursor.close()
connection.close()
```

This seems to leak quite badly when trying to connect and query a bad IP address when using the FreeTDS driver.

FreeTDS v1.1.17, their latest stable but the problem persists. It could very well be an interaction between
goodbc and the driver. Not sure at this stage.

## To develop / run the tests

    MOUNT_WORKSPACE=1 ./test.sh bash
    ./build.sh
    py.test
    
## To test the sdist package

    py.test