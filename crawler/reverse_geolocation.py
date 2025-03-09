import json
import requests
import time

# Load the JSON data
with open('/home/balerion/personal-repos/hotels-api/crawler/hotels.json', 'r') as file:
    hotels = json.load(file)

# Function to get address from coordinates
def get_address(lat, lon):
    url = f"https://nominatim.openstreetmap.org/reverse?format=json&lat={lat}&lon={lon}"
    headers = {
        'User-Agent': 'YourAppName/1.0 (your@email.com)'  # Replace with your app name and contact info
    }
    print("Requesting:", url)
    response = requests.get(url, headers=headers)
    print("Response:", response.status_code)
    if response.status_code == 200:
        data = response.json()
        address_data = data.get('address', {})
        address = f"{address_data.get('house_number', '')} {address_data.get('road', '')}, {address_data.get('city', '')}, {address_data.get('state', '')}, {address_data.get('country', '')}".strip()
        return address
    else:
        print("Error:", response.status_code, response.text)  # Log the error
        return 'Address not found'

# Iterate over each hotel and fetch the address
updated_hotels = []
for hotel in hotels:
    coordinates = hotel.get('coordinates')
    if coordinates:
        print("Fetching address for", hotel.get('name'))
        lat, lon = coordinates.split('; ')
        address = get_address(lat, lon)
        print("Address:", address)
        hotel['address'] = address
        time.sleep(1)  # Add a 1-second delay between requests
    updated_hotels.append(hotel)

# Save the updated hotels data to a new file
with open('/home/balerion/personal-repos/hotels-api/crawler/updatedHotels.json', 'w') as file:
    json.dump(updated_hotels, file, indent=2)

print("Updated hotels data saved to updatedHotels.json")