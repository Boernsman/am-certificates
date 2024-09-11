#!/usr/bin/env python3

import argparse
import base64
import csv
import configparser
import requests

def load_credentials(ini_file='config.ini'):
    config = configparser.ConfigParser()
    config.read(ini_file)

    # Fetch basic auth credentials
    username = config['auth']['username']
    password = config['auth']['password']

    return username, password

def fetch_codes_from_backend(base_url, certType, count, username, password):
    # Encode the credentials for basic auth
    auth_string = f'{username}:{password}'
    auth_bytes = auth_string.encode('utf-8')
    auth_base64 = base64.b64encode(auth_bytes).decode('utf-8')

    # Set headers for the request with Authorization
    headers = {
        'Authorization': f'Basic {auth_base64}'
    }
    
    try:
        response = requests.get(f"{base_url}/erstelle", headers=headers,
                                params={"count": count, "type": certType})
        response.raise_for_status()  # Raise an error if the response status is not 200
        return response.json()  # Return the JSON response
    except requests.exceptions.RequestException as e:
        print(f"Error fetching data from backend: {e}")
        return None

def export_to_csv(codes, output_file):
    """Export the JSON response (codes) to a CSV file."""
    with open(output_file, mode='w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(["Code", "URL"])  # Header row
        for code_data in codes:
            writer.writerow([code_data["code"], code_data["url"]])
    print(f"Codes successfully exported to {output_file}")

def main():
    # Parse command-line arguments
    parser = argparse.ArgumentParser(description='Fetch certificate codes from backend and export them to CSV.')
    parser.add_argument('-u', '--url', type=str,
                        default='https://zertifikat.austromagnum.at', help='Base URL of the backend')
    parser.add_argument('-c', '--count', type=int, required=True, help='Number of certificate codes to generate')
    parser.add_argument('-t', '--certtype', type=str, default="Kraftakt",
                        help='Type of certificate to generate')
    parser.add_argument('-o', '--output', type=str, default='codes.csv', help='Output CSV file (default: codes.csv)')
    parser.add_argument('-i', '--config', type=str, default='../tests/config.ini', help='Authentication credential file')
    
    args = parser.parse_args()
    
    username, password = load_credentials(args.config)
    json_data = fetch_codes_from_backend(args.url, args.certtype, args.count, username, password)
    
    # If valid data was returned, export it to CSV
    if json_data and "codes" in json_data:
        export_to_csv(json_data["codes"], args.output)
    else:
        print("No valid data received from backend.")

if __name__ == "__main__":
    main()
