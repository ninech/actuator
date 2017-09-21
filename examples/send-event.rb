#!/usr/bin/ruby

require 'optparse'
require 'openssl'
require 'base64'

Options = Struct.new(:filename, :secret)

class Parser
  def self.parse(options)
    args = Options.new

    opt_parser = OptionParser.new do |opts|
      opts.banner = "Usage: send-event.rb [options]"

      opts.on("-fFILE", "--file=FILE", "JSON file with example hook") do |filename|
        args.filename = filename
      end

      opts.on("-sSECRET", "--secret=SECRET", "Webhook secret to share with Actuator") do |secret|
        args.secret = secret
      end

      opts.on("-h", "--help", "Prints this help") do
        puts opts
        exit
      end
    end

    opt_parser.parse!(options)
    args
  end
end

class WebhookShot
  attr_accessor :filename, :secret

  def initialize(filename, secret)
    @filename = filename
    @secret = secret
  end

  def shoot
    system <<-CMD
  http POST localhost:8080/v1/event @#{filename} \
    'X-Hub-Signature:#{hmac_signature}' \
    'X-Github-Event:pull_request'
  CMD
  end

  private

  def hmac_signature
    data = File.read(filename)
    digest = OpenSSL::Digest.new('sha256')
    signature = "sha256=#{OpenSSL::HMAC.hexdigest(digest, secret, data)}"
  end
end

options = Parser.parse ARGV

if options.filename.nil? || options.secret.nil?
  puts 'Please provide all arguments'
  exit 1
end

WebhookShot.new(options.filename, options.secret).shoot
