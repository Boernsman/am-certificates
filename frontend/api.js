//

const API_KEY = 'YOUR_API_KEY';  // Replace with your actual API key

/**
 * Function to validate a certificate code by calling the `/valide` endpoint
 * @param {string} baseURL - The base URL of the backend
 * @param {string} code - The certificate code to validate
 * @returns {Promise} - A promise that resolves to the server's response or an error
 */
export async function validateCertificateCode(baseURL, code) {
  try {
    const url = `${baseURL}/valide?code=${encodeURIComponent(code)}`;

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'X-API-KEY': API_KEY  // API key for authentication
      }
    });

    if (!response.ok) {
      const errorMessage = await response.text();
      throw new Error(`Error: ${response.status} - ${errorMessage}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Error validating certificate:', error);
    throw error;
  }
}

/**
 * Function to generate a certificate by calling the `/certificate` endpoint (POST)
 * @param {string} baseURL - The base URL of the Go backend
 * @param {Object} certificateData - The certificate data (name, email, code)
 * @returns {Promise} - A promise that resolves to the server's response or an error
 */
export async function generateCertificate(baseURL, certificateData) {
  try {
    const url = `${baseURL}/erstelle`;

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-KEY': API_KEY  // API key for authentication
      },
      body: JSON.stringify(certificateData)
    });

    if (!response.ok) {
      const errorMessage = await response.text();
      throw new Error(`Error: ${response.status} - ${errorMessage}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Error generating certificate:', error);
    throw error;
  }
}