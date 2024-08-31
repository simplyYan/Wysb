require 'open3'
require 'logger'
require 'fileutils'
require 'optparse'

# Logger configuration
LOGGER = Logger.new('python_installer.log', 'daily')
LOGGER.level = Logger::INFO

# Configuration options
CONFIG = {
  python_version_command: 'python3 --version',
  pyinstaller_check_command: 'python3 -m PyInstaller --version',
  python_install_command: 'sudo apt-get install -y python3',
  pip_install_command: 'sudo apt-get install -y python3-pip',
  pyinstaller_install_command: 'python3 -m pip install pyinstaller',
  max_retries: 3,
  retry_delay: 5
}

# Custom error classes
class CommandExecutionError < StandardError; end
class InstallationError < StandardError; end

# Helper to run system commands with retries and logging
def run_command(command)
  retries = 0
  begin
    stdout, stderr, status = Open3.capture3(command)
    LOGGER.info("Running command: #{command}")

    if status.success?
      LOGGER.info("Command succeeded: #{command}")
      stdout.strip
    else
      LOGGER.error("Command failed: #{command}")
      raise CommandExecutionError, stderr.strip
    end
  rescue => e
    if retries < CONFIG[:max_retries]
      retries += 1
      LOGGER.warn("Retry ##{retries} for command: #{command} after error: #{e.message}")
      sleep(CONFIG[:retry_delay])
      retry
    else
      LOGGER.fatal("Max retries reached for command: #{command}. Error: #{e.message}")
      raise
    end
  end
end

# Checks if Python is installed by getting its version
def check_python_version
  version_output = run_command(CONFIG[:python_version_command])
  version = version_output.match(/Python\s(\d+\.\d+\.\d+)/)
  if version
    LOGGER.info("Python version found: #{version[1]}")
    puts "Python version: #{version[1]}"
    version[1]
  else
    LOGGER.error("Unable to determine Python version from output: #{version_output}")
    nil
  end
rescue => e
  LOGGER.error("Failed to check Python version: #{e.message}")
  nil
end

# Checks if PyInstaller is installed
def check_pyinstaller_installed
  pyinstaller_output = run_command(CONFIG[:pyinstaller_check_command])
  if pyinstaller_output.include?('PyInstaller')
    LOGGER.info("PyInstaller is installed: #{pyinstaller_output}")
    puts "PyInstaller version: #{pyinstaller_output}"
    true
  else
    LOGGER.warn("PyInstaller not installed.")
    false
  end
rescue CommandExecutionError => e
  LOGGER.warn("PyInstaller check failed, likely not installed: #{e.message}")
  false
rescue => e
  LOGGER.error("Unexpected error during PyInstaller check: #{e.message}")
  false
end

# Installs Python using the system package manager
def install_python
  LOGGER.info("Attempting to install Python...")
  run_command(CONFIG[:python_install_command])
  run_command(CONFIG[:pip_install_command])
  LOGGER.info("Python installation completed successfully.")
rescue => e
  LOGGER.fatal("Failed to install Python: #{e.message}")
  raise InstallationError, "Python installation failed"
end

# Installs PyInstaller using pip
def install_pyinstaller
  LOGGER.info("Attempting to install PyInstaller...")
  run_command(CONFIG[:pyinstaller_install_command])
  LOGGER.info("PyInstaller installation completed successfully.")
rescue => e
  LOGGER.fatal("Failed to install PyInstaller: #{e.message}")
  raise InstallationError, "PyInstaller installation failed"
end

# Parses command-line options
def parse_options
  OptionParser.new do |opts|
    opts.banner = "Usage: python_checker.rb [options]"

    opts.on("-r", "--retries RETRIES", Integer, "Set max retries (default: 3)") do |retries|
      CONFIG[:max_retries] = retries
    end

    opts.on("-d", "--delay DELAY", Integer, "Set retry delay in seconds (default: 5)") do |delay|
      CONFIG[:retry_delay] = delay
    end

    opts.on("-h", "--help", "Show this help message") do
      puts opts
      exit
    end
  end.parse!
end

# Main process to check and install Python and PyInstaller
def main
  start_time = Time.now
  LOGGER.info("Starting Python/PyInstaller check at #{start_time}")
  
  parse_options

  python_version = check_python_version
  if python_version.nil?
    LOGGER.warn("Python not found. Installing Python...")
    install_python
    python_version = check_python_version
    if python_version.nil?
      raise InstallationError, "Python installation failed and version could not be verified."
    end
  end

  unless check_pyinstaller_installed
    LOGGER.warn("PyInstaller not found. Installing PyInstaller...")
    install_pyinstaller
    unless check_pyinstaller_installed
      raise InstallationError, "PyInstaller installation failed."
    end
  end

  end_time = Time.now
  elapsed_time = end_time - start_time
  LOGGER.info("Process completed successfully at #{end_time}, taking #{elapsed_time.round(2)} seconds.")
  puts "Python and PyInstaller are correctly installed and verified."
rescue InstallationError => e
  LOGGER.fatal("Installation process failed: #{e.message}")
  puts "An error occurred during the installation process: #{e.message}. Please check the log file for details."
rescue => e
  LOGGER.fatal("Unexpected error in main process: #{e.message}")
  puts "An unexpected error occurred: #{e.message}. Please check the log file for details."
end

# Run the main process
main
