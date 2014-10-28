package = "drill"
version = "0.1-2"

source = {
  url = "http://yolo-swag.com"
}

description = {
  summary = "The drill used to create the heavens",
  detailed = [[
   A simple client to Flitter/Lagann.
  ]],
  license = "Zlib",
  homepage = "http://github.com/Xe/flitter",
  maintainer = "xena@yolo-swag.com"
}


dependencies = {
  "penlight",
  "elfs",
  "dkjson",
  "luasocket",
  "moonscript"
}


build = {
  type = "none",
  install = {
    bin = {
      "drill"
    }
  }
}

