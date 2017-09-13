#!/usr/bin/ruby

require 'openssl'
require 'base64'

# This script uses HTTPie to post some example data to the server

def hmac_signature
  data = File.read(filename)
  digest = OpenSSL::Digest.new('sha256')
  "sha256=#{OpenSSL::HMAC.hexdigest(digest, key, data)}"
end

def key
  ARGV[0] || 'supersecret'
end

def filename
  ARGV[1] || File.join(File.dirname(__FILE__), "pull-request-event.json")
end

def run_command()
  command = <<CMD
http POST localhost:8080/v1/event @#{filename} \
  'X-Hub-Signature:#{hmac_signature}' \
  'X-Github-Event:pull_request'
CMD

  puts command
  system(command)
end

run_command
