import json
import os
import random
from datetime import datetime, timedelta

import mysql.connector
import numpy as np
from dotenv import load_dotenv

load_dotenv()

korean_reviews = [
    "아파트 단지 안에 있어요.",
    "ㅋㅋ",
    "",
    "그냥 철봉",
    "밤에 어두워요 ㅠㅠ",
    "추천 안함",
    "모래랑 있어요~",
    "약간 삐뚤어요...",
    "철봉 3개 있어요",
    "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
]


# Database connection parameters - replace these with your actual database information
db_config = {
    "host": os.getenv("DB_HOST"),
    "user": os.getenv("DB_USERNAME"),
    "password": os.getenv("DB_PASSWORD"),
    "database": os.getenv("DB_NAME"),
}


def generate_random_review():
    # Randomly select a review from the predefined list
    return random.choice(korean_reviews)


def connect_to_database(config):
    try:
        cnx = mysql.connector.connect(**config)
        return cnx
    except mysql.connector.Error as err:
        print(f"Failed to connect to database: {err}")
        return None


def drop_markers_table(cursor):
    try:
        cursor.execute("DROP TABLE IF EXISTS Markers;")
        print("Markers table dropped successfully.")
    except mysql.connector.Error as err:
        print(f"Failed to drop Markers table: {err}")


def delete_all_markers(cursor):
    try:
        cursor.execute("DELETE FROM Markers;")
        print("All rows from MarkersTest table deleted successfully.")
    except mysql.connector.Error as err:
        print(f"Failed to delete rows from MarkersTest table: {err}")


def generate_random_lat_lon():
    # South Korea coordinates limits
    lat_north = 38.61500000
    lat_south = 33.10000000
    lon_east = 132.00000000
    lon_west = 124.00000000

    lat = np.random.uniform(lat_south, lat_north)
    lon = np.random.uniform(lon_west, lon_east)

    return lat, lon


def insert_markers(cursor, num_markers):
    insert_query = """
    INSERT INTO Markers (Location, Description) 
    VALUES (ST_GeomFromText('POINT(%s %s)', 4326), %s);
    """
    values = []

    for _ in range(num_markers):
        lat, lon = generate_random_lat_lon()
        # Generate a random review text for each marker
        description = generate_random_review()
        values.append((lat, lon, description))

    # Insert markers in batches
    batch_size = 1000
    for i in range(0, len(values), batch_size):
        batch = values[i : i + batch_size]
        cursor.executemany(insert_query, batch)

    print(f"{num_markers} markers inserted successfully.")


def insert_markers_json(cursor, markers):
    insert_query = """
    INSERT INTO Markers (UserID, Location, Description, CreatedAt) 
    VALUES (1, ST_GeomFromText('POINT(%s %s)', 4326), '', %s);
    """
    values = []

    for marker in markers:
        lat, lon, created_at = marker
        # Generate a random review text for each marker
        values.append((lat, lon, created_at))

    # Insert markers in batches
    batch_size = 1000
    for i in range(0, len(values), batch_size):
        batch = values[i : i + batch_size]
        cursor.executemany(insert_query, batch)

    print(f"{len(markers)} markers inserted successfully.")


# Function to process each marker and extract/transform needed data
def process_markers(json_data):
    processed_markers = []
    for marker in json_data:
        try:
            latitude = float(marker["latitude"])
            longitude = float(marker["longitude"])
        except:
            continue
        # Transform the date, set time to 00:00:00 if date is present
        date = marker["date"]
        if date:  # Check if date is not empty
            date_time = datetime.strptime(date, "%Y-%m-%d")
            # Formatting to include time as 00:00:00
            date_str = date_time.strftime("%Y-%m-%d %H:%M:%S")
        else:
            # Set to one year minus from now with time as 00:00:00
            one_year_ago = datetime.now() - timedelta(days=365)
            date_str = one_year_ago.strftime("%Y-%m-%d %H:%M:%S")

        processed_markers.append((latitude, longitude, date_str))

    return processed_markers


# Function to read JSON data from a file
def read_json_file(file_path):
    with open(file_path, "r", encoding="utf-8") as file:
        return json.load(file)


def mysql_query():
    # Connect to the database
    cnx = connect_to_database(db_config)
    if cnx is None:
        return

    cursor = cnx.cursor()

    delete_all_markers(cursor)

    # Insert markers
    insert_markers(cursor, 10000)

    # Commit the transactions and close the connection
    cnx.commit()
    cursor.close()
    cnx.close()


def json_insert():
    # Specify the path to your JSON file
    file_path = "markers.json"

    # Read the JSON data from the file
    json_data = read_json_file(file_path)

    # Process the markers
    processed_markers = process_markers(json_data)

    # Connect to the database
    cnx = connect_to_database(db_config)
    if cnx is None:
        return

    cursor = cnx.cursor()

    # delete_all_markers(cursor)

    # Insert markers
    insert_markers_json(cursor, processed_markers)

    # Commit the transactions and close the connection
    cnx.commit()
    cursor.close()
    cnx.close()


def main():
    # mysql_query()
    json_insert()


if __name__ == "__main__":
    main()
