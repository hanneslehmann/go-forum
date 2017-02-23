# insert test data
host=localhost
port=8100

post_topic2="curl -k --insecure -X POST https://localhost:8100/topic/2 -d"
$post_topic2 '{"id":"1","text":"A NEW response","author":"Jane","created":"2017-02-20T15:52:31.520895831+01:00"}'

post_topic1="curl -k --insecure -X POST https://localhost:8100/topic/1 -d"
$post_topic1 '{"id":"1","text":"A NEW response","author":"Jane","created":"2017-02-20T15:52:31.520895831+01:00"}'
