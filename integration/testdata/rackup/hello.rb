class HelloWorld
  def call(env)
    [200, {"content-type" => "text/plain"}, ["Hello world!"]]
  end
end
