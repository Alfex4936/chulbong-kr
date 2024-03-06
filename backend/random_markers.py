import os

import mysql.connector
import numpy as np
from dotenv import load_dotenv

load_dotenv()

# Database connection parameters - replace these with your actual database information
db_config = {
    "host": os.getenv("DB_HOST"),
    "user": os.getenv("DB_USERNAME"),
    "password": os.getenv("DB_PASSWORD"),
    "database": os.getenv("DB_NAME"),
}


def connect_to_database(config):
    try:
        cnx = mysql.connector.connect(**config)
        return cnx
    except mysql.connector.Error as err:
        print(f"Failed to connect to database: {err}")
        return None


def drop_markers_table(cursor):
    try:
        cursor.execute("DROP TABLE IF EXISTS MarkersTest;")
        print("Markers table dropped successfully.")
    except mysql.connector.Error as err:
        print(f"Failed to drop Markers table: {err}")


def delete_all_markers(cursor):
    try:
        cursor.execute("DELETE FROM MarkersTest;")
        print("All rows from MarkersTest table deleted successfully.")
    except mysql.connector.Error as err:
        print(f"Failed to delete rows from MarkersTest table: {err}")


def generate_random_lat_lon():
    # South Korea coordinates limits
    lat_north = 38.45000000
    lat_south = 33.10000000
    lon_east = 131.87222222
    lon_west = 125.06666667

    lat = np.random.uniform(lat_south, lat_north)
    lon = np.random.uniform(lon_west, lon_east)

    return lat, lon


def insert_markers(cursor, num_markers):
    insert_query = """
INSERT INTO MarkersTest (Location, Description) 
VALUES (ST_GeomFromText('POINT(%s %s)', 4326), '');
"""

    batch_size = 1000
    values = []

    for _ in range(num_markers):
        lat, lon = generate_random_lat_lon()
        values.append((lat, lon))
        if len(values) >= batch_size:
            cursor.executemany(insert_query, values)
            values = []

    if values:  # Insert any remaining markers
        cursor.executemany(insert_query, values)

    print(f"{num_markers} markers inserted successfully.")


def main():
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


if __name__ == "__main__":
    main()
