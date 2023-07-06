AMQP: Advanced Message Queuing Protocol
========================

# Core concepts

**Broker or message broker**: a software which receives messages from one application or service, and delivers them to
another application, service, or broker.

**Virtual host**: a vhost is a group of resources, which exists within the broker. The resources include users, exchanges
queues and so on. Users can have different access privileges to different vhosts.

**Connection**: a physical network (TCP) between the application (producer/consumer) and broker.

**Channel**: a virtual connection insides a connection. The messages are published and consumed over a channel which reuses
a connection. One Connection can contain many channels.

**Exchange**: the exchange can apply routing rules for messages, making sure that messages are reaching their final destination.
Routing rules include *direct* (point to point), *topic* (publish-subscribe), *fanout* (multicast), and *header exchanges*.

**Queue**: a sequence of messages.

**Binding**: a virtual link between an exchange and a queue. It enables messages flow from an exchange to a queue.

# Traffic cost

Establishing a connection is costly. The handshake process for an AMQP connection requires at least seven TCP packets. Each
cost of actions is as follows.

| Action                 | Cost                                    |
|:-----------------------|:----------------------------------------|
| establish a connection | 7 TCP packets                           |
| create a channel       | 2 TCP packets                           |
| Publish a message      | 1 TCP packet (more for larger messages) |
| close a channel        | 2 TCP packets                           |
| close a connection     | 2 TCP packets                           |

