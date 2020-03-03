---
title: Socket
---

Socket typed values hold one endpoint of a two way communication i.e. network/unix sockets, named-pipes etc...

Sockets...
- are immutable, i.e. assigning to an socket results in a copy of the original socket
- can be passed in/out of ops via [socket parameters](../structure/op-directory/op/parameter/socket.md)