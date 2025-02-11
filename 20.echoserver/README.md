# これは何？

シンプルなエコーサーバを用意し、bashの```/dev/tcp```ファイル機能を利用して通信をするサンプルです。

```sh
$ task
task: [default] go build -o echoserver .
task: [default] ./echoserver &
task: [default] sleep 1
05:14:24 Starting server at :8888
task: [default] ./client.sh
05:14:25 Accept client (127.0.0.1:46332)
helloworld
05:14:26 Close  client (127.0.0.1:46332)
task: [default] pkill -INT echoserver
05:14:26 Shutdown server started
05:14:26 Shutdown server completed
```
