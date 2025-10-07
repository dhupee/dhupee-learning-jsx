import os
import random
import string


# Generate random 8-character hash
def generate_random_hash():
    return "".join(random.choices(string.ascii_lowercase + string.digits, k=8))


# Create the output directory
def create_output_dir():
    random_hash = generate_random_hash()
    output_dir = f"/tmp/{random_hash}-yt-dlp-output"
    os.makedirs(output_dir, exist_ok=True)
    print(f"Created output directory: {output_dir}")
    return output_dir
