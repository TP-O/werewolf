CREATE KEYSPACE IF NOT EXISTS werewolf WITH
    REPLICATION = {
        'class': 'SimpleStrategy',
        'replication_factor': 1,
    }
AND DURABLE_WRITES = false;
