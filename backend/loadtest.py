import random

from locust import HttpUser, between, task


class WebsiteUser(HttpUser):
    wait_time = between(
        1, 5
    )  # Simulate users waiting between 1 to 5 seconds between tasks

    # Define the range for latitude and longitude within South Korea
    SouthKoreaMinLat = 33.0
    SouthKoreaMaxLat = 38.615
    SouthKoreaMinLong = 124.0
    SouthKoreaMaxLong = 132.0

    headers = {
        "Accept-Encoding": "gzip, deflate, br",
    }

    @task
    def get_markers(self):
        # Generate random latitude and longitude within the specified range
        latitude = random.uniform(self.SouthKoreaMinLat, self.SouthKoreaMaxLat)
        longitude = random.uniform(self.SouthKoreaMinLong, self.SouthKoreaMaxLong)

        # Specify the distance parameter
        distance = 5000  # in meters

        # Construct the query parameters with the random latitude, longitude, and distance
        params = {"latitude": latitude, "longitude": longitude, "distance": distance}

        # Make a GET request with the query parameters
        self.client.get("/api/v1/markers/close", params=params)
