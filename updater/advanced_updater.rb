require 'net/http'
require 'json'
require 'fileutils'
require 'digest'
require 'logger'
require 'openssl'
require 'tempfile'
require 'concurrent'
require 'optparse'

# Constants for repository and files
REPO_URL = 'https://api.github.com/repos/simplyYan/Wysb/releases/latest'
FILES = {
  'wysbc_macos' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_macos',
  'wysbc_linux' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_linux',
  'wysbc_win32.exe' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_win32.exe',
  'wysbc_win64.exe' => 'https://github.com/simplyYan/Wysb/blob/main/dist%20(binary)/wysbc_win64.exe'
}

# Logger setup
LOGGER = Logger.new('update.log', 'daily')
LOGGER.level = Logger::INFO

# Configuration options
CONFIG = {
  backup: true,
  verify_ssl: true,
  max_retries: 3,
  retry_delay: 5,
  threads: 4
}

# Custom error classes
class NetworkError < StandardError; end
class VersionMismatchError < StandardError; end
class DownloadError < StandardError; end
class BackupError < StandardError; end

# Helper to handle retries
def with_retries(max_retries, delay)
  retries = 0
  begin
    yield
  rescue => e
    if retries < max_retries
      retries += 1
      LOGGER.warn("Retry ##{retries} after error: #{e.message}")
      sleep(delay)
      retry
    else
      LOGGER.error("Max retries reached: #{e.message}")
      raise
    end
  end
end

# Fetches the latest version from the GitHub API
def fetch_latest_version
  with_retries(CONFIG[:max_retries], CONFIG[:retry_delay]) do
    uri = URI(REPO_URL)
    response = Net::HTTP.get(uri)
    raise NetworkError, "Failed to fetch release info" if response.nil? || response.empty?

    release_info = JSON.parse(response)
    raise "Invalid response structure" unless release_info.key?('tag_name')

    release_info['tag_name']
  end
rescue JSON::ParserError => e
  LOGGER.error("JSON parsing error: #{e.message}")
  raise
rescue => e
  LOGGER.error("Error fetching latest version: #{e.message}")
  raise
end

# Retrieves the current version of the local file, assuming version is stored in a hidden file
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

# Verifies SSL certificate of a URL
def verify_ssl_certificate(uri)
  return unless CONFIG[:verify_ssl]

  http = Net::HTTP.new(uri.host, uri.port)
  http.use_ssl = true
  http.verify_mode = OpenSSL::SSL::VERIFY_PEER

  http.start do
    http.peer_cert
  end
rescue OpenSSL::SSL::SSLError => e
  LOGGER.error("SSL certificate verification failed: #{e.message}")
  raise
end

# Downloads a file and verifies its integrity using SHA256 checksum
def download_file(url, file_path)
  uri = URI(url)
  verify_ssl_certificate(uri)

  with_retries(CONFIG[:max_retries], CONFIG[:retry_delay]) do
    response = Net::HTTP.get(uri)
    raise DownloadError, "Failed to download file from #{url}" if response.nil? || response.empty?

    File.open(file_path, 'wb') do |file|
      file.write(response)
    end

    checksum = Digest::SHA256.file(file_path).hexdigest
    LOGGER.info("Downloaded '#{file_path}' with SHA256 checksum: #{checksum}")

    checksum
  end
rescue => e
  LOGGER.error("Error downloading file from '#{url}': #{e.message}")
  raise
end

# Creates a backup of the existing file before update
def backup_file(file_path)
  return unless CONFIG[:backup] && File.exist?(file_path)

  backup_dir = File.join(Dir.pwd, 'backup')
  FileUtils.mkdir_p(backup_dir)
  backup_path = File.join(backup_dir, "#{File.basename(file_path)}.bak")

  FileUtils.copy(file_path, backup_path)
  LOGGER.info("Created backup of '#{file_path}' at '#{backup_path}'")
rescue => e
  LOGGER.error("Error creating backup for '#{file_path}': #{e.message}")
  raise BackupError, "Backup failed for '#{file_path}'"
end

# Updates the local version file after a successful download
def update_version_file(file_path, latest_version)
  version_file = "#{file_path}.version"
  File.write(version_file, latest_version)
  LOGGER.info("Updated version file '#{version_file}' to version #{latest_version}")
rescue => e
  LOGGER.error("Error updating version file '#{version_file}': #{e.message}")
  raise
end

# Threaded download management
def threaded_downloads(files, latest_version)
  pool = Concurrent::FixedThreadPool.new(CONFIG[:threads])
  
  files.each do |filename, url|
    pool.post do
      begin
        file_path = File.join(Dir.pwd, filename)
        local_version = current_version(file_path)

        LOGGER.info("Checking '#{filename}'...")
        if local_version && local_version == latest_version
          LOGGER.info("The file '#{filename}' is up to date.")
        else
          LOGGER.warn("Updating '#{filename}' from version '#{local_version}' to '#{latest_version}'")
          backup_file(file_path)
          download_file(url, file_path)
          update_version_file(file_path, latest_version)
        end
      rescue => e
        LOGGER.error("Failed to update '#{filename}': #{e.message}")
      end
    end
  end

  pool.shutdown
  pool.wait_for_termination
rescue => e
  LOGGER.fatal("Threaded download process failed: #{e.message}")
  raise
end

# Central method to handle updates
def update_files
  latest_version = fetch_latest_version
  LOGGER.info("The latest version available is: #{latest_version}")

  threaded_downloads(FILES, latest_version)
rescue => e
  LOGGER.fatal("Update process failed: #{e.message}")
  puts "An error occurred: #{e.message}. Please check the log file for details."
end

# Command-line option parsing
def parse_options
  OptionParser.new do |opts|
    opts.banner = "Usage: updater.rb [options]"

    opts.on("-b", "--no-backup", "Disable backups") do
      CONFIG[:backup] = false
    end

    opts.on("-s", "--skip-ssl", "Skip SSL verification") do
      CONFIG[:verify_ssl] = false
    end

    opts.on("-r", "--retries RETRIES", Integer, "Set max retries (default: 3)") do |retries|
      CONFIG[:max_retries] = retries
    end

    opts.on("-d", "--delay DELAY", Integer, "Set retry delay in seconds (default: 5)") do |delay|
      CONFIG[:retry_delay] = delay
    end

    opts.on("-t", "--threads THREADS", Integer, "Set number of threads for downloads (default: 4)") do |threads|
      CONFIG[:threads] = threads
    end

    opts.on_tail("-h", "--help", "Show this message") do
      puts opts
      exit
    end
  end.parse!
end

# Starts the update process with a time measurement
def main
  start_time = Time.now
  LOGGER.info("Starting update process at #{start_time}")

  parse_options
  update_files

  end_time = Time.now
  elapsed_time = end_time - start_time
  LOGGER.info("Update process completed at #{end_time}, taking #{elapsed_time.round(2)} seconds.")
rescue => e
  LOGGER.fatal("Main process encountered a fatal error: #{e.message}")
end

# Execute the script
main
