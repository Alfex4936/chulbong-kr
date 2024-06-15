from locust import HttpUser, task, between, events
import random
import os
import logging

# Configure logging
logging.basicConfig(level=logging.INFO)

class MyUser(HttpUser):
    wait_time = between(1, 5)  # wait time between tasks (1 to 5 seconds)

    @task
    def search_marker(self):
        terms = ["경기도", "인천광역시", "제주특별자치도"]
        term = random.choice(terms)
        self.client.get(f"/api/v1/search/marker?term={term}")

    def on_start(self):
        logging.info("Starting load test...")

    def on_stop(self):
        logging.info("Stopping load test...")

@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    logging.info("Test is starting...")

@events.test_stop.add_listener
def on_test_stop(environment, **kwargs):
    logging.info("Test is stopping...")

if __name__ == "__main__":
    os.system("locust -f loadtest.py --host http://localhost:8080")
