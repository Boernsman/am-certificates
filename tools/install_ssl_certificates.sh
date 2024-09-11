#!/bin/bash

# Paths where the certificate and key will be installed
CERT_DIR="/etc/ssl/certs"
KEY_DIR="/etc/ssl/private"

# Check if the script is run as root (required for setting permissions in /etc/ssl)
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root or with sudo"
    exit 1
fi

# Check if the required files are passed as arguments
if [ $# -ne 2 ]; then
    echo "Usage: $0 <path-to-public.crt> <path-to-private.key>"
    exit 1
fi

PUBLIC_CRT=$1
PRIVATE_KEY=$2

# Check if the provided certificate and key files exist
if [ ! -f "$PUBLIC_CRT" ]; then
    echo "Error: Public certificate file '$PUBLIC_CRT' does not exist."
    exit 1
fi

if [ ! -f "$PRIVATE_KEY" ]; then
    echo "Error: Private key file '$PRIVATE_KEY' does not exist."
    exit 1
fi

# Copy the public certificate to the certificate directory
echo "Installing public certificate..."
cp "$PUBLIC_CRT" "$CERT_DIR/"
if [ $? -ne 0 ]; then
    echo "Error: Failed to copy public certificate to $CERT_DIR"
    exit 1
fi

# Set appropriate permissions for the public certificate
echo "Setting permissions for public certificate..."
chmod 644 "$CERT_DIR/$(basename $PUBLIC_CRT)"
chown root:root "$CERT_DIR/$(basename $PUBLIC_CRT)"
if [ $? -ne 0 ]; then
    echo "Error: Failed to set permissions for public certificate."
    exit 1
fi

# Copy the private key to the key directory
echo "Installing private key..."
cp "$PRIVATE_KEY" "$KEY_DIR/"
if [ $? -ne 0 ]; then
    echo "Error: Failed to copy private key to $KEY_DIR"
    exit 1
fi

# Set appropriate permissions for the private key
echo "Setting permissions for private key..."
chmod 600 "$KEY_DIR/$(basename $PRIVATE_KEY)"
chown root:root "$KEY_DIR/$(basename $PRIVATE_KEY)"
if [ $? -ne 0 ]; then
    echo "Error: Failed to set permissions for private key."
    exit 1
fi

echo "Certificate and key have been installed successfully with the correct permissions."
echo "----- DONE -----"
