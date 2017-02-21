# insert test data
host=localhost
port=8100
post_topic="curl -k --insecure -X POST https://$host:$port/topic -d"
$post_topic '{"id":"1","title":"Topic 1","author":"Rob","comments":0,"created":"2017-02-20T15:22:25.843881514+01:00","posts":[],"closed":true}'
$post_topic '{"id":"2","title":"Topic 2","author":"Dave","comments":0,"created":"2017-02-20T15:22:25.843881514+01:00","posts":[],"closed":true}'
$post_topic '{"id":"3","title":"Topic 3","author":"Jane","comments":0,"created":"2017-02-20T15:22:25.843881514+01:00","posts":[],"closed":true}'

post_topic1="curl -k --insecure -X POST https://localhost:8100/topic/1 -d"
$post_topic1 '{"id":"1","text":"Some text for Topic 1","author":"Rob","created":"2017-02-20T15:52:31.520895831+01:00"}'
$post_topic1 '{"id":"1","text":"Some response","author":"Jane","created":"2017-02-20T15:52:31.520895831+01:00"}'
$post_topic1 '{"id":"1","text":"Another response","author":"Jane","created":"2017-02-20T15:52:31.520895831+01:00"}'

post_topic2="curl -k --insecure -X POST https://localhost:8100/topic/2 -d"
$post_topic2 '{"id":"1","text":"Some text for Topic 2","author":"Dave","created":"2017-02-20T15:52:31.520895831+01:00"}'
$post_topic2 '{"id":"1","text":"Some response","author":"Jane","created":"2017-02-20T15:52:31.520895831+01:00"}'
