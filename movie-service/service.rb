
require 'json'
require 'sinatra'
require 'logger'
require 'aphorism'
require 'color-generator'

class MovieService < Sinatra::Base
	set :show_exceptions, false

	configure do
    LOGGER = Logger.new("sinatra.log")
    enable :logging, :dump_errors
    set :raise_errors, true
  end

	def initialize
    super
    @logger = ::Logger.new(service_log_file)
    @generator = ColorGenerator.new saturation: 0.5, lightness: 0.5
  end

	get '/movies' do
		read_data
	end

	get '/movies/:id' do
		movies = JSON.parse(read_data)
		movie = movies.select { |movie| movie["imdbID"] == params[:id]}.first
		unless movie.nil?
			return movie.to_json
		end

		return {:error => "Movie with requested ID not found"}.to_json
	end

  get '/' do
    color = @generator.create_hex
    orator = Aphorism::Orator.new
    message = orator.say.split("-")
    "</br></br></br></br></br></br><p
      style=\"text-align: center;
      font-size: 2em; color: #{color}\">
      #{message[0]} </br> - <span style=\"font-style: italic;\">#{message[1]} <span>
    </p>"
  end

	not_found do
	  "There is nothing to do here.! 404!"
  end

  error 500 do
    "Something went wrong!!"
  end

	private
	def read_data
		open("movies.json").read
	end

	private
  	def service_log_file
	    file = ::File.new("payment.log","a+")
	    file.sync = true
	    file
  	end
end
