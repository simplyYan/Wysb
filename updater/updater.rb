require 'net/http'
require 'json'
require 'fileutils'

# Define URLs for the latest versions
REPO_URL = 'https://api.github.com/repos/simplyYan/Wysb/releases/latest'
FILES = {
  'wysbc_macos' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_macos',
  'wysbc_linux' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_linux',
  'wysbc_win32.exe' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_win32.exe',
  'wysbc_win64.exe' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_win64.exe'
}

def fetch_latest_version
  uri = URI(REPO_URL)
  response = Net::HTTP.get(uri)
  release_info = JSON.parse(response)
  release_info['tag_name']
end

def current_version(file_path)
  if File.exist?(file_path)
    # Adapt this according to how the version is stored
    # Assuming the version is in a file or part of the file
    File.read(file_path).strip
  else
    nil
  end
end

def download_file(url, file_path)
  uri = URI(url)
  response = Net::HTTP.get(uri)
  File.open(file_path, 'wb') do |file|
    file.write(response)
  end
end

def main
  latest_version = fetch_latest_version
  puts "The latest version available is: #{latest_version}"

  FILES.each do |filename, url|
    file_path = File.join(Dir.pwd, filename)
    local_version = current_version(file_path)

    puts "Checking '#{filename}'..."
    if local_version
      puts "Local version: #{local_version}"
      if local_version == latest_version
        puts "The file '#{filename}' is up to date."
      else
        puts "The file '#{filename}' is outdated."
        print "Do you want to update? (y/n): "
        answer = gets.chomp.downcase
        if answer == 'y'
          puts "Downloading the latest version..."
          download_file(url, file_path)
          puts "Update completed."
        end
      end
    else
      puts "File '#{filename}' not found."
      print "Do you want to download the latest version? (y/n): "
      answer = gets.chomp.downcase
      if answer == 'y'
        puts "Downloading the latest version..."
        download_file(url, file_path)
        puts "Download completed."
      end
    end
  end
end

main
