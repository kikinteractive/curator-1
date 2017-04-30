Note: this is (obviously) a fork from the original curator-go. This fork was created to host a missing recipe - the ConnectionWrapper. rerfer to this commit - 8a0724b587dee7a06671fc12970f1d534bf410bd. The connection wrapper is a gutted version of the cache recipe. Most of the functionality was stripped away leaving just a wrapper around the zookeeper connection, which allows listening to connection-level events (reconnection etc). This functionality was needed for the xiphias-go-registry repo.

# What is Curator?

Curator n ˈkyoor͝ˌātər: a keeper or custodian of a museum or other collection - A ZooKeeper Keeper.

![curator](http://curator.apache.org/images/ph-quote.png "Patrick Hunt Quote")

## What is curator-go?

Curator-go is a Golang porting for Curator, which base on the [go-zookeeper](https://github.com/samuel/go-zookeeper/).

# Getting Started

## Learn ZooKeeper

Curator-go users are assumed to know ZooKeeper. A good place to start is [ZooKeeper Getting Started Guide](http://zookeeper.apache.org/doc/trunk/zookeeperStarted.html)

## Install Curator-go

> $ go get github.com/curator-go/curator

## Using Curator

Curator-go is available from [github.com](https://github.com/curator-go/curator). You can easily include Curator-go into your code.

```
import (
	"github.com/curator-go/curator"
)
```

## Getting a Connection
Curator uses Fluent Style. If you haven't used this before, it might seem odd so it's suggested that you familiarize yourself with the style.

Curator connection instances (CuratorFramework) are allocated from the CuratorFrameworkBuilder. You only need one CuratorFramework object for each ZooKeeper cluster you are connecting to:

```
curator.NewClient(connString, retryPolicy)
```

This will create a connection to a ZooKeeper cluster using default values. The only thing that you need to specify is the retry policy. For most cases, you should use:

```
retryPolicy := curator.NewExponentialBackoffRetry(time.Second, 3, 15*time.Second)

client := curator.NewClient(connString, retryPolicy)

client.Start()
defer client.Close()
```

The client must be started (and closed when no longer needed).

## Calling ZooKeeper Directly

Once you have a CuratorFramework instance, you can make direct calls to ZooKeeper in a similar way to using the raw ZooKeeper object provided in the ZooKeeper distribution. E.g.:

```
client.Create().ForPathWithData(path, payload)
```

The benefit here is that Curator manages the ZooKeeper connection and will retry operations if there are connection problems.

## Recipes
### Distributed Lock

```
lock := curator.NewInterProcessMutex(client, lockPath)

if ( lock.Acquire(maxWait, waitUnit) )
{
    defer lock.Release()

    // do some work inside of the critical section here
}
```

### Leader Election

```
listener := curator.NewLeaderSelectorListener(func(CuratorFramework client) error {
    // this callback will get called when you are the leader
    // do whatever leader work you need to and only exit
    // this method when you want to relinquish leadership
}))

selector := curator.NewLeaderSelector(client, path, listener)
selector.AutoRequeue()  // not required, but this is behavior that you will probably expect
selector.Start()
```

# Examples
This module contains example usages of various Curator features. Each directory in the module is a separate example.

- [leader](examples/leader/) Example leader selector code
- [cache](examples/cache/) Example PathChildrenCache usage
- [locking](examples/locking/) Example of using InterProcessMutex
- [discovery](examples/discovery/) Example usage of the Curator's ServiceDiscovery
- [framework](examples/framework/) A few examples of how to use the CuratorFramework class

See the [examples](examples/) source repo for each example.

# Components

- [Recipes](doc/recipes.md) Implementations of some of the common ZooKeeper "recipes". The implementations are built on top of the Curator Framework.
- [Framework](doc/framework.md) The Curator Framework is a high-level API that greatly simplifies using ZooKeeper. It adds many features that build on ZooKeeper and handles the complexity of managing connections to the ZooKeeper cluster and retrying operations.
- [Utilities](doc/utilities.md) Various utilities that are useful when using ZooKeeper.
- [Client](doc/client.md) A replacement for the bundled ZooKeeper class that takes care of some low-level housekeeping and provides some useful utilities.
- [Errors](doc/errors.md) How Curator deals with errors, connection issues, recoverable exceptions, etc.
