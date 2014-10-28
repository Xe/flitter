json = require "dkjson"
http = require "socket.http"
elfs = require "elfs"

stringx = require "pl.stringx"
stringx.import()

export STORAGE_PATH = os.getenv("HOME") .. "/.drill.json"

httpPost = (url, table) ->
  data = json.encode(table)

  http.request url, data

--- doCommand returns the output and return status of a command.
doCommand = (command) -> --> string, number
  n = os.tmpname!
  code = os.execute command .. " 2>/dev/null > " .. n
  lines = {}

  for line in io.lines n
    table.insert lines, line

  os.remove n

  lines, code

store = (lagann) ->
  data = {
    :lagann
    username: os.getenv "USER"
  }

  out = json.encode data

  do
    fout = io.open STORAGE_PATH, "w"
    fout\write out
    fout\close!

  true

restore = ->
  ret = {}

  do
    fin = io.open STORAGE_PATH, "r"
    temp = fin\read "*a"
    fin\close!

    ret = json.decode temp

  ret

--- pairsByKeys is an iterator for a table aplhabetically
--  adapted from here: http://www.lua.org/pil/19.3.html
pairsByKeys = (t) -> --> function
  a = {}
  for n in pairs t
    table.insert a, n

  table.sort a

  i = 0
  iter = ->
    i += 1
    if a[i] == nil
      return nil
    else
      return a[i], t[a[i]]
  iter

export usage = ->
  print "drill version 0.1\n"

  print "Usage: drill <command> ...[args]\n"

  print "Available commands:"
  for name,cmd in pairsByKeys commands
    print "%12s   %s"\format(name, cmd[1])

  os.exit 1

export commands = {
  help: { "Print detailed help", "help [command]", (args) ->
    if #args == 1
      print "drill " .. commands[args[1]\lower!][2]
    else
      usage!
  }
  register: { "Registers with lagann", "register <lagann>", (args) ->
    if #args ~= 1
      print "Need host for lagann"

    sshkey = ""
    do
      fin = io.open os.getenv("HOME") .. "/.ssh/id_rsa.pub", "r"
      sshkey = fin\read "*a"
      fin\close!

    fingerprint = ""
    do
      proc = io.popen "ssh-keygen -lf " .. os.getenv("HOME") .. "/.ssh/id_rsa.pub"
      fp = proc\read "*a"
      proc\close!

      fingerprint = fp\split(" ")[2]

    user = {
      name: os.getenv "USER"
      sshkeys: {
        {
          key: sshkey
          :fingerprint
          comment: "Lagann! Spin on!"
        }
      }
    }

    body, code, header = httpPost(args[1].."/register", user)
    reply = json.decode body

    if reply.code ~= 200
      print "Error: code " .. reply.code
      if reply.data
        print reply.data

    print reply.message

    if reply.code == 200
      store args[1]
  }

  login: {"Logs into lagann", "login <lagann>", (args) ->
    user = {
      name: os.getenv "USER"
    }

    body, code, header = httpPost(args[1].."/login", user)
    reply = json.decode body

    if reply.code ~= 200
      print "Error: code " .. reply.code
      if reply.data
        print reply.data

    print reply.message

    if reply.code == 200
      store args[1]
  }

  create: {"Creates app in lagann", "create [appname]", (args) ->
    name = elfs.GenName!
    if #args > 0
      name = args[1]

    controller = restore!

    app = {
      name: name
      users: {controller.username}
    }

    enc = json.encode app

    output, _ = doCommand "curl --data '#{enc}' #{controller.lagann}/create"
    reply = json.decode table.concat output, " "

    if reply.code ~= 200
      print "Error: code " .. reply.code
      if reply.data
        print reply.data
      print reply.message

      os.exit 1

    print reply.message

    os.execute "git remote rm flitter ; git remote add flitter ssh://git@172.17.42.1:2232/#{app.name}.git"
  }
}

if #arg == 0
  usage!

command = arg[1]\lower!

if commands[command] ~= nil
  table.remove(arg, 1)

  commands[command][3] arg
else
  usage!