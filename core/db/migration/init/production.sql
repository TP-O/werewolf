CREATE KEYSPACE IF NOT EXISTS werewolf WITH
    REPLICATION = {
        'class': 'NetworkTopologyStrategy',
        'replication_factor': 3,
    }
AND DURABLE_WRITES = true;
