import csv
import qrcode
import os

# Function to generate QR code from a URL
def generate_qr_code(url, output_dir):
    # Create the QR code
    qr = qrcode.QRCode(
        version=1,  # Controls the size of the QR code
        error_correction=qrcode.constants.ERROR_CORRECT_L,  # Controls the error correction
        box_size=10,  # Size of each box in the QR code
        border=4,  # Thickness of the border
    )
    qr.add_data(url)
    qr.make(fit=True)

    # Create an image from the QR code
    img = qr.make_image(fill='black', back_color='white')

    # Save the QR code as an image file
    filename = os.path.join(output_dir, f"{url.replace('http://', '').replace('https://', '').replace('/', '_')}.png")
    img.save(filename)
    print(f"QR code saved as {filename}")

# Function to read URLs from CSV and generate QR codes
def generate_qr_codes_from_csv(csv_file, output_dir):
    # Create output directory if it doesn't exist
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    with open(csv_file, 'r') as file:
        reader = csv.reader(file)
        for row in reader:
            if len(row) > 0:
                url = row[0]  # Assuming the URL is in the first column
                generate_qr_code(url, output_dir)

csv_file = 'urls.csv'  # Path to your CSV file
output_dir = 'qr_codes'  # Directory where QR codes will be saved
generate_qr_codes_from_csv(csv_file, output_dir)

