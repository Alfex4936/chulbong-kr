import os
from pathlib import Path

from dotenv import load_dotenv


def str_to_bool(value):
    return value.lower() in {"true", "1", "yes", "on"}


# Determine the root of the project
BASE_DIR = Path(__file__).resolve().parent.parent.parent

# Load the .env file from the root directory
load_dotenv(dotenv_path=BASE_DIR / ".env")


# Access environment variables
DATABASE_URL = os.getenv("DATABASE_URL")
SECRET_KEY = os.getenv("SECRET_KEY")
DEBUG = str_to_bool(os.getenv("DEBUG", "False"))
