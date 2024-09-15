#!/bin/bash

# Function to display program information
function show_info {
    echo "Program Name: Wysb 32Toaster"
    echo "Creator: Wesley Yan Soares Brehmer"
    echo "All rights reserved."
}

# Function to check system architecture
function check_architecture {
    local arch
    arch=$(uname -m)
    if [[ "$arch" == "x86_64" ]]; then
        echo "64-bits"
    elif [[ "$arch" == "i386" ]] || [[ "$arch" == "i686" ]]; then
        echo "32-bits"
    else
        echo "Unknown architecture"
    fi
}

# Function to download files
function download_file {
    local url=$1
    local dest=$2

    echo "Downloading file from $url..."
    curl -L -o "$dest" "$url"
}

# Function to validate user input
function validate_input {
    local input=$1
    local prompt=$2
    local valid_options=$3

    while true; do
        read -rp "$prompt" input
        if [[ "$valid_options" == *"$input"* ]]; then
            echo "$input"
            break
        else
            echo "Invalid input. Please try again."
        fi
    done
}

# Display program information
show_info

# Ask user for their system type
system_type=$(validate_input "" "What is your system type (Linux or macOS)? " "Linux macOS")

# Ask user for system architecture
arch=$(validate_input "" "What is your system architecture (32-bits or 64-bits)? " "32-bits 64-bits")

if [[ "$arch" == "64-bits" ]]; then
    echo "You don't need to use this program. Please proceed to https://github.com/simplyYan/Wysb?tab=readme-ov-file#-how-to-install"
    exit 0
fi

# Ask user if they agree with the license
read -rp "Do you agree with the license? (Read at https://github.com/simplyYan/Wysb/blob/main/LICENSE) (y/n) " agree

if [[ "$agree" != "y" && "$agree" != "Y" ]]; then
    echo "You must agree to the license to continue. Program terminated."
    exit 1
fi

# Determine download URL based on the system
if [[ "$system_type" == "macOS" ]]; then
    file_url="https://github.com/simplyYan/Wysb/raw/main/dist%20(binary)/wysbc_macos32"
elif [[ "$system_type" == "Linux" ]]; then
    file_url="https://github.com/simplyYan/Wysb/raw/main/dist%20(binary)/wysbc_linux32"
fi

# Download the file
temp_file="wysbc_temp_download"
download_file "$file_url" "$temp_file"

# Ask user for installation directory
while true; do
    read -rp "Which directory would you like to install the program in? (Make sure the directory exists) " install_dir
    if [[ -d "$install_dir" ]]; then
        break
    else
        echo "The specified directory does not exist. Please try again."
    fi
done

# Move the file to the installation directory
mv "$temp_file" "$install_dir"

echo "Thank you for downloading Wysb 32Toaster!"

exit 0
