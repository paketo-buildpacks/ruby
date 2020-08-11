use Rack::ContentLength

app = proc do |env|
  [200, {'Content-Type' => 'text/plain'}, ['Hello world!']]
end

run app
