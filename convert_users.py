# this is a quick script to convert between the new and old users.json format
import json

OLD_FILE = "users-old.json"
NEW_FILE = "users-new.json"

with open(OLD_FILE, 'r') as f:
    user_data = json.load(f)

output_data = {
    "names-channel": user_data["names-channel"]
}

output_data["users"] = {old_user: {
    "id": old_user,
    "first-name": user_data["users"][old_user].split(" ")[0],
    "last-name": " ".join(user_data["users"][old_user].split(" ")[1:]),
    "grade": -1
} for old_user in user_data["users"]}

with open(NEW_FILE, 'w') as f:
    json.dump(output_data, f)