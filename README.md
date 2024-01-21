## Infotecs trainee task

RestAPI Service for processing payments

To use uuid in postgres you need to install extension:
```genericsql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

To start application run the command:
```shell
make compose-up
```

Database structure described in file ```./migration/20240120_trainee.app.sql```

Requests to test api is located in ```./test.http```