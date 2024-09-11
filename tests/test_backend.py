import requests

# Base URL of the Go backend API
BASE_URL = "http://localhost:8080"
#BASE_URL = "http://zertifikat.austromagnum.at"

# Test case for validating a certificate code
def test_validate_code_valid():
    """
    Test for validating a valid certificate code.
    """
    # Replace with a valid certificate code for testing
    valid_code = "validcode123"
    
    response = requests.get(f"{BASE_URL}/certificate", params={"code": valid_code})
    
    # Check if the response status code is 200 OK
    assert response.status_code == 200, f"Expected 200, got {response.status_code}"
    
    # Validate the response message
    response_json = response.json()
    assert "message" in response_json, "Response doesn't contain 'message'"
    assert response_json["message"] == "Code valid. Please submit your details."


# Test case for validating an invalid certificate code
def test_validate_code_invalid():
    """
    Test for validating an invalid certificate code.
    """
    invalid_code = "invalidcode456"
    
    response = requests.get(f"{BASE_URL}/certificate", params={"code": invalid_code})
    
    # Check if the response status code is 404 Not Found
    assert response.status_code == 404, f"Expected 404, got {response.status_code}"
    
    # Validate the error message
    response_json = response.json()
    assert "error" in response_json, "Response doesn't contain 'error'"
    assert response_json["error"] == "Invalid or used code"


# Test case for generating a certificate
def test_generate_certificate():
    """
    Test for generating a certificate with valid name, email, and code.
    """
    payload = {
        "code": "validcode123",  # Replace with a valid code for testing
        "name": "John Doe",
        "email": "johndoe@example.com"
    }
    
    response = requests.post(f"{BASE_URL}/certificate", json=payload)
    
    # Check if the response status code is 200 OK
    assert response.status_code == 200, f"Expected 200, got {response.status_code}"
    
    # Validate the success message and the jpeg URL
    response_json = response.json()
    assert "message" in response_json, "Response doesn't contain 'message'"
    assert response_json["message"] == "Certificate generated successfully!"
    assert "jpeg_url" in response_json, "Response doesn't contain 'jpeg_url'"
    assert response_json["jpeg_url"].startswith("/certificates/"), "Invalid jpeg URL format"


# Test case for generating a certificate with missing fields
def test_generate_certificate_missing_fields():
    """
    Test for generating a certificate with missing fields in the request.
    """
    payload = {
        "code": "validcode123",  # Replace with a valid code for testing
        "name": "John Doe"
        # Missing email
    }
    
    response = requests.post(f"{BASE_URL}/certificate", json=payload)
    
    # Check if the response status code is 400 Bad Request
    assert response.status_code == 400, f"Expected 400, got {response.status_code}"
    
    # Validate the error message
    response_json = response.json()
    assert "error" in response_json, "Response doesn't contain 'error'"
    assert response_json["error"] == "Invalid request"

