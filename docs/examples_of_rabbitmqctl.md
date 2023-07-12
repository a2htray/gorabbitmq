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

The `set_permissions` subcommand could be used to set user's permissions for the considered virtual host. The permissions 
consist of *configure*, *write* and *read*. It uses regular expressions to grant user the permissions, and the `.*` regular
expression means all the permission. The format of `set_permissions` subcommand is as follows.

```bash
rabbitmqctl [--node <node>] [--longnames] [--quiet] set_permissions [--vhost <vhost>] <username> <conf> <write> <read>
```

The `--vhost` option is used to set virtual host which the `<username>` could access to.

**set all permissions to a user**

```bash
$ rabbitmqctl set_permissions -p goapp-vhost goadmin ".*" ".*" ".*"
Setting permissions for user "goadmin" in vhost "goapp-vhost" ...
```

## Policies

### set_policy

The `set_policy` subcommand is one of RabbitMQ command line tool, which is used to set or update a policy. A policy means
a set of rules used to define queue behavior. The format is as follows.

```bash
rabbitmqctl [--node <node>] [--longnames] [--quiet] set_policy [--vhost <vhost>] [--priority <priority>] [--apply-to <apply-to>] <name> <pattern> <definition>
```

`<name>` is a unique policy name, `<pattern>` is a regular expression pattern that will be used to match queues, exchanges, etc.
`<definition>` is the policy definition which format must be JSON.

**set TTL to all queues**

```bash
$ rabbitmqctl set_policy ttl-3s "app.*" '{"message-ttl": 3000}' --apply-to queues --vhost goapp-vhost
Setting policy "ttl-3s" for pattern "app.*" to "{"message-ttl": 3000}" with priority "0" for vhost "goapp-vhost" ...
```

## Monitoring, observability and health checks

### list_queues

The `list_queues` subcommand can be used to statistic queue information that is showed in `Table` style. The command syntax
is as follows.

```bash
rabbitmqctl [--node <node>] [--longnames] [--quiet] list_queues [--vhost <vhost>] [--online] [--offline] [--local] [--no-table-headers] [<column>, ...] [--timeout <timeout>]
```

The `<column>` argument is the column set of the output table and must be one of arguments, auto_delete, consumer_capacity, 
consumer_utilisation, consumers, disk_reads, disk_writes, durable, effective_policy_definition, exclusive, exclusive_consumer_pid, 
exclusive_consumer_tag, head_message_timestamp, leader, members, memory, message_bytes, message_bytes_persistent, message_bytes_ram, 
message_bytes_ready, message_bytes_unacknowledged, messages, messages_persistent, messages_ram, messages_ready, messages_ready_ram, 
messages_unacknowledged, messages_unacknowledged_ram, mirror_pids, name, online, operator_policy, owner_pid, pid, policy, 
slave_pids, state, synchronised_mirror_pids, synchronised_slave_pids, type.

***show message count of each queue**

```bash
$ rabbitmqctl list_queues -p goapp-vhost name messages durable
Timeout: 60.0 seconds ...
Listing queues for vhost goapp-vhost ...
name	messages	durable
app.queue.b	12	true
app.queue.a	8	true
```
