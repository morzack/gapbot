# this is a quick script to convert between the new and old users.json format
import json

OLD_FILE = "users-old.json"
NEW_FILE = "users-new.json"

with open(OLD_FILE, 'r') as f:
    user_data = json.load(f)

for user in user_data['users']:
    user_data['users'][user]['last-fm'] = ""

with open(NEW_FILE, 'w') as f:
    json.dump(user_data, f)