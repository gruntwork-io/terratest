# A simple web app built on top of Ruby and Sinatra.

require 'sinatra'
require 'json'

if ARGV.length != 2
  raise 'Expected exactly two arguments: SERVER_PORT SERVER_TEXT'
end

server_port = ARGV[0]
server_text = ARGV[1]

set :port, server_port
set :bind, '0.0.0.0'
set :server, 'puma'

get '/' do
  server_text
end
