require 'net/http'
require 'json'
require 'fileutils'
require 'digest'
require 'logger'

REPO_URL = 'https://api.github.com/repos/simplyYan/Wysb/releases/latest'
FILES = {
  'wysbc_macos' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_macos',
  'wysbc_linux' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_linux',
  'wysbc_win32.exe' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_win32.exe',
  'wysbc_win64.exe' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_win64.exe'
}

LOGGER = Logger.new('update.log')
LOGGER.level = Logger::INFO

class NetworkError < StandardError; end
class VersionMismatchError < StandardError; end
class DownloadError < StandardError; end

def fetch_latest_version
  uri = URI(REPO_URL)
  response = Net::HTTP.get(uri)
  raise NetworkError, "Failed to fetch release info" if response.nil? || response.empty?

  release_info = JSON.parse(response)
  raise "Invalid response structure" unless release_info.key?('tag_name')

  release_info['tag_name']
rescue JSON::ParserError => e
  LOGGER.error("JSON parsing error: #{e.message}")
  raise
rescue => e
  LOGGER.error("Error fetching latest version: #{e.message}")
  raise
end

def current_version(file_path)
  version_file = "#{file_path}.version"
  if File.exist?(version_file)
    File.read(version_file).strip
  else
    nil
  end
rescue => e
  LOGGER.error("Error reading current version from '#{version_file}': #{e.message}")
  nil
end

def download_file(url, file_path)
  uri = URI(url)
  response = Net::HTTP.get(uri)
  raise DownloadError, "Failed to download file from #{url}" if response.nil? || response.empty?

  File.open(file_path, 'wb') do |file|
    file.write(response)
  end

  checksum = Digest::SHA256.file(file_path).hexdigest
  LOGGER.info("Downloaded '#{file_path}' with SHA256 checksum: #{checksum}")

  checksum
rescue => e
  LOGGER.error("Error downloading file from '#{url}': #{e.message}")
  raise
end

def update_version_file(file_path, latest_version)
  version_file = "#{file_path}.version"
  File.write(version_file, latest_version)
  LOGGER.info("Updated version file '#{version_file}' to version #{latest_version}")
rescue => e
  LOGGER.error("Error updating version file '#{version_file}': #{e.message}")
  raise
end

def update_files
  latest_version = fetch_latest_version
  LOGGER.info("The latest version available is: #{latest_version}")

  FILES.each do |filename, url|
    file_path = File.join(Dir.pwd, filename)
    local_version = current_version(file_path)

    begin
      puts "Checking '#{filename}'..."
      if local_version
        LOGGER.info("Local version of '#{filename}': #{local_version}")
        if local_version == latest_version
          puts "The file '#{filename}' is up to date."
          LOGGER.info("File '#{filename}' is up to date.")
        else
          raise VersionMismatchError, "Version mismatch for '#{filename}'"
        end
      else
        puts "File '#{filename}' not found."
        LOGGER.warn("File '#{filename}' not found.")
        raise VersionMismatchError, "Missing file '#{filename}'"
      end
    rescue VersionMismatchError
      print "Do you want to download/update the file '#{filename}'? (y/n): "
      answer = gets.chomp.downcase
      if answer == 'y'
        puts "Downloading the latest version..."
        download_file(url, file_path)
        update_version_file(file_path, latest_version)
        puts "Update completed."
        LOGGER.info("File '#{filename}' updated to version #{latest_version}.")
      else
        puts "Skipping update for '#{filename}'."
        LOGGER.info("Update for '#{filename}' was skipped by the user.")
      end
    end
  end
rescue => e
  LOGGER.fatal("Update process failed: #{e.message}")
  puts "An error occurred: #{e.message}. Please check the log file for details."
end

def main
  start_time = Time.now
  LOGGER.info("Starting update process at #{start_time}")

  update_files

  end_time = Time.now
  elapsed_time = end_time - start_time
  LOGGER.info("Update process completed at #{end_time}, taking #{elapsed_time.round(2)} seconds.")
rescue => e
  LOGGER.fatal("Main process encountered a fatal error: #{e.message}")
end

main
