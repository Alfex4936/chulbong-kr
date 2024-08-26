from locust import HttpUser, TaskSet, between, task


class UserBehavior(TaskSet):
    @task
    def get_root(self):
        self.client.get("/")


class WebsiteUser(HttpUser):
    tasks = [UserBehavior]
    wait_time = between(
        1, 5
    )  # Simulate a wait time between 1 and 5 seconds between tasks
