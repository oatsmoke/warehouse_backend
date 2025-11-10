CREATE TABLE IF NOT EXISTS kafka_stream
(
    action String
)
    ENGINE = Kafka
        SETTINGS
            kafka_broker_list = 'kafka:9094',
            kafka_topic_list = 'actions',
            kafka_group_name = 'clickhouse_group',
            kafka_format = 'JSONEachRow',
            kafka_num_consumers = 1;

CREATE TABLE IF NOT EXISTS events
(
    action String
)
    ENGINE = MergeTree()
        ORDER BY tuple();

CREATE MATERIALIZED VIEW IF NOT EXISTS kafka_to_events
    TO events
AS
SELECT *
FROM kafka_stream;
