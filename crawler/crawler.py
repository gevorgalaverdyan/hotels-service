import requests
from bs4 import BeautifulSoup
import json
import re

# Function to scrape hotel details from the main page
def scrape_hotels_list(url):
    response = requests.get(url)
    soup = BeautifulSoup(response.content, 'html.parser')
    hotels = []

    # Find all hotel entries in the list
    for section in soup.find_all('div', class_='div-col'):
        for li in section.find_all('li'):
            hotel_name = li.a.text
            hotel_link = li.a['href']
            city = li.a.find_next('a').text
            province = li.find_previous('h2').text
            hotels.append({
                'name': hotel_name,
                'wikiLink': f"https://en.wikipedia.org{hotel_link}",
                'city': city,
                'province': province
            })
    return hotels

# Function to extract simplified coordinates
def extract_simplified_coordinates(coords_text):
    # Use regex to extract latitude and longitude in decimal format
    match = re.search(r"(\d+\.\d+)°N\s+([-\d]+\.\d+)°W", coords_text)
    if match:
        latitude = match.group(1)
        longitude = match.group(2)
        return f"{latitude}; {longitude}"
    return ""

# Function to scrape individual hotel details
def scrape_hotel_details(hotel_url):
    response = requests.get(hotel_url)
    soup = BeautifulSoup(response.content, 'html.parser')
    infobox = soup.find('table', class_='infobox vcard')
    details = {}

    if infobox is None:
        print(f"No infobox found for {hotel_url}")
        return details

    # Scrape image
    image_tag = infobox.find('img')
    if image_tag:
        details['image'] = f"https:{image_tag['src']}"

    # Scrape coordinates
    coords_tag = infobox.find('span', class_='geo-inline')
    if coords_tag:
        coords_text = coords_tag.text.strip()
        details['coordinates'] = extract_simplified_coordinates(coords_text)

    # Scrape website URL
    website_tag = infobox.find('span', class_='url')
    if website_tag and website_tag.a:
        details['website'] = website_tag.a['href']

    return details

# Main function to collect all data
def main():
    main_url = "https://en.wikipedia.org/wiki/List_of_hotels_in_Canada"
    hotels = scrape_hotels_list(main_url)

    for hotel in hotels:
        print(f"Scraping details for {hotel['name']}...")
        details = scrape_hotel_details(hotel['wikiLink'])
        hotel.update(details)

    # Save to JSON file
    with open('hotels.json', 'w') as f:
        json.dump(hotels, f, indent=4)

if __name__ == "__main__":
    main()
