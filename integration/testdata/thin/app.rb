require 'sinatra'
configure { set :server, :thin }

get '/' do
  'Hello world!'
end
