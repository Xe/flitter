express = require "express"
getenv = require "getenv"
etcd = require "etcd"

app = do express

client = new etcd {
  url: "http://" + getenv("HOST", "172.17.42.1") + ":4001"
}

app.use require("morgan")("combined")

app.get "/", (req, res) ->
	res.send "hi\n"

server = app.listen getenv.int("PORT", 3000), ->
	console.log "Listening on port " + server.address().port
