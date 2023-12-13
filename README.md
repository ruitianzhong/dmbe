## Driver Management Back End

### Setup

Modify the configuration in [config.template.yaml](./config.template.yaml).

After finishing all configuration,create config.yaml file and copy all the content from
[config.template.yaml](./config.template.yaml).

Make sure the executable is in same directory with config.yaml

Attention:

Don't leak your config.yaml to the public

#### MySQL

Configure MySQL.

example:

```yaml
db:
  username: your username
  password: your password
  dbname: database name
  address: localhost
  port: 3306
```

#### Server

example:

```yaml
app:
  host: 127.0.0.1
  port: 8080
```

#### Session Key

Configure the session key for gorilla/sessions,
and it must be keep secret.

```yaml
auth:
  session_key: your session key
```

### Build

#### Prerequisite

+ Download [Go](https://go.dev/dl/) on your machine.

+ If you have network problem,configure the GOPROXY:

```shell
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

+ Your database has been set up with the SQL in [sql/init.sql](sql/init.sql)

#### Command

```shell
cd dmbe
go build dmbe
```

### Run

Run in the background:

```shell
nohup ./dmbe &
```