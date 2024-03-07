from locust import HttpUser, task, between

class WebsiteUser(HttpUser):
    wait_time = between(1, 5)  # Simulate users waiting between 1 to 5 seconds between tasks

    headers = {
        'Accept-Encoding': 'gzip, deflate, br',
    }

    @task
    def get_markers(self):
        self.client.get("/api/v1/markers", headers=self.headers)