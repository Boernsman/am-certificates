#!/usr/bin/env python3

import argparse
import base64
import csv
import configparser
import json
import requests

def load_credentials(ini_file='config.ini'):
    config = configparser.ConfigParser()
    config.read(ini_file)
    apiKey = config['api_keys']['key1']
    return apiKey

def generate_certificate(base_url, code, name, email, apiKey):

    headers = {
        'X-API-Key': f'Bearer {apiKey}'
    }
    data = {
        "code": code,
        "name": name,
        "email": email
    }
    json_payload = json.dumps(data)
    
    try:
        response = requests.post(f"{base_url}/generiere", headers=headers, data=json_payload)
        response.raise_for_status()  # Raise an error if the response status is not 200
        return response.json()  # Return the JSON response
    except requests.exceptions.RequestException as e:
        print(f"Error fetching data from backend: {e}")
        return None

def main():
    # Parse command-line arguments
    parser = argparse.ArgumentParser(description='Fetch certificate codes from backend and export them to CSV.')
    parser.add_argument('-u', '--url', type=str,
                        default='https://zertifikat.austromagnum.at', help='Base URL of the backend')
    parser.add_argument('-n', '--name',   type=str, required=True, help='Name of the certified person')
    parser.add_argument('-e', '--email',  type=str, required=True, help='Email of the certified person')
    parser.add_argument('-c', '--code',   type=str, required=True, help='Code of the certificate')
    parser.add_argument('-i', '--config', type=str, default='../tests/config.ini', help='Authentication credential file')
    
    args = parser.parse_args()
    
    apiKey = load_credentials(args.config)
    json_data = generate_certificate(args.url, args.code, args.name, args.email, apiKey)
    print(json_data)

if __name__ == "__main__":
    main()
