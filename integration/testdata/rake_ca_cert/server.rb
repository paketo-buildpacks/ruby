#!/usr/local/bin/ruby

require 'webrick'
require 'webrick/https'
require 'openssl'

# server key pair
key = OpenSSL::PKey::RSA.new(File.read('certs/key.pem'))
cert = OpenSSL::X509::Certificate.new(File.read('certs/cert.pem'))

# configure cert store and make sure it uses ca certs provided by buildpack
cert_store = OpenSSL::X509::Store.new
cert_store.set_default_paths

server = WEBrick::HTTPServer.new(
  Port: 8080,
  Logger: WEBrick::Log.new($stderr, WEBrick::Log::WARN),
  SSLEnable: true,
  SSLVerifyClient: OpenSSL::SSL::VERIFY_PEER | OpenSSL::SSL::VERIFY_FAIL_IF_NO_PEER_CERT,
  SSLCertificate: cert,
  SSLPrivateKey: key,
  SSLCertificateStore: cert_store,
)

server.mount_proc '/' do |req, res|
  res.body = 'Hello world, Authenticated User!'
end

trap 'INT' do server.shutdown end

server.start
