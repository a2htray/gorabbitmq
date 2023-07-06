Examples of RabbitMQ command line tool `rabbitmqctl`
===============================

# Introduction

# Sub commands

## Users

### add_user

`add_user` subcommand could create a new user in the internal database, who has no permission for any virtual hosts by default.
Its format is as follows.

```bash
rabbitmqctl [--node <node>] [--longnames] [--quiet] add_user <username> <password>
```

**basic usage**

```bash
$ rabbitmqctl add_user admin admin123
Adding user "admin" ...
Done. Don't forget to grant the user permissions to some virtual hosts! See 'rabbitmqctl help set_permissions' to learn more.
```

### change_password

The `change_password` subcommand is simple and its format is as follows.

```bash
rabbitmqctl [--node <node>] [--longnames] [--quiet] change_password <username> <password>
```

**basic usage**

```bash
$ rabbitmqctl change_password admin newpassword123
Changing password for user "admin" ...
```

### set_user_tags

RabbitMQ allows developers to set or update a specific user to have certain tags by using `set_user_tags`. There are five 
tags can be used.

* administrator
* impersonator
* management
* monitoring
* policymaker

And, its format is as follows.

```bash
rabbitmqctl [--node <node>] [--longnames] [--quiet] set_user_tags <username> <tag> [...]
```

**set the administrator tag to admin**

```bash
$ rabbitmqctl set_user_tags admin administrator
Setting tags for user "admin" to [administrator] ...
```

## Virtual hosts

### add_vhost

The `add_vhost` subcommand could create a virtual host from the terminal. And its format is as follows.

```bash
rabbitmqctl [--node <node>] [--longnames] [--quiet] add_vhost <vhost> [--description <description> --tags "<tag1>,<tag2>,<...>" --default-queue-type <quorum|classic|stream>]
```

The `vhost` argument is the virtual host name which would be created, and developers can add description and set tags by using 
`--description` and `--tags` options respectively.

**add a virtual host**

```bash
$ rabbitmqctl add_vhost goapp-vhost --description 'virtual host for go test application' --tags goapp
Adding vhost "goapp-vhost" ...
```

## Access control

### set_permissions




