class ApplicationController < ActionController::Base
  def index
    render plain: 'Hello World!'
  end
end
